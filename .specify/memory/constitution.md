<!--
Sync Impact Report:
- Version change: N/A → 1.0.0 (initial constitution)
- Modified principles: N/A (new)
- Added sections: Core Principles (5), Plugin Architecture, Development Workflow, Governance
- Removed sections: None
- Templates requiring updates:
  - ✅ plan-template.md (Constitution Check section exists, principles align)
  - ✅ spec-template.md (requirements structure compatible)
  - ✅ tasks-template.md (phase structure compatible with plugin development)
- Follow-up TODOs: None
-->

# Apache DevLake Constitution

## Core Principles

### I. Plugin Independence
All plugins MUST be self-contained and independently deployable. Plugins reside in `backend/plugins/<name>/` with no cross-plugin imports. Each plugin MUST implement required interfaces (`PluginMeta`, `PluginTask`, `PluginModel`, `PluginMigration`) as documented in [AGENTS.md](../../AGENTS.md). Plugin independence enables modular data source integration and parallel development.

### II. Three-Layer Data Model
Data transformation MUST follow the Raw → Tool → Domain pipeline:
- **Raw Layer** (`_raw_*` tables): Store unmodified API responses as JSON for replay/debugging
- **Tool Layer** (`_tool_<plugin>_*` tables): Plugin-specific models extracted from raw data
- **Domain Layer** (standardized tables): Normalized models in `backend/core/models/domainlayer/`

Collectors write to Raw, Extractors write to Tool, Converters write to Domain. This separation ensures data lineage and enables cross-plugin correlation.

### III. Test-Driven Development
Tests MUST be written before implementation for all new functionality:
- **Unit tests**: `*_test.go` files alongside source, use mocks from `backend/mocks/`
- **E2E tests**: CSV fixtures in `e2e/` directory for extractor/converter validation
- **Integration tests**: Use `backend/test/helper/` client for full-stack testing

All models MUST be registered in `GetTablesInfo()` or CI fails via `plugins/table_info_test.go`.

### IV. Migration-First Schema Changes
Database schema changes MUST use migration scripts:
- Located in `models/migrationscripts/`
- Registered in `register.go`'s `All()` function
- Version format: `YYYYMMDD_description.go`
- Migrations run sequentially on server startup

Never modify existing migrations; create new ones for schema evolution.

### V. Apache Compliance
All contributions MUST comply with Apache Software Foundation requirements:
- Apache 2.0 license header required on ALL source files
- No proprietary dependencies or incompatible licenses
- Follow contribution guidelines at [devlake.apache.org/community](https://devlake.apache.org/community)

## Plugin Architecture

Each Go plugin in `backend/plugins/<name>/` MUST follow this structure:
```
api/         # REST endpoints (connections, scopes, scope-configs)
impl/        # Plugin implementation (implements core interfaces)
models/      # Tool layer models + migrationscripts/
tasks/       # Collectors, Extractors, Converters
e2e/         # Integration tests with CSV fixtures
```

Required interfaces per plugin type:
- **Data Source Plugins**: `PluginMeta`, `PluginTask`, `PluginModel`, `PluginMigration`, `PluginSource`, `DataSourcePluginBlueprintV200`
- **Helper Plugins** (e.g., dora, refdiff): `PluginMeta`, `PluginTask`

API changes require running `make swag` to regenerate Swagger documentation.

## Development Workflow

### Build Commands
```bash
make dep        # Install dependencies
make build      # Build plugins + server
make dev        # Build + run server on :8080
make unit-test  # Run unit tests
make e2e-test   # Run E2E tests
make lint       # Run golangci-lint (from backend/)
make swag       # Regenerate Swagger docs (from backend/)
```

### Quality Gates
Before submitting PR:
1. All tests pass (`make unit-test`)
2. Linting passes (`make lint`)
3. New models registered in `GetTablesInfo()`
4. New migrations registered in `All()`
5. Swagger updated if API changed (`make swag`)

## Governance

This constitution supersedes informal practices and establishes non-negotiable standards for Apache DevLake development.

**Amendment Process**:
1. Propose changes via GitHub Issue or mailing list discussion
2. Obtain consensus from project maintainers
3. Update this document with version increment
4. Document migration plan for existing code if principles change

**Compliance Verification**:
- All PRs MUST be reviewed for principle compliance
- CI/CD enforces model registration and test requirements
- Use [AGENTS.md](../../AGENTS.md) for detailed development guidance

**Version**: 1.0.0 | **Ratified**: 2025-12-05 | **Last Amended**: 2025-12-05
