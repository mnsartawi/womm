<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="public/womm-round.png">
    <img src="public/womm-round.png" alt="womm" width="250" />
  </picture>
</p>



# Works on my machine — Tauri + Next.js

GUI companion for the [womm CLI](https://github.com/mnsartawi/womm). Makes it easier to generate, edit, and validate `.womm` config files without leaving a visual interface.

## Prerequisites

- [womm CLI](https://github.com/mnsartawi/womm) installed on your system
- Node.js 18+
- Rust toolchain (for Tauri)

## Getting started

```bash
npm install
npm run tauri dev
```

## Build

```bash
npm run tauri build
```

## Features

- Generate starter `.womm` files
- Validate environment against a `.womm` config
- Edit config files inline
- Runs the CLI under the hood — same checks, nicer output
