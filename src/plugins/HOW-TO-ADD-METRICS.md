## So you want to add a new metric...

...you've never had a better idea than that!

You love getting data from a source like GitLab, but perhaps you're not getting the 
the metric you want. Good news! You can get any metric you like by contributing a little
bit of code.

Let's say we are getting data on GitLab for projects, and commits, but you want to see
data on issues.

To add this metric, all you need to do is visit the documentation, and find the API routes
that support your metrics. Don't forget authentication and pagination. It will take some research, but you'll get it!

As it turns out, GitLab has a lot of docs on their Issues API: (https://docs.gitlab.com/ee/api/issues.html)

Let's say you want to see number of comments per issue. After a little reading you will notice
that comments are found through the Notes API rather than the Issues API. This is why a little
reading goes a long way! (https://docs.gitlab.com/ee/api/issues.html#comments-on-issues)

So you've done some research, and you've pinpointed the endpoint you'd like to call:

GET /projects/:id/issues/:issue_iid/notes

This one looks correct. But wait... you need a project id AND an issue iid! This means you'll
have to provide a project ID somehow, and then get all the issues from the project, then loop
through the issues to get all the notes! 

Once you've done all that, you need to store the data in the DB. Next step is enrichment.



