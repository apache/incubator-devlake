ClickUp
=======

Implementation
--------------

- Spaces are mapped as "scopes".
- Folders are not *currently* supported.

Current Requirements
--------------------

### Bug and Incident Tracking

- Add a custom `DropDown` field called `type` in ClickUp
- Add at least options for `Bug` and `Incident`

Manual Testing
--------------

The backend plugin can be run directly for development purposes:

```
go run ./plugins/clickup/clickup.go -c<connection id> -s<scope id>
```

Where connection ID is the ID of a `_tools_clickup_connections` (configured
via. the GUI) and `scope id` is the space ID. For example:

```
./plugins/clickup/clickup.go -c1 -s90060207111
```
