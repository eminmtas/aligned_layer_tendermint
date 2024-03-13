#!/bin/bash
NODE_NAME=$1
VALIDATOR=$2
STAKING_AMOUNT=$3

NODE_HOME=$HOME/.alignedlayer
CHAIN_BINARY=alignedlayerd
CHAIN_ID=alignedlayer
PEER_ADDR=blockchain-1

git clone https://github.com/yetanotherco/aligned_layer_tendermint.git
cd aligned_layer_tendermint
ignite chain build 

$CHAIN_BINARY init $NODE_NAME --chain-id $CHAIN_ID
curl $PEER_ADDR:26657/genesis | jq '.result.genesis' > $NODE_HOME/config/genesis.json

NODEID=$(curl -s $PEER_ADDR:26657/status | jq -r '.result.node_info.id')
$CHAIN_BINARY config set config p2p.seeds "$NODEID@$PEER_ADDR:26656" --skip-validate
$CHAIN_BINARY config set config p2p.persistent_peers "$NODEID@$PEER_ADDR:26656" --skip-validate
$CHAIN_BINARY config set app minimum-gas-prices 0.25stake --skip-validate

#ADD ASK FOR TOKENS

x=$($CHAIN_BINARY keys add $VALIDATOR)
NODE_ADDR=$(echo $x | awk '{print $3}')
VALIDATOR_KEY=$($CHAIN_BINARY tendermint show-validator)

cd $NODE_HOME/config
touch validator.json

echo '{"pubkey": '$VALIDATOR_KEY',
	"amount": "'$STAKING_AMOUNT'stake",
	"moniker": "'$VALIDATOR'",
	"commission-rate": "0.1",
	"commission-max-rate": "0.2",
	"commission-max-change-rate": "0.01",
	"min-self-delegation": "1"}' > validator.json

open validator.json
echo $NODE_ADDR

$CHAIN_BINARY tx staking create-validator $HOME/$NODE_HOME/config/validator.json --from $NODE_ADDR --node tcp://$PEER_ADDR:26656

$CHAIN_BINARY start
