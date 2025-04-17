<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/storj/.github/main/assets/storj-logo-full-white.png">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/storj/.github/main/assets/storj-logo-full-color.png">
  <img alt="Storj logo" src="https://raw.githubusercontent.com/storj/.github/main/assets/storj-logo-full-color.png" width="140">
</picture>

# uplink-c

C library for Storj V3 Network.

Storj is building a decentralized cloud storage network.
[Check out our white paper for more info!](https://storj.io/white-paper)

---

Storj is an S3-compatible platform and suite of decentralized applications that
allows you to store data in a secure and decentralized manner. Your files are
encrypted, broken into little pieces and stored in a global decentralized
network of computers. Luckily, we also support allowing you (and only you) to
retrieve those files!

# Build

Download and install the latest release of Go (at least Go 1.13) at [golang.org](https://golang.org/).

When ready, building the shared library is as easy as executing:

```
make build
```

The output is in the `.build` folder.

If you specifically need GPLv2 compatibility, you can use `GPL2=true make
build` instead, which will compile the library without any Apache v2
dependencies (sadly, Apache v2 is incompatible with the GPLv2). Currently this
results in slower hashing performance (no github.com/minio/sha256-simd) and
reduced debugging and analysis infrastructure.

## Cross-Compilation

Cross-compilation is supported. The [`zig`](https://ziglang.org) compiler is required.

After all pre-requisites are installed, cross-compile by executing:

```sh
GOARCH="target-arch" GOOS="target-os" make build
```

For example, to cross-compile to an ARM64 Linux target:

```sh
GOARCH="arm64" GOOS="linux" make build
```

# API Documentation

Documentation of the stable C API is at [storj.github.io/uplink-c](https://storj.github.io/uplink-c/)

# API resource management

Functions that return a struct have allocated memory and possibly handles for
that struct.
There is a function associated with the struct that the caller must
use to free those resources.
Such a function can be recognized by the "_free_" in its name.

The rest of parameters of the functions follow the c-convention, the caller owns
them unless they mention it explicitly.

In summary:

- The caller owns the data.
- Some functions allocate on behalf of the caller (in which case there's a
  corresponding free that needs to be called).

# Examples

For some example code please take a look at [testsuite](testsuite/testplanet) folder.
Where [example_test.c](testsuite/testplanet/example_test.c) shows the most common use cases.

# License

This library is distributed under the
[MIT/expat](https://opensource.org/licenses/MIT) license.

# Support

If you have any questions or suggestions please reach out to us on
[our community forum](https://forum.storj.io/) or
email us at support@tardigrade.io.
