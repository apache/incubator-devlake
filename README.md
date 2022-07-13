<div align="center">
<br/>
<img src="img/logo.svg" width="120px">
<br/>

# Apache DevLake(Incubating)

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
![badge](https://github.com/apache/incubator-devlake/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/incubator-devlake)](https://goreportcard.com/report/github.com/apache/incubator-devlake)
[![Slack](https://img.shields.io/badge/slack-join_chat-success.svg?logo=slack)](https://join.slack.com/t/devlake-io/shared_invite/zt-17b6vuvps-x98pqseoUagM7EAmKC82xQ)
</div>
<br>
<div align="left">

### What is Apache DevLake?
Apache DevLake is an open-source dev data platform that ingests, analyzes, and visualizes the fragmented data from DevOps tools to distill insights for engineering productivity.

Apache DevLake is designed for developer teams looking to make better sense of their development process and to bring a more data-driven approach to their own practices. You can ask Apache DevLake many questions regarding your development process. Just connect and query.

### Demo
See [demo](https://grafana-lake.demo.devlake.io/d/0Rjxknc7z/demo-homepage?orgId=1). The data in the demo comes from this repo.


<br/>

<div align="left">
<img src="https://user-images.githubusercontent.com/14050754/145056261-ceaf7044-f5c5-420f-80ca-54e56eb8e2a7.png" width="100%" alt="User Flow" style="border-radius:15px;"/>
<p align="center">User Flow</p>

<br/>


## What can be accomplished with Apache DevLake?
1. Collect DevOps data across the entire Software Development Life Cycle (SDLC) and connect the siloed data with a standard [data model](https://devlake.apache.org/docs/DataModels/DevLakeDomainLayerSchema).
2. Provide out-of-the-box engineering [metrics](https://devlake.apache.org/docs/EngineeringMetrics) to be visualized in a sereis of dashboards.
3. Allow a flexible [framework](https://devlake.apache.org/docs/Overview/Architecture) for data collection ad ETL to support customizable data analysis.


## Supported Data Sources

| Data Source                                                | Domain                                                     | Versions                             |
| ---------------------------------------------------------- | ---------------------------------------------------------- | ------------------------------------ |
| [Feishu](https://devlake.apache.org/docs/Plugins/feishu)   | Documentation                                              | Cloud                                |
| [GitHub](https://devlake.apache.org/docs/Plugins/github)   | Source Code Management, Code Review, Issue/Task Management | Cloud                                |
| [Gitlab](https://devlake.apache.org/docs/Plugins/gitlab)   | Source Code Management, Code Review, Issue/Task Management | Cloud, Community Edition 13.x+       |
| [Jenkins](https://devlake.apache.org/docs/Plugins/jenkins) | CI/CD                                                      | 2.263.x+                             |
| [Jira](https://devlake.apache.org/docs/Plugins/jira)       | Issue/Task Management                                      | Cloud, Server 8.x+, Data Center 8.x+ |
| TAPD                                                       | Issue/Task Management                                      | Cloud                                |


## Quick Start
- [Deploy Locally](https://devlake.apache.org/docs/QuickStart/LocalSetup)
- [Deploy to Kubernetes](https://devlake.apache.org/docs/QuickStart/KubernetesSetup)
- [Deploy in Temporal Mode](https://devlake.apache.org/docs/UserManuals/TemporalSetup)
- [Deploy in Developer Mode](https://devlake.apache.org/docs/DeveloperManuals/DeveloperSetup)


## Project Roadmap
- <a href="https://devlake.apache.org/docs/Overview/Roadmap" target="_blank">Roadmap 2022</a>: Detailed project roadmaps for 2022.
- <a href="https://devlake.apache.org/docs/EngineeringMetrics" target="_blank">Supported engineering metrics</a>: provide rich perspectives to observe and analyze SDLC.


## How to Contribute
This section lists all the documents to help you contribute to the repo.

- [Architecture](https://devlake.apache.org/docs/Overview/Architecture): Architecture of Apache DevLake
- [Data Model](https://devlake.apache.org/docs/DataModels/DevLakeDomainLayerSchema): Domain Layer Schema
- [Add a Plugin](/plugins/README.md): Guide to add a plugin
- [Add Metrics](/plugins/HOW-TO-ADD-METRICS.md): Guide to add metrics in a plugin
- [Contribution Guidelines](https://devlake.apache.org/community): Start from here if you want to make contribution


## Community

- <a href="https://join.slack.com/t/devlake-io/shared_invite/zt-18uayb6ut-cHOjiYcBwERQ8VVPZ9cQQw" target="_blank">Slack</a>: Message us on Slack
- <a href="https://github.com/apache/incubator-devlake/wiki/FAQ" target="_blank">FAQ</a>: Frequently Asked Questions
- Wechat Community:<br>
  ![](img/wechat_community_barcode.png)


## License<a id="license"></a>

This project is licensed under Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
