# So you want to add a new metric...
...you've never had a better idea than that!

You love getting data from a source like GitLab, but perhaps you're not getting the the metric you want. Good news! You can get any metric you like by contributing a little bit of code.

Each plugin has an Enricher. Just navigate into the plugin folder to find the enrichment folder. There, you should find all the various enrichment methods. The job of the Enricher is to grab the data that was collected by the Collector, augment it, and add it to a new DB in accordance with a new schema (which you will have to create).

- [ ] Details of how to add metrics with new go setup
