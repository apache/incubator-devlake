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