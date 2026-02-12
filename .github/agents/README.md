<!--
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Custom GitHub Copilot Agents for DevLake

This directory contains custom GitHub Copilot agents that provide specialized assistance for working with Apache DevLake.

## Available Agents

### Azure Deployment Agent

**File:** `azure-deployment.agent.md`

**Purpose:** Provides interactive, step-by-step guidance for deploying Apache DevLake to Microsoft Azure.

**Supported Deployment Methods:**
- Docker container on Azure Container Instances (ACI)
- Azure Kubernetes Service (AKS) for production deployments

**Key Features:**
- Interactive workflow with user prompts at each stage
- Prerequisites validation (Azure CLI, authentication, subscriptions)
- Automated resource creation (ACR, databases, AKS clusters)
- **Azure MCP Tools Integration** - Direct access to Azure resources via MCP server
  - Resource management (resource groups, deployments, templates)
  - Storage operations (blobs, containers, file shares)
  - AKS cluster management and scaling
  - Cosmos DB queries and management
  - Key Vault secret retrieval
  - App Configuration key-value management
  - Azure Monitor log queries and metrics
- Security best practices and recommendations
- Post-deployment validation and monitoring setup
- Comprehensive troubleshooting guide

## How to Use Custom Agents

### Prerequisites
1. **GitHub Copilot subscription** - You need an active GitHub Copilot subscription
2. **VS Code with GitHub Copilot extension** - Install the GitHub Copilot extension in VS Code
3. **Repository access** - Clone this repository to your local machine

### Using the Azure Deployment Agent

1. **Open VS Code** in the DevLake repository root directory

2. **Invoke GitHub Copilot Chat** (Ctrl+Shift+I or Cmd+Shift+I)

3. **Reference the agent** by mentioning it in your prompt:
   ```
   @azure-deployment help me deploy DevLake to Azure
   ```
   or simply:
   ```
   I want to deploy DevLake to Azure
   ```
   
   If the agent has `infer: false` (as currently configured), you need to explicitly mention it.

4. **Follow the interactive prompts** - The agent will:
   - Verify your prerequisites
   - Ask about deployment preferences
   - Guide you through resource creation
   - Execute deployment commands
   - Validate the deployment

### Example Interaction

```
You: @azure-deployment I want to deploy DevLake to Azure using AKS

Agent: I'd be happy to help you deploy Apache DevLake to Azure using AKS! 
Let me walk you through the process step by step.

First, let me verify some prerequisites:
1. Do you have the Azure CLI (az) installed and are you logged in?
2. Do you have an active Azure subscription with permissions to create resources?

Once we confirm these, I'll ask you a few questions about your deployment preferences.

You: Yes, I'm logged in with az cli and have an active subscription

Agent: Great! Now let's gather some information about your deployment...
[continues with interactive workflow]
```

## Agent Configuration

Custom agents are configured using YAML frontmatter in Markdown files. Key configuration options:

- **name**: Display name for the agent
- **description**: Brief description of the agent's purpose
- **target**: Where the agent can be used (e.g., `github-copilot`, `vscode`)
- **tools**: List of allowed tools the agent can use
- **infer**: Whether Copilot can automatically invoke this agent (true/false)
- **metadata**: Additional metadata for categorization and filtering

## Best Practices

1. **Be Specific**: When invoking an agent, clearly state your goal
2. **Provide Context**: Share relevant information about your environment
3. **Follow Instructions**: The agents are designed to be interactive - answer their questions
4. **Review Commands**: Always review commands before executing them, especially those that create Azure resources
5. **Security**: Never share sensitive information like passwords or API keys in agent interactions

## Customizing Agents

You can customize existing agents or create new ones:

1. Create a new `.agent.md` file in this directory
2. Add YAML frontmatter with agent configuration
3. Write comprehensive instructions, examples, and boundaries
4. Test the agent in VS Code with GitHub Copilot

## Troubleshooting

**Agent not appearing in Copilot:**
- Ensure you're using the latest GitHub Copilot extension
- Check that the agent file has correct YAML frontmatter
- Try reloading VS Code window

**Agent not responding as expected:**
- Review the agent's instructions in the `.agent.md` file
- Ensure you're providing the information requested
- Check that prerequisites are met

**Commands failing:**
- Verify Azure CLI is installed and you're logged in
- Check that you have proper permissions
- Review error messages and follow troubleshooting steps in the agent guide

## Contributing

To contribute a new agent or improve existing ones:

1. Follow the agent configuration format
2. Include comprehensive instructions and examples
3. Set clear boundaries and limitations
4. Test thoroughly before submitting
5. Update this README with agent documentation

## Resources

- [GitHub Copilot Documentation](https://docs.github.com/en/copilot)
- [Custom Agents Configuration](https://docs.github.com/en/copilot/reference/custom-agents-configuration)
- [Apache DevLake Documentation](https://devlake.apache.org/)
- [Azure CLI Documentation](https://docs.microsoft.com/en-us/cli/azure/)

## License

Licensed under the Apache License, Version 2.0. See the LICENSE file in the repository root for details.
