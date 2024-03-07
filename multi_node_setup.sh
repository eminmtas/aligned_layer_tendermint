#!/bin/bash

if [ $# -lt 1 ]; then
    echo "Usage: $0 <node1> [<node2> ...]"
    exit 1
fi

echo "Creating directories for nodes..."
rm -rf prod-sim
for node in "$@"; do
    mkdir -p prod-sim/$node
done

echo "Initializing chains..."
for node in "$@"; do
    echo "Initializing $node..."
    docker run -v $(pwd)/prod-sim/$node:/root/.lambchain -it lambchaind_i init lambchain --chain-id lambchain
done

echo "Creating keys..."
echo "password\npassword" > docker run --rm -it -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i keys --keyring-backend file --keyring-dir /root/.lambchain/keys add alice