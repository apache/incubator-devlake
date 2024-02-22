# DevLake Jupyter Playground

[DevLake](https://devlake.apache.org/) offers an abundance of data for exploration.
This playground contains a basic set-up to interact with the data using [Jupyter Notebooks](https://jupyter.org/) and [Pandas](https://pandas.pydata.org/).


# How to play

## Prerequisites
- [Python >= 3.11](https://www.python.org/downloads/)
- [Poetry](https://python-poetry.org/docs/#installation)
- Access to a DevLake database
- A place to run a Jupyter Notebook (e.g. [VS Code](https://code.visualstudio.com/))


## Usage
1. Have a local clone of this repository.
2. Run `poetry install` in the `playground` directory.
3. Open an example Jupyter Notebook from the `notebooks` directory in your preferred Jupyter Notebook tool.
4. Make sure the notebook uses the virtual environment created by poetry.
5. Configure your database URL in the notebook code.
6. Run the notebook.
7. Start exploring the data in your own notebooks!


## Create your own Jupyter Notebook

A good starting point for creating a new notebook is `template.ipynb`.
It contains the basic steps you need to go from query to output.

To define a query, use the [Domain Layer Schema](https://devlake.apache.org/docs/DataModels/DevLakeDomainLayerSchema#schema-diagram) to get an overview of the available tables and fields.

Use [Pandas](https://pandas.pydata.org/) api to organize, transform, and analyze the query results.


## Predefined notebooks and utilities

A notebook might offer a valuable perspective on the data not available within the capabilities of a Grafana dashboard.
In this case, it's worthwhile to contribute this notebook to the community as a predefined notebook, e.g., `process_analysis.ipynb` (it depends on [graphviz](https://graphviz.org/) for its visualization.)

The same goes for utility methods with, for example, predefined Pandas data transformations offering an interesting view on the data.

Please check the [contributing guidelines](https://github.com/apache/incubator-devlake/blob/main/README.md#-how-to-contribute).
