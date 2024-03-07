#!/bin/bash

password="password"
token="stake"
initial_balance=1000000000$token
initial_stake = 60000000$token


if [ $# -lt 1 ]; then
    echo "Usage: $0 <node1> [<node2> ...]"
    exit 1
fi

echo "Creating directories for nodes..."
rm -rf prod-sim
for node in "$@"; do
    mkdir -p prod-sim/$node
done

for node in "$@"; do
    echo "Initializing $node..."
    docker run -v $(pwd)/prod-sim/$node:/root/.lambchain -it lambchaind_i init lambchain --chain-id lambchain

    docker run --rm -it -v $(pwd)/prod-sim/$node:/root/.lambchain --entrypoint sed lambchaind_i -i 's/"stake"/"$token"/g' /root/.lambchain/config/genesis.json

done

echo "Creating key for alice in $1..."
printf "$password\n$password\n" | docker run --rm -i -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i keys --keyring-backend file --keyring-dir /root/.lambchain/keys add alice

alice_address=$(echo $password | docker run --rm -i -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i keys --keyring-backend file --keyring-dir /root/.lambchain/keys show alice --address)

echo "Giving alice some tokens..."
docker run --rm -it -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i genesis add-genesis-account $alice_address $initial_balance

echo "Giving alice some stake..."
echo $password | docker run --rm -i -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i genesis gentx alice $initial_stake --keyring-backend file --keyring-dir /root/.lambchain/keys --account-number 0 --sequence 0 --chain-id lambchain --gas 1000000 --gas-prices 0.1upawn

echo "Collecting genesis transactions..."
docker run --rm -it -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i genesis collect-gentxs

if ! docker run --rm -it -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i validate-genesis; then
    echo "Invalid genesis"
    exit 1
fi

echo "Copying genesis file to other nodes..."
for node in "${@:2}"; do
    cp prod-sim/$1/config/genesis.json prod-sim/$node/config/genesis.json
done