# Pokemon Pond

## Plan

1. Get the plugin boilerplate in
1. Write the pokemon model
   - what attributes do we want to store in the model?
1. Write pokemon collect task
   - Do we want to query for every pokemon every time we run this? We could query for what we have locally and only get what's missing. e.g. in the case where there's new pokemon. Or we can
1. Write item collect task
   - this should query for the pokemon locally, and then query the pokemon API for their items
   - plan for now is to update the pokemon model with a cost attribute
1. Write some tests
