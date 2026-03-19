# Archaeological checklist

Where non-discoverable knowledge hides in a codebase. Work through each section systematically. Check items as you complete them.

## Build and test traps

Scan: `Makefile`, `Justfile`, `Taskfile.yml`, `scripts/`, `package.json` scripts, CI configs (`.github/workflows/`, `.gitlab-ci.yml`, `Jenkinsfile`, `.circleci/`)

- [ ] Test commands that differ from the obvious one (`pnpm vitest run` vs `pnpm test`)
- [ ] Build steps that require specific environment setup
- [ ] Scripts that look standard but have side effects
- [ ] CI-only steps that developers must replicate locally
- [ ] Pre-build or post-build hooks that are not obvious

## Environment and tooling traps

Scan: `.nvmrc`, `.node-version`, `.tool-versions`, `.python-version`, `.env.example`, `.env.template`, `docker-compose.yml`, `Dockerfile`

- [ ] Required runtime versions not obvious from package manager configs
- [ ] Environment variables that must exist for the app to start
- [ ] Docker-specific networking or volume mount requirements
- [ ] Database seed/migration commands that differ from framework defaults
- [ ] Private registries or authentication needed for package installation

## Git and workflow traps

Scan: `.husky/`, `.git/hooks/`, `.pre-commit-config.yaml`, `CONTRIBUTING.md`, `DEVELOPMENT.md`, `.github/pull_request_template.md`

- [ ] Pre-commit hooks that reject certain patterns
- [ ] Branch naming conventions
- [ ] Commit message formats enforced by hooks (not documented elsewhere)
- [ ] PR review requirements or ownership rules
- [ ] Protected branches or required status checks

## Monorepo structure

Scan: `pnpm-workspace.yaml`, `lerna.json`, `nx.json`, `turbo.json`, `rush.json`

- [ ] Package dependency relationships that are not obvious
- [ ] Build order constraints
- [ ] Shared configs that must not be overridden per-package
- [ ] Workspace-level scripts vs package-level scripts

## Code archaeology

Search for: `HACK`, `FIXME`, `XXX`, `WORKAROUND`, `DO NOT`, `NEVER`, `WARNING`, `DEPRECATED` in comments. Search for files named `legacy`, `deprecated`, `compat`.

- [ ] Areas of the codebase that are off-limits
- [ ] Workarounds for known issues that look wrong but are intentional
- [ ] Deprecated code that is still imported by production modules
- [ ] Comments explaining "why" not "what" (historical decisions)

## Architecture traps

Read actual source code in key directories (controllers, middleware, services, models).

- [ ] Custom middleware that resembles standard framework middleware but behaves differently
- [ ] Database wrappers or ORM customizations
- [ ] Auth patterns that deviate from framework defaults
- [ ] API response wrappers that all controllers must use
- [ ] Shared cache key prefixes, queue names, or event names that couple services
- [ ] Custom error classes or error handling middleware

## Configuration traps

Scan: `.npmrc`, `.yarnrc.yml`, `composer.json` config section, `.editorconfig`, `.dockerignore`, `.gitignore`

- [ ] Private registry URLs
- [ ] Resolution overrides or patched dependencies
- [ ] Ignored files that matter (generated code directories, auto-generated types)
- [ ] Non-default package manager settings (strict-peer-deps, hoist patterns)
