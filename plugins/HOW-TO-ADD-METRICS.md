# So you want to add a new metric...
...you've never had a better idea than that!

You love getting data from a source like GitLab, but perhaps you're not getting the the metric you want. Good news! You can get any metric you like by contributing a little bit of code.

1. Decide if you are going to need new data for your metric, or if you just need to enrich the data that already exists.

2.  If you need new data, then you will need to:
  a.  Create a model for the data you'd like to capture. 
  b.  Add your new model to the plugin's init file. This will allow GORM to 
      auto-migrate the definition to create a new DB table.
  c.  Create a collector in the "tasks" folder. You will need to make API calls to gather and save your new data. 
      You may need to do some research to figure out what fields are returned from the API to specifically capture
      everything you would like to calculate your metric.
  d.  Add any additional "enrichment" calculations on top of the data that you are fetching.
  e.  Add your "collection" method to the plugin's execute function in the main package.
  f.  Start the project, trigger an API request that triggers your plugin, and watch the data flow in!
  g.  Congrats, you are done!

3.  If you do not need new data then all you need to do is add any additional "enrichment" calculations on top of the data that 
    we are already are fetching!

Thanks for your contribution!

-- The Dev Lake Team

