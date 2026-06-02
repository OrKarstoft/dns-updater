---
title: Logging
sidebar_position: 3
---

# Logging

`dns-updater` uses **zerolog** for logging. Logging is configured in `config.yaml` under the `log` section.

## Configuration

```yaml
log:
  level: info
  type: pretty
```

## `log.level`

log.level controls the global log level. It is parsed using zerolog.ParseLevel(...).

If you provide an invalid value, dns-updater falls back to info and emits a warning.
Common values include: debug, info, warn, error.

Example:

```yaml
log:
  level: debug
  type: pretty
```

## `log.type`

log.type controls the log output format/destination.

Supported values:

- `pretty`
  Human-readable console output.

- `json`
  Structured JSON logs to stdout (this is also the default if log.type is not recognized).

- `file`
  Currently not a true file logger. The code contains commented example logic for writing to a file, but the current implementation still logs to stdout.

## Defaults

If log is omitted, these defaults are applied:

```yaml
log:
  level: info
  type: pretty
```
