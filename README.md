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

# How It Works

A blockchain can be created using the igniteCLI, which collects various files from CosmosSDK, making it easier to deploy a blockchain to production. CosmosSDK manages communication between the Application Layer and the Consensus Layer via the ABCI. The Consensus Layer, in turn, is responsible for CometBFT, which uses Tendermint as the consensus client.

In the Application Layer, messages are defined to be included in transactions. The application then forwards the message's information to a Module, which is essentially a Go function triggered when the message is included in a transaction.

Abstracting these concepts, we have the following diagram: a person sends a transaction to the blockchain, the blockchain processes it, and redirects the information to the appropriate module based on the message's information. The module we aim to implement is a Snarks validator, so the transaction needs to include the proof, the validator key, and the public witness.

<p align="center">
  <img src="imgs/Diagram_Cosmos.svg">
</p>
