# Basic Usage Guide

## Quick Start

**Linux (incl. iOS iSH) / macOS / Android Termux**:
(Requires write permissions for /usr/bin, /usr/local/bin, or $PREFIX/bin)

```bash
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)
```

**Windows PowerShell** (Administrator):

```ps
irm https://unlock.icmp.ing/scripts/download_test.ps1 | iex
```

## Command Line Arguments

### Basic Options

|Flag|Description|Example|
|---|---|---|
|`-m`|Mode: 0=Auto (Default), 4=IPv4 Only, 6=IPv6 Only|`-m 4` Only test IPv4|
|`-v`|Show version and exit|`-v`|
|`-u`|Check and update to latest version|`-u`|

### Performance Optimization

|Flag|Description|Example|
|---|---|---|
|`-conc`|Max concurrent tests (0=unlimited)|`-conc 50` Limit to 50 concurrent|
|`-cache`|Enable caching and serial region execution|`-cache` Enable caching|
|`-show-active`|Show active tests in progress bar|`-show-active=false` Disable active display|

### Debugging & Testing

|Flag|Description|Example|
|---|---|---|
|`-debug`|Enable debug mode (verbose error output)|`-debug`|
|`-nf`|Test Netflix availability only|`-nf`|
|`-test`|Run specific tests|`-test`|

## Common Use Cases

```bash
# Run all tests (Default)
./unlock-test

# Run only IPv4 tests
./unlock-test -m 4

# Limit concurrency to 30 (Good for low-end devices)
./unlock-test -conc 30

# Enable debug mode to see detailed errors
./unlock-test -debug
```
