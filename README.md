<div align="center">

# SPV Wallet: Go Client

[![Release](https://img.shields.io/github/release-pre/bitcoin-sv/spv-wallet-go-client.svg?logo=github&style=flat&v=2)](https://github.com/bitcoin-sv/spv-wallet-go-client/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/bitcoin-sv/spv-wallet-go-client/run-tests.yml?branch=main&v=2)](https://github.com/bitcoin-sv/spv-wallet-go-client/actions)
[![Report](https://goreportcard.com/badge/github.com/bitcoin-sv/spv-wallet-go-client?style=flat&v=2)](https://goreportcard.com/report/github.com/bitcoin-sv/spv-wallet-go-client)
[![codecov](https://codecov.io/gh/bitcoin-sv/spv-wallet-go-client/branch/main/graph/badge.svg?v=2)](https://codecov.io/gh/bitcoin-sv/spv-wallet-go-client)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/bitcoin-sv/spv-wallet-go-client&style=flat&v=2)](https://mergify.io)
<br>

[![Go](https://img.shields.io/github/go-mod/go-version/bitcoin-sv/spv-wallet-go-client?v=2)](https://golang.org/)
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod&v=2)](https://gitpod.io/#https://github.com/bitcoin-sv/spv-wallet-go-client)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat&v=2)](https://github.com/RichardLitt/standard-readme)
[![Makefile Included](https://img.shields.io/badge/Makefile-Supported%20-brightgreen?=flat&logo=probot&v=2)](Makefile)


<br/>
</div>

## Table of Contents
- [SPV Wallet: Go Client](#spv-wallet-go-client)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Documentation](#documentation)
      - [Built-in Features](#built-in-features)
    - [Automatic Releases on Tag Creation (recommended)](#automatic-releases-on-tag-creation-recommended)
    - [Manual Releases (optional)](#manual-releases-optional)
  - [Usage](#usage)
    - [Examples \& Tests](#examples--tests)
    - [Benchmarks](#benchmarks)
  - [Code Standards](#code-standards)
  - [Usage](#usage-1)
  - [Contributing](#contributing)
  - [License](#license)

<br/>

## Installation

**spv-wallet-go-client** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/bitcoin-sv/spv-wallet-go-client
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/bitcoin-sv/spv-wallet-go-client)

For in-depth information and guidance, please refer to the [SPV Wallet Documentation](https://docs.bsvblockchain.org/network-topology/applications/spv-wallet).

[![GoDoc](https://godoc.org/github.com/bitcoin-sv/spv-wallet-go-client?status.svg&style=flat&v=2)](https://pkg.go.dev/github.com/bitcoin-sv/spv-wallet-go-client)

<br/>

<details>
<summary><strong><code>Repository Features</code></strong></summary>
<br/>

This repository was created using [MrZ's `go-template`](https://github.com/mrz1836/go-template#about)

#### Built-in Features
- Continuous integration via [GitHub Actions](https://github.com/features/actions)
- Build automation via [Make](https://www.gnu.org/software/make)
- Dependency management using [Go Modules](https://github.com/golang/go/wiki/Modules)
- Code formatting using [gofumpt](https://github.com/mvdan/gofumpt) and linting with [golangci-lint](https://github.com/golangci/golangci-lint) and [yamllint](https://yamllint.readthedocs.io/en/stable/index.html)
- Unit testing with [testify](https://github.com/stretchr/testify), [race detector](https://blog.golang.org/race-detector), code coverage [HTML report](https://blog.golang.org/cover) and [Codecov report](https://codecov.io/)
- Releasing using [GoReleaser](https://github.com/goreleaser/goreleaser) on [new Tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging)
- Dependency scanning and updating thanks to [Dependabot](https://dependabot.com) and [Nancy](https://github.com/sonatype-nexus-community/nancy)
- Security code analysis using [CodeQL Action](https://docs.github.com/en/github/finding-security-vulnerabilities-and-errors-in-your-code/about-code-scanning)
- Automatic syndication to [pkg.go.dev](https://pkg.go.dev/) on every release
- Generic templates for [Issues and Pull Requests](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository) in GitHub
- All standard GitHub files such as `LICENSE`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, and `SECURITY.md`
- Code [ownership configuration](.github/CODEOWNERS) for GitHub
- All your ignore files for [vs-code](.editorconfig), [docker](.dockerignore) and [git](.gitignore)
- Automatic sync for [labels](.github/labels.yml) into GitHub using a pre-defined [configuration](.github/labels.yml)
- Built-in powerful merging rules using [Mergify](https://mergify.io/)
- Welcome [new contributors](.github/mergify.yml) on their first Pull-Request
- Follows the [standard-readme](https://github.com/RichardLitt/standard-readme/blob/master/spec.md) specification
- [Visual Studio Code](https://code.visualstudio.com) configuration with [Go](https://code.visualstudio.com/docs/languages/go)
- (Optional) [Slack](https://slack.com), [Discord](https://discord.com) or [Twitter](https://twitter.com) announcements on new GitHub Releases
- (Optional) Easily add [contributors](https://allcontributors.org/docs/en/bot/installation) in any Issue or Pull-Request

</details>

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

- [stretchr/testify](https://github.com/stretchr/testify)
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

Releases are automatically created when you create a new [git tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging)!

If you want to manually make releases, please install GoReleaser:

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to GitHub and can be installed:
- **using make:** `make install-releaser`
- **using brew:** `brew install goreleaser`

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

<br/>

### Automatic Releases on Tag Creation (recommended)
Automatic releases via [GitHub Actions](.github/workflows/release.yml) from creating a new tag:
```shell
make tag version=1.2.3
```

<br/>

### Manual Releases (optional)
Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production (manually).

<br/>

</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                           Runs multiple commands
clean                         Remove previous builds and any cached data
clean-mods                    Remove all the Go mod cache
coverage                      Shows the test coverage
diff                          Show the git diff
generate                      Runs the go generate command in the base of the repo
godocs                        Sync the latest tag with GoDocs
help                          Show this help message
install                       Install the application
install-all-contributors      Installs all contributors locally
install-go                    Install the application (Using Native Go)
install-releaser              Install the GoReleaser application
lint                          Run the golangci-lint application (install if not found)
release                       Full production release (creates release in GitHub)
release                       Runs common.release then runs godocs
release-snap                  Test the full release (build binaries)
release-test                  Full production test release (everything except deploy)
replace-version               Replaces the version in HTML/JS (pre-deploy)
tag                           Generate a new tag and push (tag version=0.0.0)
tag-remove                    Remove a tag if found (tag-remove version=0.0.0)
tag-update                    Update an existing tag to current commit (tag-update version=0.0.0)
test                          Runs lint and ALL tests
test-ci                       Runs all tests via CI (exports coverage)
test-ci-no-race               Runs all tests via CI (no race) (exports coverage)
test-ci-short                 Runs unit tests via CI (exports coverage)
test-no-lint                  Runs just tests
test-short                    Runs vet, lint and tests (excludes integration tests)
test-unit                     Runs tests and outputs coverage
uninstall                     Uninstall the application (and remove files)
update-contributors           Regenerates the contributors html/list
update-linter                 Update the golangci-lint package (macOS only)
vet                           Run the Go vet application
```
</details>

<br/>

## Usage
Checkout all the [examples](examples)!

<br/>

### Examples & Tests
All unit tests and [examples](examples) run via [GitHub Actions](https://github.com/bitcoin-sv/spv-wallet-go-client/actions) and
uses [Go version 1.19.x](https://golang.org/doc/go1.19). View the [configuration file](.github/workflows/run-tests.yml).

<br/>

Run all tests (including integration tests)
```shell script
make test
```

<br/>

Run tests (excluding integration tests)
```shell script
make test-short
```

<br/>

### Benchmarks
Run the Go benchmarks:
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## Usage


```
// http example
func main() {

	// Generate keys
	keys, _ := xpriv.Generate()

	// Create a client
	client, _ := walletclient.New(
        walletclient.WithXPriv(keys.XPriv()),
        walletclient.WithHTTP("localhost:3001"),
        walletclient.WithSignRequest(true))
    
    fmt.Println(client.IsSignRequest())
}

```

Checkout all the [examples](examples)!

<br/>

## Contributing
All kinds of contributions are welcome!
<br/>
To get started, take a look at [code standards](.github/CODE_STANDARDS.md).
<br/>
View the [contributing guidelines](.github/CODE_STANDARDS.md#3-contributing) and follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

<br/>

## License

[![License](https://img.shields.io/github/license/bitcoin-sv/spv-wallet-go-client.svg?style=flat&v=2)](LICENSE)
