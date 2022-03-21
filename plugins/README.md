# So you want to Build a New Plugin...

...the good news is, it's easy!


## Basic Interface

```golang
type YourPlugin string

func (plugin YourPlugin) Description() string {
	return "To collect and enrich data from YourPlugin"
}

func (plugin YourPlugin) Execute(options map[string]interface{}, progress chan<- float32) {
	logger.Print("Starting YourPlugin execution...")

  // Check fields that are needed in options.
	projectId, ok := options["projectId"]
	if !ok {
		logger.Print("projectId is required for YourPlugin execution")
		return
	}

  // Start collecting stuff.
  if err := tasks.CollectProject(projectId); err != nil {
		logger.Error("Could not collect projects: ", err)
		return
	}
  // Handle error.
  if err != nil {
    logger.Error(err)
  }

  // Export a variable named PluginEntry for Framework to search and load
  var PluginEntry YourPlugin //nolint
}
```

## Summary

  To build a new plugin you will need a few things. You should choose an API that you'd like to see data from. Think about the metrics you would like to see first, and then look for data that can support those metrics.

## Collection

  Then you will want to write a collection to gather data. You will need to do some reading of the API documentation to figure out what metrics you will want to see at the end in your Grafana dashboard (configuring Grafana is the final step).

## Build a Fetcher to make Requests

The plugins/core folder contains an api client that you can implement within your own plugin. It has useful methods like Get(). 
Each API handles pagination differently, so you will likely need to implement a "fetch with pagination" method. One way to do
this is to use the "ants" package as a way to manage tasks concurrently: https://github.com/panjf2000/ants

Your colection methods may look something like this:

```golang
func Collect() error {
	pluginApiClient := CreateApiClient()

	return pluginApiClient.FetchWithPagination("<your_api_url>",
		func(res *http.Response) error {
			pluginApiResponse := &ApiResponse{}
      // You must unmarshal the response from the api to use the results.
			err := helper.UnmarshalResponse(res, pluginApiResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
      // loop through the results and save them to the database.
			for _, value := range *pluginApiResponse {
				pluginModel := &models.pluginModel{
					pluginId:       value.pluginId,
					Title:          value.Title,
					Message:        value.Message,
				}

				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&pluginModel).Error

				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
			}

			return nil
		})
}
```

Note the use of "upsert". This is useful for only saving modified records.

## Enrichment

  Once you have collected data from the API, you may want to enrich that data by:

  - Add fields you don't currently have
  - Compute fields you might want for metrics
  - Eliminate fields you don't need

## You're Done!

Congratulations! You have created your first plugin! 🎖 
