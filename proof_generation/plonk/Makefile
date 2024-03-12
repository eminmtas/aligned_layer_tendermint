.PHONY: plonk-test plonk-bad-test query-tx

CHAIN ?= alignedlayer
MODULE ?= verification
MESSAGE ?= verify
HASH ?= placeholder
PROOF_FILE ?= example_base64_serialized/proof.bin
PUBLIC_INPUTS_FILE ?= example_base64_serialized/pw.bin

plonk-test:
	$(CHAIN)d tx $(MODULE) $(MESSAGE) \
		$$(cat $(PROOF_FILE)) $$(cat $(PUBLIC_INPUTS_FILE)) \
		--from alice \
		--chain-id $(CHAIN)

plonk-bad-test:
	$(CHAIN)d tx $(MODULE) $(MESSAGE) \
		$$(cat $(PROOF_FILE)) produce_tx_with_code_err \
		--from alice \
		--chain-id $(CHAIN)

query-tx:
	$(CHAIN)d q tx $(HASH)
