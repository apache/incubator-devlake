
# For `make e2e-test` to run properly, the following steps must be taken:

1. `python3.9` is required by the time of this document. 
   - Try `deadsnakes` if you are using Ubuntu 22.04 or above, the `python3.9-dev` is required.
   - Use `virtualenv` if you are having multiple python versions. `virtualenv -p python3.9 path/to/venv` and `source path/to/venv/bin/activate.sh` should do the trick
2. [poetry](https://python-poetry.org/) is required. 
   - run `cd backend/python/pydevlake && poetry install`
   - run `cd backend/python/plugins/azuredevops && poetry install`
3. sqlalchemy won't work with `localhost` in the database connection string, use `127.0.0.1` instead
