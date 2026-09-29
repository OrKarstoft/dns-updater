# Documentation site (Docusaurus)

The hosted documentation for `dns-updater` is available at:

- https://orkarstoft.github.io/dns-updater

This folder contains the source for that documentation website, built with [Docusaurus](https://docusaurus.io/).

## Structure

- `docs/` (this directory): Docusaurus project root; npm dependencies are installed on demand
- `docs/docs/`: Documentation content (Markdown/MDX)
- `docs/static/`: Static assets served as-is
- `docs/src/`: Docusaurus theme/components customizations
- `docs/docusaurus.config.ts`: Site configuration
- `docs/sidebars.ts`: Sidebar configuration

## Prerequisites

- Node.js 24 and npm (matching `.github/workflows/docs.yml`)

## Install

From the repository root, change to `docs/` and run the **Install dependencies** commands in
[`../.github/workflows/docs.yml`](../.github/workflows/docs.yml). They create a
temporary, ignored `package.json` and install pinned direct dependencies without
writing a lockfile. This repository does not commit npm manifests; the workflow
is the source of truth for dependency versions.

## Local development

```bash
cd docs
./node_modules/.bin/docusaurus start
```

This starts the local dev server (by default at http://localhost:3000). Changes in `docs/docs/*` are reflected live.

## Build

```bash
cd docs
./node_modules/.bin/docusaurus build
```

Build output is generated into `docs/build`.

## Serve the production build locally

```bash
cd docs
./node_modules/.bin/docusaurus serve
```

## Notes

- The content currently in `docs/docs/*` includes the default Docusaurus tutorial pages (for example `docs/docs/intro.md`). Replace or remove these as you add real project documentation.
- With no committed lockfile, transitive dependencies can change between builds,
  and GitHub cannot report vulnerabilities in a tracked npm dependency graph.
  Pinned direct versions reduce drift but do not make builds fully reproducible.
