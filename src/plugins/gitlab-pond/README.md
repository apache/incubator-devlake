# Gitlab Pond

## Metrics

Metric Name | Description
:------------ | :-------------
Number of Merge Requests | Number of Merge Requests
Merge Request Pass Rate	| Ratio of merge requests to merged
Code Reviewer Count	| Number of Merge Reviewers
Code Review Time | Time from the first Merge Review comment until merged
Commit Author Count | Number of Contributors
Commit Count | Number of Commits
Accumulated New Lines of Code | Accumulated Number of New Lines
Accumulated Deleted Lines of Code | Accumulated Number of Removed Lines

## Finding Project Id
To get the project id for a specific Gitlab repository:
- Visit the repository page on gitlab
- Find the project id just below the title

  ![Screen Shot 2021-08-06 at 4 32 53 PM](https://user-images.githubusercontent.com/3789273/128568416-a47b2763-51d8-4a6a-8a8b-396512bffb03.png)

> Use this project id in your requests, to collect data from this project

## Create a Gitlab API Token

1. When logged into Gitlab visit `https://gitlab.com/-/profile/personal_access_tokens`
2. Give the token any name, no expiration date and all scopes (excluding write access)

    ![Screen Shot 2021-08-06 at 4 44 01 PM](https://user-images.githubusercontent.com/3789273/128569148-96f50d4e-5b3b-4110-af69-a68f8d64350a.png)

3. Click the **Create Personal Access Token** button
4. Copy the token into the `lake` repo config
