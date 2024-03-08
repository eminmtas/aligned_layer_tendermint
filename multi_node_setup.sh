#!/bin/bash

password="password"
token="stake"
initial_balance=1000000000
initial_stake=60000000


if [ $# -lt 1 ]; then
    echo "Usage: $0 <node1> [<node2> ...]"
    exit 1
fi

echo "Creating directories for nodes..."
rm -rf prod-sim
for node in "$@"; do
    mkdir -p prod-sim/$node
done

node_ids=()

for node in "$@"; do
    echo "Initializing $node..."
    docker run -v $(pwd)/prod-sim/$node:/root/.lambchain -it lambchaind_i init lambchain --chain-id lambchain > /dev/null

    docker run --rm -it -v $(pwd)/prod-sim/$node:/root/.lambchain --entrypoint sed lambchaind_i -i 's/"stake"/"'$token'"/g' /root/.lambchain/config/genesis.json

    node_id=$(docker run --rm -i -v $(pwd)/prod-sim/$node:/root/.lambchain lambchaind_i tendermint show-node-id)
    node_ids+=($node_id)

    echo "Node ID for $node: $node_id"
done

echo "Creating key for alice in $1..."
printf "$password\n$password\n" | docker run --rm -i -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i keys --keyring-backend file --keyring-dir /root/.lambchain/keys add alice > /dev/null

alice_address=$(echo $password | docker run --rm -i -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i keys --keyring-backend file --keyring-dir /root/.lambchain/keys show alice --address)

echo "Alice's address: $alice_address"

echo "Giving alice some tokens..."
docker run --rm -it -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i genesis add-genesis-account $alice_address $initial_balance$token

echo "Giving alice some stake..."
echo $password | docker run --rm -i -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i genesis gentx alice $initial_stake$token --keyring-backend file --keyring-dir /root/.lambchain/keys --account-number 0 --sequence 0 --chain-id lambchain --gas 1000000 --gas-prices 0.1$token

echo "Collecting genesis transactions..."
docker run --rm -it -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i genesis collect-gentxs > /dev/null

if ! docker run --rm -it -v $(pwd)/prod-sim/$1:/root/.lambchain lambchaind_i genesis validate-genesis; then
    echo "Invalid genesis"
    exit 1
fi

echo "Copying genesis file to other nodes..."
for node in "${@:2}"; do
    cp prod-sim/$1/config/genesis.json prod-sim/$node/config/genesis.json
done

echo "Setting node addresses in config..."
for (( i=1; i <= "$#"; i++ )); do
    other_addresses=()
    for (( j=1; j <= "$#"; j++ )); do
        if [ $j -ne $i ]; then
            other_addresses+=("${node_ids[$j - 1]}@${!j}:26656")
        fi
    done
    other_addresses=$(IFS=,; echo "${other_addresses[*]}")
    sed -i .back 's/persistent_peers = ""/persistent_peers = "'$other_addresses'"/' prod-sim/${!i}/config/config.toml
done