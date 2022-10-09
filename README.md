<div align="center">
<br/>
<img src="img/logo.svg" width="120px">
<br/>

# Apache DevLake(Incubating)

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![unit-test](https://github.com/apache/incubator-devlake/actions/workflows/test.yml/badge.svg)](https://github.com/apache/incubator-devlake/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/incubator-devlake)](https://goreportcard.com/report/github.com/apache/incubator-devlake)
[![Slack](https://img.shields.io/badge/slack-join_chat-success.svg?logo=slack)](https://join.slack.com/t/devlake-io/shared_invite/zt-17b6vuvps-x98pqseoUagM7EAmKC82xQ)

</div>
<br>
<div align="left">

## 🤔 What is Apache DevLake?

[Apache DevLake](https://devlake.apache.org) is an open-source dev data platform that ingests, analyzes, and visualizes the fragmented data from DevOps tools to extract insights for engineering excellence, developer experience, and community growth.

Apache DevLake is designed for developer teams looking to make better sense of their development process and to bring a more data-driven approach to their own practices. You can ask Apache DevLake many questions regarding your development process. Just connect and query.

## 🎯 What can be accomplished with Apache DevLake?

1. Your Dev Data lives in many silos and tools. DevLake brings them all together to give you a complete view of your Software Development Life Cycle (SDLC).
2. From DORA to scrum retros, DevLake implements metrics effortlessly with prebuilt dashboards supporting common frameworks and goals.
3. DevLake fits teams of all shapes and sizes, and can be readily extended to support new data sources, metrics, and dashboards, with a flexible framework for data collection and transformation.

## 👉 Live Demo
[DORA Dashboard](https://grafana-lake.demo.devlake.io/grafana/goto/r1EVQH4Vk?orgId=1)

[All Dashbaords](https://devlake.apache.org/docs/LiveDemo)

## 💪 Supported Data Sources

| Data Source  | Domain(s)                                           |Supoorted Versions                    | Plugins |
|--------------|-----------------------------------------------------|--------------------------------------|-------- |
| GitHub       | Source Code Management, Code Review, Issue Tracking | Cloud                                |`github`, `gitextractor`, `refdiff` |
| GitLab       | Source Code Management, Code Review, Issue Tracking | Cloud, Community Edition 13.x+       |`gitlab`, `gitextractor`, `refdiff` |
| BitBucket    | Source Code Management, Code Review, Issue Tracking | Cloud                                |`bitbucket`, `gitextractor`, `refdiff` |
| Gitee        | Source Code Management, Code Review, Issue Tracking | Cloud                                |`gitee`, `gitextractor`, `refdiff` |
| Jira         | Issue Tracking                                      | Cloud, Server 8.x+, Data Center 8.x+ |`jira` |
| TAPD         | Issue Tracking                                      | Cloud                                | `tapd` |
| Jenkins      | CI/CD                                               | 2.263.x+                             |`jenkins` |
| GitHub Action| CI/CD                                               | Cloud                                |`github`|
| GitLab CI    | CI/CD                                               | Cloud, Community Edition 13.x+       |`gitlab`|
| Feishu       | Calendar                                            | Cloud                                |`feishu` |
| AE           | Source Code Management                              |                                      | `ae`    |

## 🚀 Getting Started

- [Install via Docker Compose](https://devlake.apache.org/docs/GettingStarted/DockerComposeSetup)
- [Install via Kubernetes](https://devlake.apache.org/docs/GettingStarted/KubernetesSetup)
- [Install via Helm ](https://devlake.apache.org/docs/GettingStarted/HelmSetup)
- [Install in Temporal Mode](https://devlake.apache.org/docs/GettingStarted/TemporalSetup)
- [Install in Developer Mode](https://devlake.apache.org/docs/DeveloperManuals/DeveloperSetup)

## 🤓 How do I use DevLake?

### 1. Set up DevLake

You can set up Apache DevLake by following our step-by-step instructions for [Install via Docker Compose](https://devlake.apache.org/docs/GettingStarted/DockerComposeSetup) or [Install via Kubernetes](https://devlake.apache.org/docs/GettingStarted/KubernetesSetup). 
Please see detailed instructions [here](https://devlake.apache.org/docs/Overview/Introduction#how-do-i-use-devlake), and ask the community if you get stuck at any point.

### 2. Create a Blueprint

The DevLake Configuration UI will guide you through the process (a Blueprint) to define the data connections, data scope, transformation and sync frequency of the data you wish to collect.

### 3. Track the Blueprint's progress

You can track the progress of the Blueprint you have just set up.

### 4. View the pre-built dashboards

Once the first run of the Blueprint is completed, you can view the corresponding dashboards.

### 5. Customize the dashboards with SQL

If the pre-built dashboards are limited for your use cases, you can always customize or create your own metrics or dashboards with SQL.

## 😍 How to Contribute

Please read the [contribution guidelines](https://devlake.apache.org/community) before you make contribution. The following docs list the resources you might need to know after you decided to make contribution.

- [Create an Issue](https://devlake.apache.org/community/make-contribution/fix-or-create-issues): Report a bug or feature request to Apache DevLake
- [Submit a PR](https://devlake.apache.org/community/make-contribution/development-workflow): Start with [good first issues](https://github.com/apache/incubator-devlake/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) or [issues with no assignees](https://github.com/apache/incubator-devlake/issues?q=is%3Aissue+is%3Aopen+no%3Aassignee)
- [Join Mailing list](https://devlake.apache.org/community/subscribe): Initiate or participate in project discussions on the mailing list
- [Write a Blog](https://devlake.apache.org/community/make-contribution/BlogSubmission): Write a blog to share your use cases about Apache DevLake
- [Contribute a Plugin](https://devlake.apache.org/docs/DeveloperManuals/PluginImplementation): [Add a plugin](https://github.com/apache/incubator-devlake/issues?q=is%3Aissue+is%3Aopen+label%3Aadd-a-plugin+) to integrate Apache DevLake with more data sources for the community

## 📄 Contributing to Documentation:
  
- Apache DevLake's documentation is hosted at [devlake.apache.org](https://devlake.apache.org/)
- **We have a separate GitHub repository for Apache DevLake's documentation:** [github.com/apache/incubator-devlake-website](https://github.com/apache/incubator-devlake-website)

## ⌚ Project Roadmap

- <a href="https://devlake.apache.org/docs/Overview/Roadmap" target="_blank">Roadmap 2022</a>: Detailed project roadmaps for 2022.

## 💙 Community

- <a href="https://join.slack.com/t/devlake-io/shared_invite/zt-18uayb6ut-cHOjiYcBwERQ8VVPZ9cQQw" target="_blank">Slack</a>: Message us on Slack
- <a href="https://github.com/apache/incubator-devlake/wiki/FAQ" target="_blank">FAQ</a>: Frequently Asked Questions
- Wechat Community:<br/>
  ![](img/wechat_community_barcode.png)

## 📄 License<a id="license"></a>

This project is licensed under Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
