## E2E Tests

## Why?

To ensure data integrity of the application, we need to make sure that the result
data matches what we expect to get from accessing a real API as if a real user is
using it. Automated tests allow us to do this in a very convenient, low cost way 
that is easily repeatable.

## How it works

1. Automatically or Manually trigger all collection / enrichment / conversion tasks
2. Tests access all key data models from our DB to determine if the expected number 
of rows were collected and processed or not.
