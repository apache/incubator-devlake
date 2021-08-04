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

## Step by Step

Let's walk through a short example.

1. Choose an API you're interested in.

Let's say you want to see data from the Movie Database. 

2. Choose some metrics you would like to see about movies.

You would like to know how movie production has increased over the last 50 years.
Let's start with how many movies have been produced each year as a metric.

3. Look for the data to support this from the Movie DB API 

Lets assume this data is available through an endpoint like (this may not be the case):

GET /movies

4. Get an API key for authentication.

You will need an API key to access the data. 

5. Understand Rate Limits and Pagination

There are limits to accessing API data. You will need to work around these limits.

6. Build the collector to fetch data from the API and store the raw data in a DB.

7. Build the enricher to find data from the raw data DB, perform enrichment, and store it in
   a relational DB with a new schema.

8. Build a graph you'd like to see within Grafana
   You may now use your enriched DB as a data source within Grafana. Once you've connected
   the data source, you can write SQL within Grafana to get the data and present it in a
   variety of ways!

9. Congrats on building your first plugin!



