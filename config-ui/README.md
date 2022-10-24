# DevLake - Configuration UI

The **Config-UI Application** is a **React.js** SPA (Single-Page-Application) that manages the setup and configuration of a **DevLake** Instance.

#### Technology / Stack Overview
- React.js
- BlueprintJS
- Cypress
- Webpack

## Development
In order to develop on this project you will need a properly working **React Development Environment**.

### Environment Setup
Install Package Dependencies before attempting to start the UI. The application will not start unless all packages are installed without errors.

#### Install NPM Dependencies
```
[config-ui@main] $> npm -i
```

#### Start Development Server
❗ Please ensure the **DevLake API** is **online** before starting the **UI**, otherwise the application will remain in offline mode and errors will be displayed.

```
[config-ui@main] $> npm start
```
Server will listen on `http://localhost:4000`

#### Production Build
To Build static and minified production assets to the `dist/` directory.
```
[config-ui@main] $> npm run build
```

#### TEST / RUN Production Build
Build production assets and listen to emulate a production environment. This is to verify minified bundled/assets are operating correctly.
```
[config-ui@main] $> npm run start-production
```
Server will listen on `http://localhost:9000`

For actual production use, the **Docker Image** for Config-UI should be used as outlined in the main project README.md

## Plugin Registration

### Step 1. Create Registry JSON Configuration
Depending on the nature of the plugin, see a similar integration in `registry/plugins/` folder and copy as base template. For example, to model after GitHub copy `registry/plugins/github.json` to a uniquely named JSON registry file.

```
# Example "Merico"
$> cp registry/plugins/github.json registry/plugins/merico.json"
```

**Configure Plugin Options**
The newly created registry file needs to customized with **Name** and **Type** (`Default="integration"`), and all related **Connection** property fields (_Labels_, Tooltips, etc)

Since not all Plugins created may have **UI** capabilities, the `type` property is used to distinguish _variance_ in behavior.

- integration
- plugin
- pipeline
- webhook

The **Enabled** property (`enabled`) must be set to `true` to allow the plugin the be registered. 

```
# Sample Config for "Merico" Plugin
{
  "id": "merico",
  "name": "Merico",
  "type": "integration",
  "enabled": true,
  "multiConnection": true,
  "isBeta": false,
  "isProvider": true,
  "icon": "src/images/integrations/merico.svg",
  ...
  ...
  ...
}
```

### Step 2. Register the JSON Configuration
Next, the new plugin needs to be registered with the **Integrations Manager Hook** (`@hooks/useIntegrations.jsx`). 
1. Import the new Plugin JSON File
2. Add Plugin to the bottom of `pluginRegistry` Array

```
$> vi hooks/useIntegrations.jsx`
....
....
import DbtPlugin from '@/registry/plugins/dbt.json'
import StarrocksPlugin from '@/registry/plugins/starrocks.json'
import DoraPlugin from '@/registry/plugins/dora.json'
# Register Merico Plugin
+ import MericoPlugin from '@/registry/plugins/merico.json'

function useIntegrations(
  pluginRegistry = [
    JiraPlugin,
    GitHubPlugin,
    ...
    ...
    # Load Merico Plugin
    + MericoPlugin

```

### Step 3. Define & Register Transformation Settings Component (Optional)

If this new Plugin requires **Transformation Settings**, the transformation settings component must be created and imported. You may use an existing transformation settings file as a reference. In the future there will be support to create transformation settings dynamically from configuration.

1. Create and Design React Transformation **JSX Component** in `settings/[my-plugin.jsx]`

```jsx
# Example Transformation Settings Template for "Merico"
export default function MericoSettings(props) {
  const {
    Providers,
    ProviderLabels,
    provider,
    connection,
    entities = [],
    transformation = {},
    isSaving = false,
    isSavingConnection = false,
    onSettingsChange = () => {}
  } = props

  return (
    <>
     
    </>
  )
}
```

3. Import & **Register Provider Transformation Settings** in `components/blueprints/ProviderTransformationSettings.jsx`.  Update the `TransformationComponents` Map to Load the Plugin

```jsx
# components/blueprints/ProviderTransformationSettings.jsx
import AzureSettings from '@/pages/configure/settings/azure'
import BitbucketSettings from '@/pages/configure/settings/bitbucket'
import GiteeSettings from '@/pages/configure/settings/gitee'
# Import/Register Merico Plugin
+ import MericoSettings from '@/pages/configure/settings/merico'

const ProviderTransformationSettings = (props) => { 
  ...
  ...
  ...
  // Provider Transformation Components (LOCAL)
  const TransformationComponents = useMemo(
    () => ({
      ...
      ...
      [Providers.AZURE]: AzureSettings,
      [Providers.BITBUCKET]: BitbucketSettings,
      [Providers.GITEE]: GiteeSettings
      # Load Merico Plugin
      + [Providers.MERICO]: MericoSettings
    }),
    [Providers]
  )

}
```


## Testing

### Cypress E2E Tests
The Cypress Test Runner has been installed and configured for `@config-ui`. **Integration Tests** are located at `@config-ui/cypress/integration`.

Before writing tests please read official **Cypress** Documentation at [https://docs.cypress.io/](https://docs.cypress.io/)

#### Integration Tests Coverage
Test **Specs** are organized in the following main **groups**. New tests cases should be added to the appropriate spec group file(s), or a new group should be created if necessary for new criteria, or to create alternate flows.

| Integration Group       | Tests             | Status |
| ----------------------- | ----------------- | -----: |
| 0-api                   | 5                 | PASS   |
| 1-application           | 12                | PASS   |
| 2-data-integrations     | 17                | PASS   |
| 3-pipelines             | 30                | PASS   |

Test cases will be updated and created as necessary when new features are added and important bug-fixes are applied.

#### Open/RUN Cypress
Once the **Cypress Dashboard** opens, choose a test file from the available list of integration tests.

```
[config-ui@main] $> npm run cypress
```