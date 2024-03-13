# Gnark Plonk

The Gnark Plonk verifier in the blockchain needs the following base64 encoded elements:

- Proof
- Public Inputs
- Verifying Key

## Using Gnark Plonk Verifier

Change the circuit definition and solution inside `gnark_plonk.go`:

Generate the proof and necessary elements for the verification by running:

```sh
make generate-proof
```

This will generate the necessary files in the current directory.

Send the proof to the blockchain:

```sh
make send-proof
```

This will output the transaction hash, which we an query by running:

```sh
HASH=63a... make query-tx
```

We should see an event called `verifiaction_finished` containing a `proof_verifies` attribute.

