# Dbt

<div align="center">

| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

</div>

<br>

## Summary

dbt (data build tool) enables analytics engineers to transform data in their warehouses by simply writing select statements. dbt handles turning these select statements into tables and views.
dbt does the T in ELT (Extract, Load, Transform) processes – it doesn’t extract or load data, but it’s extremely good at transforming data that’s already loaded into your warehouse.

## User setup<a id="user-setup"></a>
- If you only plan to run the product locally, this is the **ONLY** section you should need.

#### Required Packages to Install<a id="user-setup-requirements"></a>
- [python3.7+](https://www.python.org/downloads/)
- [dbt-mysql](https://pypi.org/project/dbt-mysql/#configuring-your-profile)

1. pip install dbt-mysql
2. dbt init demoapp (demoapp is project name) 
3. create your SQL transformations and data models if you need

## Convert Data By Dbt

please use the Raw JSON API to manually initiate a run using **cURL** or graphical API tool such as **Postman**. `POST` the following request to the DevLake API Endpoint.

```json
[
  [
    {
      "plugin": "dbt",
      "options": {
          "projectPath": "/Users/abeizn/demoapp",
          "projectName": "demoapp",
          "projectTarget": "dev",
          "selectedModels": ["my_first_dbt_model","my_second_dbt_model"],
          "projectVars": {
            "demokey1": "demovalue1",
            "demokey2": "demovalue2"
        }
      }
    }
  ]
]
```

- `projectPath`: the absolute path of the dbt project. (required)
- `projectName`: the name of the dbt project. (required)
- `projectTarget`: this is the default target your dbt project will use. (optional)
- `selectedModels`: a model is a select statement. Models are defined in .sql files, and typically in your models directory. (required)
And selectedModels accepts one or more arguments. Each argument can be one of:
1. a package name #runs all models in your project, example: example
2. a model name   # runs a specific model, example: my_fisrt_dbt_model
3. a fully-qualified path to a directory of models # you need update the value of model-paths in dbt_project.yml. 
- `vars`: dbt provides a mechanism variables to provide data to models for compilation. (optional) 
example: select * from events where event_type = '{{ var("event_type") }}' this sql in your model, you need set parameters "vars": "{event_type: real_value}"

### Resources:
- Learn more about dbt [in the docs](https://docs.getdbt.com/docs/introduction)
- Check out [Discourse](https://discourse.getdbt.com/) for commonly asked questions and answers
- Join the [chat](https://community.getdbt.com/) on Slack for live discussions and support
- Find [dbt events](https://events.getdbt.com) near you
- Check out [the blog](https://blog.getdbt.com/) for the latest news on dbt's development and best practices

