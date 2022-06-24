# Apache DevLake Cli Tool -- Code Generator

## How to use?

Just run by `go run`:
```bash
go run generator/main.go [command]
```

Help also integrate in it:
```bash
go run generator/main.go help
```

## Plugin Related

* [create-collector](./docs/generator_create-collector.md)	 - Create a new collector
* [create-extractor](./docs/generator_create-extractor.md)	 - Create a new extractor
* [create-plugin](./docs/generator_create-plugin.md)

Usage Gif:
![usage](https://user-images.githubusercontent.com/3294100/175464884-1dce09b0-fade-4c26-9a1b-b535d9651bc1.gif)

## Migration Related

* [init-migration](./docs/generator_init-migration.md)	     - Init migration for plugin
* [create-migration](./docs/generator_create-migration.md)	 - Create a new migration

Usage Gif:
![usage](https://user-images.githubusercontent.com/3294100/175537358-083809ce-9862-41f1-86e9-41448a44eaae.gif)

## Others

* [completion](./docs/generator_completion.md)               - Generate the autocompletion script for the specified shell
* [generator-doc](./docs/generator_generator-doc.md)         - generate document for generator
* [global options](./docs/generator.md)
