# Documentation site (Docusaurus)

The hosted documentation for `dns-updater` is available at:

- https://orkarstoft.github.io/dns-updater

This folder contains the source for that documentation website, built with [Docusaurus](https://docusaurus.io/).

## Structure

- `docs/` (this directory): Docusaurus project root
- `docs/docs/`: Documentation content (Markdown/MDX)
- `docs/static/`: Static assets served as-is
- `docs/src/`: Docusaurus theme/components customizations
- `docs/docusaurus.config.ts`: Site configuration
- `docs/sidebars.ts`: Sidebar configuration

## Prerequisites

- **Node.js >= 20** (as required by `package.json`)
- npm (recommended here because a `package-lock.json` is committed)

## Install

From the repository root:

```bash
cd docs
npm ci
```

## Local development

```bash
cd docs
npm run start
```

This starts the local dev server (by default at http://localhost:3000). Changes in `docs/docs/*` are reflected live.

## Build

```bash
cd docs
npm run build
```

Build output is generated into `docs/build`.

## Serve the production build locally

```bash
cd docs
npm run serve
```

## Notes

- The content currently in `docs/docs/*` includes the default Docusaurus tutorial pages (for example `docs/docs/intro.md`). Replace or remove these as you add real project documentation.
- If you prefer Yarn, you can use it, but keep in mind this repository currently includes `package-lock.json`, so `npm ci` will produce the most reproducible installs.
