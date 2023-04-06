# DevLake - Configuration UI

The **Config-UI Application** is a **React.js** SPA (Single-Page-Application) that manages the setup and configuration of a **DevLake** Instance.

#### Technology / Stack Overview

- React
- Blueprint
- Vite
- TypeScript
- Yarn3

## Development

In order to develop on this project you will need a properly working **React Development Environment**.

### Environment Setup

Install Package Dependencies before attempting to start the UI. The application will not start unless all packages are installed without errors.

#### Install Dependencies

```
$ yarn
```

#### Start Development Server

❗ Please ensure the **DevLake API** is **online** before starting the **UI**, otherwise the application will remain in offline mode and errors will be displayed.

```
$ yarn start
```

Server will listen on `http://localhost:4000`

#### Production Build

To Build static and minified production assets to the `dist/` directory.

```
$ yarn build
```

#### TEST / RUN Production Build

Build production assets and listen to emulate a production environment. This is to verify minified bundled/assets are operating correctly.

```
$ yarn preview
```

For actual production use, the **Docker Image** for Config-UI should be used as outlined in the main project README.md
