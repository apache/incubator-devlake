## So you want to Build a New Plugin...

...the good news is, it's easy!

### Summary

To build a new plugin you will need a few things. You should choose an API that you'd like to see data from. Think about the metrics you would like to see first, and then look for data that can support those metrics.

Then you will want to build a collector to gather data. You will need to do some reading of the API documentation to figure out what metrics you will want to see at the end in your Grafana dashboard (configuring Grafana is the final step).  

We're working with Node, so you will want your collector to use some sort of HTTP requests to gather data from the api. A package like axios can be used here, or there may be a Node package you can download and use. Then you will want to store that raw data in a DB. If you are storing it in a document DB like Mongo, you can store it as is. If you are using a relational DB like Postgres, you will need a schema to be defined according to the fields you get from your API requests. 

Once you are able to store the raw data from your queries, you will want to enrich that data to:
a) Add fields you don't currently have
b) Compute fields you might want for metrics
c) Eliminate fields you don't need

To build the enricher, you will want to query your DB to find the raw data you've stored. Then you may perform any actions you want in terms of calculations to enrich data. You will need to define a new model for your data, and migrations for your enrichment DB. Once you have that set up, you can store your raw data in its new format in your enriched data DB. This is recommended to be a relational DB.

It is good to build a collector and an enricher together, because knowledge of the API will help with both.

Let's walk through a short example.

Let's say you want to see data from the Movie Database.

