---
sidebar_position: 2
title: Getting Started
---

# Getting Started

## Prerequisites

- A supported DNS provider account (GCP Cloud DNS or DigitalOcean DNS)
- A DNS zone already created in that provider
- A machine/network where your public IP can change (home connection, lab, etc.)

## Quick start

1. Create a `config.yaml` file (see the **Configuration** section).
2. Run `dns-updater` either:
   - as a local binary (CLI), or
   - as a Docker container.

`dns-updater` loads `config.yaml` from the **current working directory**.
