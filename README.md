# Dev Lake

## Quick start

### Required
- [Install Golang](https://golang.org/doc/install)

### Developer Setup

```shell
git clone https://github.com/merico-dev/lake.git
cd lake
make get
cp .env.example .env
make build
make compose

```
While docker is running, in a new terminal:
```
cd lake
./lake
```

Then you can make a POST request:
```
curl --location --request POST 'localhost:8080/source' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Plugin": "Jira",
    "Options": {}
    
}'
```

### Makefile

This is like the package.json file that runs all our commands

1. You can install make 
Ubuntu: `sudo apt-get install build-essential`
Windows: `http://gnuwin32.sourceforge.net/packages/make.htm`
Mac: Comes pre installed
2. Then you can run make commands like this:
`make hello`

### How to make a commit

We use https://github.com/lintingzhen/commitizen-go to author our commits. 

Then you can run:
`make commit`

### How to run the tests

You can see a sample test in /test/example
You can run the tests with `make test`
## Dev Lake


# Plugin System

Directory strucutre
```sh
plugins/
├── core
│   └── plugin.go           # plugin interface
├── jira                    # plugin directory
│   ├── jira.go             # plugin entry point source file which implemented plugin interface
│   └── jira.so             # compiled shared library
├── plugins.go              # plugins management
└── plugins_test.go         # unit testing

```

## How to write a plugin

1. Create a subdirectory within `plugins/` directory, i.e. `plugins/jira`.
2. Create a entry point file with in subdirectory from previous step, i.e. `plugins/jira.go`.
3. Declare a pseudo type in entry point file and implement `Plugin` interface defined in `plugins/core/plugin.go`.
4. Declare a `PluginEntry` variable in entry point file for Plugin Manager to search and load.
5. Compile plugin with `go build -buildmode=plugin -o jira/jira.so jira/jira.go`.

## How to use Plugins

1. Import `github.com/merico-dev/lake/plugins` module.
2. Call `LoadPlugin` with plugins directory path, i.e. `plugins/`, it will load all first `.so` file of each
   subdirectory as plugins, each plugin has the same name to the subdirectory.
3. Call `go RunPlugin` with plugin `name` along with its `options`, plus a channel named `progress`.


## How to verify

1. Change working directory to `plugins`: `cd plugins/`
2. Compile plugin with `go build -buildmode=plugin -o jira/jira.so jira/jira.go`
3. Run test with `go test . -v`, you should see something like:
   ```
   === RUN   TestPluginsLoading
   start runing plugin jira
   start jira plugin execution
   running plugin jira, progress: 10
   running plugin jira, progress: 50
   end jira plugin execution
   running plugin jira, progress: 100
   end running plugin jira
   --- PASS: TestPluginsLoading (3.00s)
   PASS
   ok      github.com/merico-dev/lake/plugins      3.008s
   ```
4. In order to debug using VSCode, plugin must be compiled with option `-gcflags="all=-N -l"` for step 2, do so will
   break test command from step 3, unless same option was added.
