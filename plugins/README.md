# Plugin System

Directory strucutre
```sh
plugins/
├── core/
│   └── plugin.go           # plugin interface
├── jira/                   # plugin directory
│   ├── models/             # entity path (or others name you like)
│   ├── jira.go             # plugin entry point source file which implemented plugin interface
│   └── jira.so             # compiled shared library
├── plugins.go              # plugins management
└── plugins_test.go         # unit testing

```

## How to start your plugin author journey

1. Create a subdirectory within `plugins/` directory, i.e. `plugins/jira`.
2. Create a entry point file with in subdirectory from previous step, i.e. `plugins/jira.go`.
3. Declare a pseudo type in entry point file and implement `Plugin` interface defined in `plugins/core/plugin.go`.
4. Declare a `PluginEntry` variable in entry point file for Plugin Manager to search and load.
5. Compile plugin with `go build -buildmode=plugin -o jira/jira.so jira/jira.go`.

### What do Plugin interface methods mean

1. `Description() string` A method that describes the purpose of the plugin
2. `Init()` Be called when plugin load
3. `Execute(options map[string]interface{}, progress chan<- float32)` Be called when plugin start to fetch/process data

## How to use Plugins

1. Import `github.com/merico-dev/lake/plugins` module.
2. Call `LoadPlugin` with plugins directory path, i.e. `plugins/`, it will load all first `.so` file of each
   subdirectory as plugins, each plugin has the same name to the subdirectory.
3. Call `go RunPlugin` with plugin `name` along with its `options`, plus a channel named `progress`.


## How to verify

1. Change working directory to `plugins`: `cd plugins/`
2. Compile the plugin with `go build -buildmode=plugin -o jira/jira.so jira/*.go`
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