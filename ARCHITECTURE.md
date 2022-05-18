## Architecture Overview</a>

![devlake-architecture](docs/v0.11-arch-dataflow.svg)
<p align="center">Architecture Diagram (data-flow perspective)</p>

![devlake-architecture](docs/v0.11-arch-component.svg)
<p align="center">Architecture Diagram (component perspective)</p>

## Stack (from low to high)

1. config
2. logger
3. models
4. plugins
5. services
6. api / cli

## Rules

1. Higher layer calls lower layer, not the other way around
2. Whenever lower layer neeeds something from higher layer, a interface should be introduced for decoupling
3. Components should be initialized in a low to high order during bootstraping
