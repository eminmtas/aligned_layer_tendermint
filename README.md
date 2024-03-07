# Bootcamp Verifying Lambchain

This repository contains a WIP zkSNARK verifier blockchain using Cosmos SDK and CometBFT and created with Ignite CLI.

The application interacts with zkSNARK verifiers built in Rust through FFI.

## Requirements

- Go
- Ignite

## Single Node Usage

To run a single node blockchain, run:

```sh
ignite chain serve
```

This command installs dependencies, builds, initializes, and starts your blockchain in development.

To send verify message (transaction), run:

```sh
lambchaind tx lambchain verify --from alice --chain-id lambchain <proof>
```

To get the transaction result, run:

```sh
lambchaind query tx <txhash>
```

## Configure

The blockchain in development can be configured with `config.yml`.

## How It Works

A blockchain can be created using the ignite CLI, which generates boilerplate for a Cosmos SDK application, making it easier to deploy a blockchain to production. Cosmos SDK is built on top of the consensus layer, implementing the ABCI (Application BlockChain Interface). By default, CometBFT (a fork of Tendermint) is used as the consensus layer.

Transaction's information is sent to a specific Cosmos' Module, which in our case is a zk-SNARK verifier.

<p align="center">
  <img src="imgs/Diagram_Cosmos.svg">
</p>

### Transaction Lifecycle

A transaction can be created with ignite CLI, using the following command:

```sh
lambchaind tx lambchain verify --from alice --chain-id --generate-only lambchain "base64-encoded proof"
```

The transaction is structured as a JSON containing metadata and a set of **messages**. A message contains the fully-qualified name of the method that should handle it, and it's parameters.

Custom modules define a message service (through a proto file) containing handlers which are then registered in the app builder. Each module also defines a `Keeper` which encapsulates module state, tipically through a `KvStore`.

```json
{
    "body": {
        "messages": [
            {
                "@type": "/lambchain.lambchain.MsgVerify",
                "creator": "cosmos1524vzjchy064rr98d2de7u6uvl4qr3egfq67xn",
                "proof": "base64-encoded proof"
            }
        ],
        "memo": "",
        "timeout_height": "0",
        "extension_options": [],
        "non_critical_extension_options": []
    },
    "auth_info": {
        "signer_infos": [],
        "fee": {
            "amount": [],
            "gas_limit": "200000",
            "payer": "",
            "granter": ""
        },
        "tip": null
    },
    "signatures": []
}
```

After encoding it with protobuf, the transaction is sent to the node (cometBFT). The transaction is then relayed to the application through the ABCI methods `checkTx` and `deliverTx`.

- In `checkTx`, the default `BaseApp` implementation checks that a handler exists for every message and the `AnteHandler`'s are executed (by default verifying transaction authentication and gas fees).

- In `deliverTx`, in addition to the procedure previously mentioned, the handler is called for every message with it's corresponding parameters, and the `PostHandler`'s are executed.

The response is then encoded in the transaction result, and added to the
blockchain.
