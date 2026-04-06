---
title: Cache
sidebar_position: 4
---

# Cache

Caching stores the last observed public IP so the updater can skip DNS updates when the IP hasn’t changed.

```yaml
cache:
  enabled: true
  filePath: "dns-updater.cache"
```

## `cache.enabled: true`

A local file cache is used:

the last IP is written to cache.filePath
on the first run the cache file may not exist yet, which is expected.

## `cache.enabled: false`

A no-op cache is used:

GetLastIP always returns nil
this effectively forces an update whenever the updater runsn

## `cache.filePath: "dns-updater.cache"`

The path to the cache file. This can be an absolute or relative path. If the file does not exist, it will be created when the updater runs.
