# Pydevlake

A framework to write data collection plugins for [DevLake](https://devlake.apache.org/).


# How to create a new plugin

## Create plugin project


Make sure you have [Poetry](https://python-poetry.org/docs/#installation) installed.
Move to `python/plugins` and execute `poetry new myplugin`.
This will generate a new directory for your plugin.

In the `pyproject.toml` file and add the following line at the end of the `[tool.poetry.dependencies]` section:
```
pydevlake = { path = "../../pydevlake", develop = false }
```

Now run `poetry install`.

## Create `main` file

Create a `main.py` file with the following content:

```python
from pydevlake import Plugin, Connection


class MyPluginConnection(Connection):
    pass


class MyPlugin(Plugin):
    @property
    def connection_type(self):
        return MyPluginConnection

    def test_connection(self, connection: MyPluginConnection):
        pass

    @property
    def streams(self):
        return []


if __name__ == '__main__':
    MyPlugin.start()
```

This file is the entry point to your plugin.
It specifies three things:
- the parameters that your plugin needs to collect data, e.g. the url and credentials to connect to the datasource or custom options
- how to validate that some given parameters allows to connect to the datasource, e.g. test whether the credentials are correct
- the list of data streams that this plugin can collect


### Connection parameters

The parameters of your plugin are defined as class attributes of the connection class.
For example to add a `url` parameter of type `str` edit `MyPLuginConnection` as follow:

```python
class MyPluginConnection(Connection):
    url: str
```


## Add a new data stream

A data stream groups the logic for:
- collecting the raw data from the datasource
- extracting this raw data into a tool-specific model
- converting the tool model into an equivalent [DevLake domain model](https://devlake.apache.org/docs/next/DataModels/DevLakeDomainLayerSchema)


### Create a tool model

Create a `models.py` file.
Then create a class that modelizes the data your stream is going to collect.

```python
from pydevlake.model import ToolModel

class User(ToolModel, table=True):
    id: str = Field(primary_key=True)
    name: str
    email: str
```

Your tool model must declare at least one attribute as a primary key, like `id` in the example above.
It must inherit from `ToolModel`, which in turn inherit from `SQLModel`, the base class of an [ORM of the same name](https://sqlmodel.tiangolo.com/).
You can use `SQLModel` features like [declaring relationships with other models](https://sqlmodel.tiangolo.com/tutorial/relationship-attributes/).


### Create the stream class

Create a new file for your first stream in a `streams` directory.

```python
from pydevlake import Stream, DomainType
from pydevlake.domain_layer.crossdomain import User as DomainUser

from myplugin.models import User as ToolUser


class Users(Stream):
    tool_model = ToolUser
    domain_types = [DomainType.CROSS]

    def collect(self, state, context) -> Iterable[Tuple[object, dict]]:
        pass

    def convert(self, user: ToolUser, context) -> Iterable[DomainUser]:
        pass
```

This stream will collect raw user data, e.g. as parsed JSON objects, extract this raw data as your
tool-specific user model, then convert them into domain-layer user models.

The `tool_model` class attribute declares the tool model class that is extracted by this strem.
The `domain_types` class attribute is a list of domain types this stream is about.

The `collect` method takes a `state` dictionary and a context object and yields tuples of raw data and new state.
The last state that the plugin yielded for a given connection will be reused during the next collection.
The plugin can use this `state` to store information necessary to perform incremental collection of data.


The `convert` method takes a tool-specific user model and convert it into domain level user models.
Here the two models align quite well, the conversion is trivial:

```python
def convert(self, user: ToolUser, context: Context) -> Iterable[DomainUser]:
    yield DomainUser(
        id=user.id,
        name=user.name
        email=user.email
    )
```

### Create an API wrapper

Lets assume that your datasource is a REST API.
We can create the following class to define it.

```python
from pydevlake.api import API


class MyAPI(API):
    def __init__(self, url: str):
        self.url = url

    def users(self):
        return self.get(f'{self.url}/users')
```

By inheriting `API` you get access to facilities to wrap REST APIs.
Here the `users` method will return a `Response` object that contains the results of calling `GET` on `<url>/users`.

Now we can go back to our stream file and implement `collect`:

```python
from myplugin.api import MyAPI

...

    def collect(self, state, context) -> Iterable[Tuple[object, dict]]:
        api = MyAPI(context.connection.url)
        for user in api.users().json():
            yield user, state

...
```

If the API responds with a list of JSON object with properties matching your `User` model attributes, you're done!.
Indeed extraction has a default implementation that takes of this common case.
This is why it is important to make tool models that align with the data you collect.

If this is not the case, e.g. the attribute case mismatch, you can redefine the `extract` method:

```python
...

class Users(Stream):
    ...

    def extract(self, raw_data: dict) -> ToolModel:
        return ToolUser(
            id=raw_data["ID"],
            name=raw_data["Name"],
            email=raw_data["Email"]
        )

    ...
```


#### Request and response hook

For each request sent and response received by your API wrapper,
you can register hooks. Hooks allows you to implement
authentication, pagination, and generic API error handling.

For example, lets assume that we are dealing with an API that
require user to authenticate via a token set in a request header.

A request hook is declared in your API with a `@request_hook` decorator.

```python
...
class MyAPI(API):
    def __init__(self, url, token):
        self.url = url
        self.token = token

    ...
    @request_hook
    def authenticate(self, request):
        if self.token:
            request.headers['Token'] = self.token
    ...
```

Here the method `authenticate` is a hook that is run on each request.
Similarly you can declare response hooks with `@response_hook`.
Multiple hooks are executed in the order of their declaration.
The `API` base class declares some hooks that are executed first.


#### Pagination

One usage of a response hook is for handling paginated results.
A response hook can be used to wrap the `Response` object in a
`PagedResponse` object that support iteration and fetching other pages.
This response hook is actually defined in `API` base class and expect
your API wrapper to declare a `paginator` property.

You can subclass `Paginator` to provide API specific logic or reuse an
existing implementation such as `TokenPaginator`.
A token paginator assumes the API paginated responses are JSON object with one
property that is an array of items and another that contains the token to the next
page.

For example, the following paginator fetch items from the `'results'` attribute,
the next page from the `'nextPage'` attribute and will issue requests with a `page`
query parameter.

```
...
class MyAPI(API):
    ...
    paginator = TokenPaginator(
        items_attr='results',
        next_page_token_attr='nextPage',
        next_page_token_param='page'
    )
    ...
```

## Substreams

With REST APIs, you often need to fetch a stream of items, and then to collect additional
data for each of those items.

For example you might want to collect all `UserComments` written by each user collected via the `Users` stream.

To handle such cases, `UserComments` would inherit from `Substream` instead of `Stream`.
A substream needs to specify which is his parent stream. The `collect` method
of a substream will be called with each item collected from the parent stream.

```python
...
from pydevlake import Substream
from myplugin.streams.users import Users

class UserComments(Substream):
    parent_stream = Users # Must specify the parent stream
    ...
    def collect(self, state: dict, context, user: User):
        """
        This method will be called for each user collected from parent stream Users.
        """
        for json in MyPluginAPI(context.connection.token).user_comments(user.id):
            yield json, state
    ...
```


# Test the plugin standalone

To test your plugin manually, you can run your `main.py` file with different commands.
You can find all those commands with `--help` cli flag:

```console
poetry run myplugin/main.py --help
```

For testing, the interesting commands are `collect`/`extract`/`convert`.
Each takes a context and a stream name.
The context is a JSON object that must at least contain:
- a `db_url`, e.g. you can use `"sqlite+pysqlite:///:memory:"` for an in-memory DB
- a `connection` object containing the same attributes than your plugin connection type

Also, python plugins communicate with go side over an extra file descriptor 3, so you should
redirect to stdout when testing your plugin.

```
console
CTX='{"db_url":"sqlite+pysqlite:///:memory:", "connection": {...your connection attrs here...}}'
poetry run myplugin/main.py $CTX users 3>&1
```


# Test the plugin with DevLake

To test your plugin together with DevLake, you first need to create a connection for your plugin and get its id.
One easy way to do that is to run DevLake with `make dev` and then to create the connection with a POST
request to your plugin connection API:

```console
curl -X 'POST' \
  'http://localhost:8080/plugins/myplugin/connections' \
  -d '{...connection JSON object...}'
```

You should get the created connection with his id (which is 1 for the first created connection) in the response.

Now that a connection for your plugin exists in DevLake database, we can try to run your plugin using `backend/server/services/remote/run/run.go` script:

```console
cd backend
go run server/services/remote/run/run.go  -c 1 -p python/plugins/myplugin/myplugin/main.py
```

This script takes a connection id (`-c` flag) and the path to your plugin `main.py` file (`-p` flag).
You can also send options as a JSON object (`-o` flag).

# Automated tests
Make sure you have unit-tests written for your plugin code. The test files should end with `_test.py`, and are discovered and
executed by the `run_tests.sh` script by the CICD automation. The test files should be placed inside the plugin project directory.


# Debugging Python plugins
You need to have a PyCharm debugger installed to debug the Python code. This capability is controlled by the environment
variable `ENABLE_PYTHON_DEBUGGER` which defaults to `false`. Set to `true` to enable it. You will further have to set
the environment variables `PYTHON_DEBUG_HOST` (The hostname/IP on which your debugger is running relative to the environment
in which the plugin is running) and `PYTHON_DEBUG_PORT` (The corresponding debugger port). The variables should be set in the
Go integration tests written in `backend/test/integration/remote` or Docker container/server env configuration. Once done,
set breakpoints in the Python plugin code in your IDE, turn on the debugger in it, and those breakpoints should get hit.