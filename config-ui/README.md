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
# DevLake - Configuration UI

The **Config-UI Application** is a **React.js** SPA (Single-Page-Application) that manages the setup and configuration of a **DevLake** Instance.

#### Technology / Stack Overview

- React
- Antd
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
