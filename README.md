# xk6-cli-wrapper

A k6 extension for wrapping CLI app

## Build

To build a `k6` binary with this plugin, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Install `xk6`:

  ```shell
  go install github.com/k6io/xk6/cmd/xk6@latest
  ```

2. Build the binary:

  ```shell
  xk6 build master \
    --with github.com/nanang-ab/xk6-cli-wrapper
  ```

## Example

Please refer the [sample k6 script](cli-sample-k6-script.js)
