generate:
	lambchaind tx lambchain verify $(PROOF) \
		--from alice \
		--chain-id lambchain \
		--generate-only > transaction.json

sign: transaction.json
	lambchaind tx sign transaction.json \
		--chain-id lambchain \
		--keyring-backend test \
		--from alice \
		--overwrite \
		--output-document signed.json

encode: signed.json
	lambchaind tx encode signed.json > encoded.txt

send-without-encode:
	grpcurl -plaintext \
	    -d "{\"tx_bytes\":\"`cat encoded.txt`\",\"mode\":\"BROADCAST_MODE_SYNC\"}" \
	    localhost:9090 \
	    cosmos.tx.v1beta1.Service/BroadcastTx

send: | encode send-without-encode

# Output result code for HASH transaction
result:
	lambchaind query tx $(HASH) \
		--output json | jq .code

clean:
	rm -f transaction.json signed.json encoded.txt
