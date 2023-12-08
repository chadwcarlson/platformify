<p align="center">
<a href="https://www.upsun.com/">
<img src="internal/utils/logo.svg" width="500px">
</a>
</p>

<p align="center">
<a href="https://github.com/platformsh/platformify/issues">
<img src="https://img.shields.io/github/issues/platformsh/platformify.svg?style=for-the-badge&labelColor=f4f2f3&color=6046FF&label=Issues" alt="Open issues" />
</a>&nbsp&nbsp
<a href="https://github.com/platformsh/platformify/pulls">
<img src="https://img.shields.io/github/issues-pr/platformsh/platformify.svg?style=for-the-badge&labelColor=f4f2f3&color=6046FF&label=Pull%20requests" alt="Open PRs" />
</a>&nbsp&nbsp
<a href="https://github.com/platformsh/platformify/blob/main/LICENSE">
<img src="https://img.shields.io/static/v1?label=License&message=MIT&style=for-the-badge&labelColor=f4f2f3&color=6046FF" alt="License" />
</a>&nbsp&nbsp
<br /><br />

<h2 align="center"><code>project:init</code> for Platform.sh and Upsun</h2>

Get your project ready to be deployed in Platform.sh and Upsun!

This project supplies the `project:init` and `validate` subcommands imported into both the [Platform.sh](https://docs.platform.sh/administration/cli.html) and [Upsun CLI](https://docs.upsun.com/administration/cli.html).

## Contributions

1. [Create a fork of this repository](https://github.com/platformsh/platformify/fork), and clone locally.
1. Create a new branch for your contribution.

1. Make your changes, then [build locally to test](#building).
1. Before opening a pull request, you can [run tests locally](#tests).


## Building

In order to build the binary, use the following:

```console
go build ./cmd/platformify/
<!-- or -->
go build ./cmd/upsunify/
```

Run the CLI with:

```console
./platformify
<!-- or -->
./upsunify
```

Alternatively, you can build both commands with `make`:

```bash
make build
```

## Tests

1. Clean

    ```bash
    make clean
    ```

1. Tidy    

    ```bash
    go mod tidy
    ```

1. Format

    ```bash
    gofmt -s -w .
    ```

1. Build

    ```bash
    make build
    ```

1. Linter

    ```bash
    make lint
    ```

    > [!NOTE]
    > It may be necessary on failure to [install golangci-lint](https://golangci-lint.run/usage/install/#local-installation). 
    > For Mac, this is `brew install golangci-lint`, and `export PATH="/Users/yourusername/go/bin/:$PATH"`.
    > Replace the path above for your settings.
    > This needs to be the full path. `~` is not enough. Use `GOPATH` if available, otherwise spell out the full path.


    > [!NOTE]
    > If `make lint` results in a failure for a file you've modified (`File is not gotfmt-ed with -s`), run `gofmt -s -w .`.

1. Run tests

    ```bash
    make test
    ```
