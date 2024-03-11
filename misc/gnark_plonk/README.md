# Using PLONK with gnark

This example demonstrates how to use gnark.go to run a [PLONK](https://eprint.iacr.org/archive/2019/953/20240223:124519) proof. It follows the documentation to instantiate the circuit and verify it.

To run only the verification step, you'll need the proof and public witness. You may need to serialize these components beforehand. The verifying key can be hardcoded into the program if the prover and the verifier must agree and utilize the same setup.

[Introduction to PLONK](https://medium.com/@lucafra92/what-is-plonk-29c56f326cf6)

## How to Run It

```sh
go run gnark_plonk.go
```

## How Verify a Proof

Protocol:
- Setup(Just Once, "circuit agnostic"):
    - Generate public_key and verifying key with the Circuit's Constraint System(specific to the problem) &rarr; pk, vk = setup(ccs, plonk_specific_data(KZGSRS))
- Prover:
    - Calculate &rarr; witnessFull(inputs) & witnessPublic(inputs)
    - Generate &rarr; proof(ccs, pk, witness_full)
- Verifier:
    - Verify &rarr; proof(proof, vk, witnessPublic)

Following this scheme, if the verifier waits for a proof to verify, the prover and the verifier must agree and use the same setup. Then the prover generates the witness and proof and sends it to the verifier, which can already have the vk.





