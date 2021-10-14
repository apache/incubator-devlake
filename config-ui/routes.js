const routes = require('next-routes')

module.exports = routes()
//? Name is URL route
//? Page is page (file) name in /pages
//? Pattern is for dynamic routes eg. '/user/:id'

// Main Setup
.add({name: 'configuration', page: '/'})
.add({name: 'triggers', page: '/triggers'})

// Plugins
.add({name: 'jira', page: '/plugins/jira'})
.add({name: 'gitlab', page: '/plugins/gitlab'})
.add({name: 'jenkins', page: '/plugins/jenkins'})
