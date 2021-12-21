# Help with Debugging

Use this document to help setup Go debugging for this project, including plugins code.

[Delve](https://github.com/go-delve/delve) is the primary tool for debugging code in Go. See the project documentation for general information about the debugger.

## Building Plugins

In order to successfully step into plugins code, the `-gcflags="all=-N -l"` option must be passed to the command that builds the plugins. This will ensure that necessary debug symbols are included in the output binaries. Run `make build-plugin-debug` to build all plugins with this option.

## Open Bugs

There are currently some open issues related to debugging that should be noted.

* [Cannot debug Go plugins on Mac](https://github.com/go-delve/delve/issues/1628)

This issue makes it impossible to debug plugins on Mac only. It will still be possible to debug non-plugin code in the usual way.

* [Breakpoints in plugins are not working if set before plugin is loaded](https://github.com/go-delve/delve/issues/1653)

This issue is a nuisance that affects all OS, but can be worked around. When starting the debugger, any previously set breakpoints in plugin code will be disabled. In VS Code, a message appears stating e.g. cannot find file plugins/example.go.

To get around this, set a breakpoint on the line immediately after plugins are loaded. Currently this is done in `main.go`, in the very first line of `main`. With the debugger paused on the next line `if err != nil {`, reset any breakpoints and they should work as expected. In VS Code, this can be done for all breakpoints at once with the command `Toggle Activate Breakpoints`.

```go
func main() {
	err := plugins.LoadPlugins(plugins.PluginDir())
	if err != nil {
		panic(err)
	}
	api.CreateApiService()
	println("Hello, lake")
}
```

## Configure in VS Code

### Install VS Code Go Extension

Install the [VS Code Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go), which provides many tools for Go development, including debugging.

Once installed, update the Go tools to ensure that you have access to all the latest debugging features:

* Open the command pallet (*Ctrl+Shift+P*)
* Type and select *Go: Install/Update Tools*
* Select all tools, and click *OK*

### Setup a Build Task

Create a file `.vscode/tasks.json` that includes a task that runs the command to build plugins for debugging. In this example, the task is named `build-plugin`.

```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "build-plugin",
            "type": "shell",
            "command": "make build-plugin-debug"
        }
    ]
}
```

### Setup a Launch Configuration

Create a file `.vscode/launch.json` that instructs VS Code to launch the Go debugger on the file `main.go`, after running the `build-plugin` task. This will ensure that plugins are always built before debugging starts. If not working on plugin code, the pre-launch task can be disabled to save time.

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "main.go",
            "preLaunchTask": "build-plugin"
        }
    ]
}
```
