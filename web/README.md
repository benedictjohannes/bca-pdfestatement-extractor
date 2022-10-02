# Excel Export: Web Drag-and-Drop

## Overview

In this package, the Go program is compiled into WebAssembly, which is then called from the frontend layer.

The frontend layer itself is a [Preact](https://preactjs.com/) application written using Parcel bundler.

## Building

Both `go` and `node` (with `npm`) needs to be installed to build. As usual, the prerequisite `npm install` needs to be run beforehand.

The `go` layer needs to be built first to be processed by the frontend. As such, the `build` script inside `package.json` actually consists of three steps scripted (for Linux, my dev env):

-   `build:go:generate`: calls go generate (to extract `wasm_exec.js` for the active Go environtemnt into the `src` folder)
-   `build:go:build`: calls go build (resulting in `src/exportExcelPDFeStatement.wasm`)
-   `build:fe`: full frontend (parcel) build

## Development

In my developing the web version, I use the excellent [`gow`](https://github.com/mitranim/gow) command for automatic rebuilds. I run these commands for development inside this directory:

1. `GOOS=js GOARCH=wasm go generate`
2. `GOOS=js GOARCH=wasm gow build -o src/exportExcelPDFeStatement.wasm` (continuously rebuild on any `.go` file changes)
3. (in another terminal) `npm run start`
