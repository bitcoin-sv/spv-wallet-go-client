<div align="center">

# 🧰&nbsp;&nbsp;spv-wallet-go-client

**Lightweight Go client library for interacting with the SPV Wallet**

<br/>

<a href="https://github.com/bsv-blockchain/spv-wallet-go-client/releases"><img src="https://img.shields.io/github/release-pre/bsv-blockchain/spv-wallet-go-client?include_prereleases&style=flat-square&logo=github&color=black" alt="Release"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/bsv-blockchain/spv-wallet-go-client?style=flat-square&logo=go&color=00ADD8" alt="Go Version"></a>
<a href="https://github.com/bsv-blockchain/spv-wallet-go-client/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-OpenBSV-blue?style=flat-square" alt="License"></a>

<br/>

<table align="center" border="0">
  <tr>
    <td align="right">
       <code>CI / CD</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/bsv-blockchain/spv-wallet-go-client/actions"><img src="https://img.shields.io/github/actions/workflow/status/bsv-blockchain/spv-wallet-go-client/fortress.yml?branch=main&label=build&logo=github&style=flat-square" alt="Build"></a>
       <a href="https://github.com/bsv-blockchain/spv-wallet-go-client/actions"><img src="https://img.shields.io/github/last-commit/bsv-blockchain/spv-wallet-go-client?style=flat-square&logo=git&logoColor=white&label=last%20update" alt="Last Commit"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Quality</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://codecov.io/gh/bsv-blockchain/spv-wallet-go-client"><img src="https://codecov.io/gh/bsv-blockchain/spv-wallet-go-client/branch/main/graph/badge.svg?style=flat-square" alt="Coverage"></a>
    </td>
  </tr>

  <tr>
    <td align="right">
       <code>Security</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://scorecard.dev/viewer/?uri=github.com/bsv-blockchain/spv-wallet-go-client"><img src="https://api.scorecard.dev/projects/github.com/bsv-blockchain/spv-wallet-go-client/badge?style=flat-square" alt="Scorecard"></a>
       <a href=".github/SECURITY.md"><img src="https://img.shields.io/badge/policy-active-success?style=flat-square&logo=security&logoColor=white" alt="Security"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Community</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/bsv-blockchain/spv-wallet-go-client/graphs/contributors"><img src="https://img.shields.io/github/contributors/bsv-blockchain/spv-wallet-go-client?style=flat-square&color=orange" alt="Contributors"></a>
       <a href="https://deepwiki.com/bsv-blockchain/spv-wallet-go-client"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki"></a>
    </td>
  </tr>
</table>

</div>

<br/>
<br/>

<div align="center">

### <code>Project Navigation</code>

</div>

<table align="center">
  <tr>
    <td align="center" width="33%">
       📦&nbsp;<a href="#-installation"><code>Installation</code></a>
    </td>
    <td align="center" width="33%">
       🧪&nbsp;<a href="#-examples--tests"><code>Examples&nbsp;&&nbsp;Tests</code></a>
    </td>
    <td align="center" width="33%">
       📚&nbsp;<a href="#-documentation"><code>Documentation</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🤝&nbsp;<a href="#-contributing"><code>Contributing</code></a>
    </td>
    <td align="center">
       🛠️&nbsp;<a href="#-code-standards"><code>Code&nbsp;Standards</code></a>
    </td>
    <td align="center">
       ⚡&nbsp;<a href="#-benchmarks"><code>Benchmarks</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🤖&nbsp;<a href="#-ai-usage--assistant-guidelines"><code>AI&nbsp;Usage</code></a>
    </td>
    <td align="center">
       📝&nbsp;<a href="#-license"><code>License</code></a>
    </td>
    <td align="center">
       👥&nbsp;<a href="#-maintainers"><code>Maintainers</code></a>
    </td>
  </tr>
</table>
<br/>

## 📦 Installation

**spv-wallet-go-client** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/bsv-blockchain/spv-wallet-go-client
```

<br/>

## 📚 Documentation

- **API Reference** – Dive into the godocs at [pkg.go.dev/github.com/bsv-blockchain/spv-wallet-go-client](https://pkg.go.dev/github.com/bsv-blockchain/spv-wallet-go-client)
- **Usage Examples** – Browse practical patterns either the [examples directory](examples) or the example tests
- **Test Suite** – Review both the unit tests and fuzz tests(powered by [`testify`](https://github.com/stretchr/testify))
- **SPV Wallet Docs** - please refer to the [SPV Wallet Documentation](https://docs.bsvblockchain.org/network-topology/spv-wallet)

<br/>

<details>
<summary><strong><code>Quick Start Guide</code></strong></summary>

## Quick Start Guide

The implementation enforces separation of concerns by isolating admin and non-admin APIs, requiring separate initialization for their respective clients. This ensures clarity and modularity when utilizing the exposed functionality.

### `UserAPI` Initialization Methods:

### 1. [`NewUserAPIWithAccessKey`](/user_api.go#L468)
- **Description:** Initializes a `UserAPI` instance using an access key for authentication.
- **Note:** Requests made with this instance will be securely signed, ensuring integrity and authenticity.

### 2. [`NewUserAPIWithXPriv`](/user_api.go#L449)
- **Description:** Initializes a `UserAPI` instance using an extended private key (xPriv) for authentication.
- **Note:** Requests made with this instance will also be securely signed.
- **Recommendation:** This option offers a high level of security, making it a preferred choice alongside the access key option.

### 3. [`NewUserAPIWithXPub`](/user_api.go#L435)
- **Description:** Initializes a `UserAPI` instance using an extended public key (xPub).
- **Note:** Requests made with this instance will not be signed.
- **Security Advisory:** For enhanced security, it is strongly recommended to use either `NewUserAPIWithAccessKey` or `NewUserAPIWithXPriv` instead, as unsigned requests may be less secure.


### `AdminAPI` Initialization Methods:

### 1. [`NewAdminAPIWithXPriv`](/admin_api.go#L375)
- **Description:** Initializes a `AdminAPI` instance using an extended private key (xPriv) for authentication.
- **Note:** Requests made with this instance will be securely signed, ensuring integrity and authenticity.

### 2. [`NewAdminAPIWithXPub`](/admin_api.go#L390)
- **Description:** Initializes a `AdminAPI` instance using an extended public key (xPub).
- **Note:** Requests made with this instance will not be signed.
- **Security Advisory:** For enhanced security, it is strongly recommended to use either `NewAdminAPIWithXPriv`instead, as unsigned requests may be less secure.

**Code snippets:**
- [AdminAPI example](/examples/admin_add_user/admin_add_user.go)
- [UserAPI example](/examples/list_transactions/list_transactions.go)

</details>

<details>
<summary><strong><code>Compatibility & Support</code></strong></summary>

### Compatibility and Support

#### Deprecation Notice
The client **does not support** the following:
- **Admin and non-admin old endpoints** of the SPV Wallet API based on the `/v1/` prefix.
- Deprecated methods for building query parameters for HTTP requests.

#### Current Compatibility
The client is designed for full compatibility with the newer `/api/v1/` endpoints exposed by the SPV Wallet API. It focuses on aligning with the latest standards and structure provided by the API.

</details>

<details>
<summary><strong><code>API Endpoints Compatibility</code></strong></summary>

### API Admin Endpoints Compatibility

#### Access Keys API
| HTTP Method | Endpoint                 | Action             | Support Status | API Code                                                        | Pagination |
|-------------|--------------------------|--------------------|----------------|-----------------------------------------------------------------|------------|
| GET         | /api/v1/admin/users/keys | Search access keys | ✅              | [API](/internal/api/v1/admin/accesskeys/access_keys_api.go#L25) | ✅          |

#### Contacts API
| HTTP Method | Endpoint                             | Action          | Support Status | API Code                                                   | Pagination |
|-------------|--------------------------------------|-----------------|----------------|------------------------------------------------------------|------------|
| GET         | /api/v1/admin/contacts               | Search contacts | ✅              | [API](/internal/api/v1/admin/contacts/contacts_api.go#L42) | ✅          |
| POST        | /api/v1/admin/contacts/confirmations | Confirm contact | ✅              | [API](/internal/api/v1/admin/contacts/contacts_api.go#L83) | ❌          |
| PUT         | /api/v1/admin/contacts/{id}          | Update contact  | ✅              | [API](/internal/api/v1/admin/contacts/contacts_api.go#L68) | ❌          |
| DELETE      | /api/v1/admin/contacts/{id}          | Delete contact  | ✅              | [API](/internal/api/v1/admin/contacts/contacts_api.go#L95) | ❌          |
| POST        | /api/v1/admin/contacts/{paymail}     | Create contact  | ✅              | [API](/internal/api/v1/admin/contacts/contacts_api.go#L27) | ❌          |

#### Invitations API
| HTTP Method | Endpoint                       | Action            | Support Status | API Code                                                         | Pagination |
|-------------|--------------------------------|-------------------|----------------|------------------------------------------------------------------|------------|
| POST        | /api/v1/admin/invitations/{id} | Accept invitation | ✅              | [API](/internal/api/v1/admin/invitations/invitations_api.go#L22) | ❌          |
| DELETE      | /api/v1/admin/invitations/{id} | Reject invitation | ✅              | [API](/internal/api/v1/admin/invitations/invitations_api.go#L35) | ❌          |


#### Paymails API
| HTTP Method | Endpoint                    | Action           | Support Status | API Code                                                   | Pagination |
|-------------|-----------------------------|------------------|----------------|------------------------------------------------------------|------------|
| GET         | /api/v1/admin/paymails      | Search paymails  | ✅              | [API](/internal/api/v1/admin/paymails/paymails_api.go#L73) | ✅          |
| POST        | /api/v1/admin/paymails      | Create paymail   | ✅              | [API](/internal/api/v1/admin/paymails/paymails_api.go#L44) | ❌          |
| GET         | /api/v1/admin/paymails/{id} | Retrieve paymail | ✅              | [API](/internal/api/v1/admin/paymails/paymails_api.go#L59) | ❌          |
| DELETE      | /api/v1/admin/paymails/{id} | Delete paymail   | ✅              | [API](/internal/api/v1/admin/paymails/paymails_api.go#L27) | ❌          |

#### Stats API
| HTTP Method | Endpoint            | Action         | Support Status | API Code                                             | Pagination |
|-------------|---------------------|----------------|----------------|------------------------------------------------------|------------|
| GET         | /api/v1/admin/stats | Retrieve stats | ✅              | [API](/internal/api/v1/admin/stats/stats_api.go#L23) | ✅          |

#### Status API
| HTTP Method | Endpoint             | Action          | Support Status | API Code                                               | Pagination |
|-------------|----------------------|-----------------|----------------|--------------------------------------------------------|------------|
| GET         | /api/v1/admin/status | Retrieve status | ✅              | [API](/internal/api/v1/admin/status/status_api.go#L23) | ❌          |

#### Transactions API
| HTTP Method | Endpoint                        | Action               | Support Status | API Code                                                           | Pagination |
|-------------|---------------------------------|----------------------|----------------|--------------------------------------------------------------------|------------|
| GET         | /api/v1/admin/transactions      | Search transactions  | ✅              | [API](/internal/api/v1/admin/transactions/transactions_api.go#L39) | ✅          |
| GET         | /api/v1/admin/transactions/{id} | Retrieve transaction | ✅              | [API](/internal/api/v1/admin/transactions/transactions_api.go#L26) | ❌          |

#### UTXOs API
| HTTP Method | Endpoint            | Action       | Support Status | API Code                                             | Pagination |
|-------------|---------------------|--------------|----------------|------------------------------------------------------|------------|
| GET         | /api/v1/admin/utxos | Search UTXOs | ✅              | [API](/internal/api/v1/admin/utxos/utxos_api.go#L25) | ✅          |

#### Webhooks API
| HTTP Method | Endpoint                             | Action               | Support Status | API Code                                                   | Pagination |
|-------------|--------------------------------------|----------------------|----------------|------------------------------------------------------------|------------|
| GET         | /api/v1/admin/webhooks/subscriptions | Subscribe to webhook | ✅              | [API](/internal/api/v1/admin/webhooks/webhooks_api.go#L23) | ❌          |
| DELETE      | /api/v1/admin/webhooks/subscriptions | Unsubscribe webhook  | ✅              | [API](/internal/api/v1/admin/webhooks/webhooks_api.go#L36) | ❌          |

#### XPubs API
| HTTP Method | Endpoint            | Action       | Support Status | API Code                                             | Pagination |
|-------------|---------------------|--------------|----------------|------------------------------------------------------|------------|
| GET         | /api/v1/admin/users | Search XPubs | ✅              | [API](/internal/api/v1/admin/xpubs/xpubs_api.go#L41) | ✅          |
| POST        | /api/v1/admin/users | Create XPub  | ✅              | [API](/internal/api/v1/admin/xpubs/xpubs_api.go#L27) | ❌          |

</details>

<details>
<summary><strong><code>API Non-Admin Endpoints Compatibility</code></strong></summary>

### API Non-Admin Endpoints Compatibility

#### Access Keys API
| HTTP Method | Endpoint                        | Action              | Support Status | API Code                                                      | Pagination |
|-------------|---------------------------------|---------------------|----------------|---------------------------------------------------------------|------------|
| GET         | /api/v1/users/current/keys      | Search access keys  | ✅              | [API](/internal/api/v1/user/accesskeys/access_key_api.go#L56) | ✅          |
| POST        | /api/v1/users/current/keys      | Create access key   | ✅              | [API](/internal/api/v1/user/accesskeys/access_key_api.go#L27) | ❌          |
| GET         | /api/v1/users/current/keys/{id} | Retrieve access key | ✅              | [API](/internal/api/v1/user/accesskeys/access_key_api.go#L42) | ❌          |
| DELETE      | /api/v1/users/current/keys/{id} | Revoke access key   | ✅              | [API](/internal/api/v1/user/accesskeys/access_key_api.go#L82) | ❌          |

#### Contacts API
| HTTP Method | Endpoint                   | Action            | Support Status | API Code                                                   | Pagination |
|-------------|----------------------------|-------------------|----------------|------------------------------------------------------------|------------|
| GET         | /api/v1/contacts           | Search contacts   | ✅              | [API](/internal/api/v1/user/contacts/contacts_api.go#L27)  | ✅          |
| GET         | /api/v1/contacts/{paymail} | Retrieve contact  | ✅              | [API](/internal/api/v1/user/contacts/contacts_api.go#L53)  | ❌          |
| PUT         | /api/v1/contacts/{paymail} | Upsert contact    | ✅              | [API](/internal/api/v1/user/contacts/contacts_api.go#L67)  | ❌          |
| DELETE      | /api/v1/contacts/{paymail} | Remove contact    | ✅              | [API](/internal/api/v1/user/contacts/contacts_api.go#L89)  | ❌          |
| POST        | /api/v1/contacts/{paymail} | Confirm contact   | ✅              | [API](/internal/api/v1/user/contacts/contacts_api.go#L101) | ❌          |
| DELETE      | /api/v1/contacts/{paymail} | Unconfirm contact | ✅              | [API](/internal/api/v1/user/contacts/contacts_api.go#L113) | ❌          |

#### Invitations API
| HTTP Method | Endpoint                               | Action            | Support Status | API Code                                                        | Pagination |
|-------------|----------------------------------------|-------------------|----------------|-----------------------------------------------------------------|------------|
| POST        | /api/v1/invitations/{paymail}/contacts | Accept invitation | ✅              | [API](/internal/api/v1/user/invitations/invitations_api.go#L22) | ❌          |
| DELETE      | /api/v1/invitations/{paymail}          | Reject invitation | ✅              | [API](/internal/api/v1/user/invitations/invitations_api.go#L34) | ❌          |

#### Merkle Roots API
| HTTP Method | Endpoint            | Action              | Support Status | API Code                                                        | Pagination |
|-------------|---------------------|---------------------|----------------|-----------------------------------------------------------------|------------|
| GET         | /api/v1/merkleroots | Search Merkle roots | ✅              | [API](/internal/api/v1/user/merkleroots/merkleroots_api.go#L36) | ❌          |

#### Paymails API
| HTTP Method | Endpoint         | Action          | Support Status | API Code                                                  | Pagination |
|-------------|------------------|-----------------|----------------|-----------------------------------------------------------|------------|
| GET         | /api/v1/paymails | Search paymails | ✅              | [API](/internal/api/v1/user/paymails/paymails_api.go#L25) | ✅          |

#### Transactions API
| HTTP Method | Endpoint                    | Action               | Support Status | API Code                                                           | Pagination |
|-------------|-----------------------------|----------------------|----------------|--------------------------------------------------------------------|------------|
| GET         | /api/v1/transactions        | Search transactions  | ✅              | [API](/internal/api/v1/user/transactions/transactions_api.go#L137) | ✅          |
| POST        | /api/v1/transactions        | Record transaction   | ✅              | [API](/internal/api/v1/user/transactions/transactions_api.go#L93)  | ❌          |
| POST        | /api/v1/transactions/drafts | Draft transaction    | ✅              | [API](/internal/api/v1/user/transactions/transactions_api.go#L78)  | ❌          |
| GET         | /api/v1/transactions/{id}   | Retrieve transaction | ✅              | [API](/internal/api/v1/user/transactions/transactions_api.go#L123) | ❌          |
| PATCH       | /api/v1/transactions/{id}   | Update transaction   | ✅              | [API](/internal/api/v1/user/transactions/transactions_api.go#L108) | ❌          |

#### UTXOs API
| HTTP Method | Endpoint      | Action       | Support Status | API Code                                            | Pagination |
|-------------|---------------|--------------|----------------|-----------------------------------------------------|------------|
| GET         | /api/v1/utxos | Search UTXOs | ✅              | [API](/internal/api/v1/user/utxos/utxos_api.go#L25) | ❌          |

#### XPubs API
| HTTP Method | Endpoint              | Action                     | Support Status | API Code                                           | Pagination |
|-------------|-----------------------|----------------------------|----------------|----------------------------------------------------|------------|
| GET         | /api/v1/users/current | Retrieve current user info | ✅              | [API](/internal/api/v1/user/xpubs/xpub_api.go#L24) | ❌          |
| PATCH       | /api/v1/users/current | Update current user info   | ✅              | [API](/internal/api/v1/user/xpubs/xpub_api.go#L24) | ❌          |


</details>


<details>
<summary><strong><code>Development Build Commands</code></strong></summary>
<br/>

Get the [MAGE-X](https://github.com/mrz1836/mage-x) build tool for development:
```shell script
go install github.com/mrz1836/mage-x/cmd/magex@latest
```

View all build commands

```bash script
magex help
```

</details>

<details>
<summary><strong>Repository Features</strong></summary>
<br/>

This repository includes 25+ built-in features covering CI/CD, security, code quality, developer experience, and community tooling.

**[View the full Repository Features list →](.github/docs/repository-features.md)**

</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

This project uses [goreleaser](https://github.com/goreleaser/goreleaser) for streamlined binary and library deployment to GitHub. To get started, install it via:

```bash
brew install goreleaser
```

The release process is defined in the [.goreleaser.yml](.goreleaser.yml) configuration file.


Then create and push a new Git tag using:

```bash
magex version:bump push=true bump=patch branch=main
```

This process ensures consistent, repeatable releases with properly versioned artifacts and citation metadata.

</details>

<details>
<summary><strong><code>Pre-commit Hooks</code></strong></summary>
<br/>

Set up the Go-Pre-commit System to run the same formatting, linting, and tests defined in [AGENTS.md](.github/AGENTS.md) before every commit:

```bash
go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest
go-pre-commit install
```

The system is configured via modular env files in [`.github/env/`](.github/env/README.md) and provides 17x faster execution than traditional Python-based pre-commit hooks. See the [complete documentation](http://github.com/mrz1836/go-pre-commit) for details.

</details>

<details>
<summary><strong>GitHub Workflows</strong></summary>
<br/>

All workflows are driven by modular configuration in [`.github/env/`](.github/env/README.md) — no YAML editing required.

**[View all workflows and the control center →](.github/docs/workflows.md)**

</details>

<details>
<summary><strong><code>Updating Dependencies</code></strong></summary>
<br/>

To update all dependencies (Go modules, linters, and related tools), run:

```bash
magex deps:update
```

This command ensures all dependencies are brought up to date in a single step, including Go modules and any tools managed by [MAGE-X](https://github.com/mrz1836/mage-x). It is the recommended way to keep your development environment and CI in sync with the latest versions.

</details>

<br/>

## 🧪 Examples & Tests

All unit tests and [examples](examples) run via [GitHub Actions](https://github.com/bsv-blockchain/spv-wallet-go-client/actions) and use [Go version 1.26.x](https://go.dev/doc/go1.26). View the [configuration file](.github/workflows/fortress.yml).

Run all tests (fast):

```bash script
magex test
```

Run all tests with race detector (slower):
```bash script
magex test:race
```

<br/>

## ⚡ Benchmarks

Run the Go benchmarks:

```bash script
magex bench
```

<br/>

## 🛠️ Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## 🤖 AI Usage & Assistant Guidelines
Read the [AI Usage & Assistant Guidelines](.github/tech-conventions/ai-compliance.md) for details on how AI is used in this project and how to interact with AI assistants.

<br/>

## 👥 Maintainers
| [<img src="https://github.com/icellan.png" height="50" alt="Siggi" />](https://github.com/icellan) | [<img src="https://github.com/galt-tr.png" height="50" alt="Galt" />](https://github.com/galt-tr) | [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:--------------------------------------------------------------------------------------------------:|:-------------------------------------------------------------------------------------------------:|:------------------------------------------------------------------------------------------------:|
|                                [Siggi](https://github.com/icellan)                                 |                                [Dylan](https://github.com/galt-tr)                                 |                                [MrZ](https://github.com/mrz1836)                                 |

<br/>

## 🤝 Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.

[![Stars](https://img.shields.io/github/stars/bsv-blockchain/spv-wallet-go-client?label=Please%20like%20us&style=social&v=1)](https://github.com/bsv-blockchain/spv-wallet-go-client/stargazers)

<br/>

## 📝 License

[![License](https://img.shields.io/badge/license-OpenBSV-blue?style=flat&logo=springsecurity&logoColor=white)](LICENSE)
