## So you want to add a new metric...

...you've never had a better idea than that!

You love getting data from a source like GitLab, but perhaps you're not getting the
the metric you want. Good news! You can get any metric you like by contributing a little
bit of code.

Each plugin has an Enricher. Just navigate into the plugin folder to find the enrichment folder. There, you should find all the various enrichment methods. The job of the Enricher is to grab the data that was collected by the Collector, augment it, and add it to a new DB in accordance with a new schema (which you will have to create).

An the index file of an enricher contains something like this:

```js
async function enrich (rawDb, enrichedDb, { thingId }) {
  const args = { rawDb, enrichedDb, thingId: Number(projectId) }
  await thingNumberOne.enrich(args)
  await thingNumberTwo.enrich(args)
  await thingNumberThree.enrich(args)
}

```

It gathers all the enrich functions into one, so you can keep your implementations separate. For example:

```js
async function enrichSomeThing (rawDb, enrichedDb, id) {
  const collectionOfRawThings = await collector.getCollection(rawDb)
  const rawThing = await collectionOfRawThings.findOne({ id: id })
  const computedField = computeThings(rawThing.somethingToCompute, rawThing.otherThingToCompute)
  const enrichedThing = {
    name: rawThing.name,
    id: rawThing.id,
    aFieldYouWant: rawThing.a_field_you_want,
    computedField: computedField
  }
  await enrichedDb.ThingModel.upsert(enrichedThing)
}
```

In order for the enriched DB to have a ThingModel for you to upsert into, you will need to write a Model file in the [model directory](../../db/postgres) and a migration script for Sequelize in the [migrations directory](../../db/migrations)

<br>

---

## Other Docs

Section | Description | Link
:------------ | :------------- | :-------------
Requirements | Underlying software used | [Link](../../README.md#requirements)
User Setup | Quick and easy setup | [Link](../../README.md#user-setup)
Developer Setup | Steps to get up and running | [Link](../../README.md#developer-setup)
Plugins | Links to specific plugin usage & details | [Link](../../README.md#plugins)
Build a Plugin | Details on how to make your own | [Link](README.md)
Grafana | How to visualize the data | [Link](../../docs/GRAFANA.md)
Contributing | How to contribute to this repo | [Link](../../CONTRIBUTING.md)
