# Documentation

Several approaches exist to release a blockchain using the Cosmos ecosystem:

- [tendermint-rs](https://github.com/informalsystems/tendermint-rs)
    - Provides client libraries for Tendermint/CometBFT, which are consensus clients, along with some network settings.
- [cosmos-sdk](https://github.com/cosmos/cosmos-sdk)
    - A comprehensive framework that supplies all the code necessary to set up a blockchain; only "applications" need to be created.
- [ignite-cli](https://github.com/ignite/cli)
    - Generates all the boilerplate code required to build and deploy a chain using CosmosSDK.

## Tendermint-rs

As it is solely a library, establishing communication between the consensus client and the application itself (i.e., determining what to execute) is necessary. Managing keys is another critical aspect that may need to be developed from scratch.

An advantage of tendermint-rs is that it offers more control to the developer, albeit at the cost of speed.

## CosmosSDK

This framework contains all the necessary components. However, some logic is required to interact with transactions. Setting up multiple nodes can be tedious, as can handling keys. The documentation provides explanations for various approaches.

Some parts or concepts are not fully customizable, which is a trade-off for speed.

## Ignite CLI

By harnessing the power of CosmosSDK, ignite provides useful commands for creating multiple nodes.

It can even set up a chain, whose source code is in a GitHub repository, with just one command.
