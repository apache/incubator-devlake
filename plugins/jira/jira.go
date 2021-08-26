package main // must be main for plugin entry point

import "github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
type Jira string

func (plugin Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (jira Jira) Execute(options map[string]interface{}, progress chan<- float32) {
	// TODO: go get the first page of data
	// - User sets jira auth token, email and boardId in the .env file
	// TODO: add a ticket for disuccion about source and task in DB for clarity.
	// - Question: Why? User needs send a post request to /source ?
	// - On app load, run plugins
	// - Run the jira plugin execute method
	// - Read the token, board id and email from the .env file
	// - Base 64 encode it with their email:token
	// - Call jira api for board
	// - save response in DB
	// - call jira api for issues based on board id
	// - save all issues in the db
	// TODO: add a new ticket to handle paging
	// TODO: add a ticket to calculate lead time based on issues
	logger.Info("Execute called", options)

	close(progress)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Jira //nolint
