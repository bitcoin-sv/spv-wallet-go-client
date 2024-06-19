# Quick Guide how to run examples

In this directory you can find examples of how to use the `spv-wallet-go-client` package.

## Before you run

### Pre-requisites

-   You have access to the `spv-wallet` non-custodial wallet (running locally or remotely).
-   You have installed this package on your machine (`go install` on this project's root directory).

### Concerning the keys

-   The `ExampleAdminKey` defined in `example_keys.go` is the default one from [spv-wallet-web-backend repository](https://github.com/bitcoin-sv/spv-wallet-web-backend/blob/main/config/viper.go#L56)
    -   If in your current `spv-wallet` instance you have a different `adminKey`, you should replace the one in `example_keys` with the one you have.
-   The `ExampleXPub` and `ExampleXPriv` are just placeholders, which won't work.
    -   You should replace them by newly generated ones using `task generate_keys`,
    -   ... or use your actual keys if you have them (don't use the keys which are already added to another wallet).

> Additionally, to make it work properly, you should adjust the `ExamplePaymail` to align with your `domains` configuration in the `spv-wallet` instance.

## Proposed order of executing examples

1. `generate_keys` - generates new keys (you can copy them to `example_keys` if you want to use them in next examples)
2. `admin_add_user` - adds a new user (more precisely adds `ExampleXPub` and then `ExamplePaymail` to the wallet)

> To fully experience the next steps, it would be beneficial to transfer some funds to your `ExamplePaymail`. This ensures the examples run smoothly by demonstrating the creation of a transaction with an actual balance. You can transfer funds to your `ExamplePaymail` using a Bitcoin SV wallet application such as HandCash or any other that supports Paymail.

3. `get_balance` - checks the balance - if you've transferred funds to your `ExamplePaymail`, you should see them here
4. `create_transaction` - creates a transaction (you can adjust the `outputs` to your needs)
5. `list_transactions` - lists all transactions and with example filtering
6. `send_op_return` - sends an OP_RETURN transaction
7. `admin_remove_user` - removes the user

In addition to the above, there are additional examples showing how to use the client from a developer perspective:

-   `handle_exceptions` - presents how to "catch" exceptions which the client can throw
-   `custom_logger` - shows different ways you can configure (or disable) internal logger

## Util examples

1. `xpriv_from_mnemonic` - allows you to generate/extract an xPriv key from a mnemonic phrase. To you use it you just need to replace the `mnemonic` variable with your own mnemonic phrase.
2. `xpub_from_xpriv` - allows you to generate an xPub key from an xPriv key. To you use it you just need to replace the `xPriv` variable with your own xPriv key.
3. `generate_totp` - allows you to generate and check validity of a TOTP code for client xPriv and a contact's PKI

## How to run an example

The examples are written in Go and can be run by:

```bash
cd examples
task name_of_the_example
```

> See the `examples/Taskfile.yml` for the list of available examples and scripts
