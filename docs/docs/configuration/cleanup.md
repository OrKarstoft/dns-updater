---
title: Cleanup
sidebar_position: 6
---

# Cleanup

Use cleanup to delete managed DNS records and their safemode TXT records for everything listed in the configuration file.

## Run cleanup

Run the binary with the `--clean` flag:

```bash
./dns-updater --clean
```

With Docker:

```bash
docker run -v /path/to/config.yaml:/config.yaml ghcr.io/orkarstoft/dns-updater:latest --clean
```

## Safety guard

Cleanup relies on [Safe Mode](./safe-mode.md) TXT records.

For each configured record, `dns-updater` checks that the matching safemode TXT record exists and contains the expected ownership data. It only deletes the A record and its safemode TXT record when that ownership check passes.

If the TXT record is missing or its ownership data does not match, `dns-updater` skips that record and logs an error. This prevents cleanup from deleting records that `dns-updater` does not own.
