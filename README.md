# ESF Products Manager

## Requirements

- [Go](https://go.dev)
- [CockroachDB](https://www.cockroachlabs.com)
- [pnpm](https://pnpm.js.org)

## Start

```bash
# Install Editor Dependencies && Build
cd Editor && pnpm install && pnpm run build && cd ../
# Install DB Migration Tool
./_tasks launch_setup
# Start DB In Background
./_tasks launch_cockroachdb
# Migrate DB
./_tasks do_migrate up
# Start Server
./_tasks run
```
