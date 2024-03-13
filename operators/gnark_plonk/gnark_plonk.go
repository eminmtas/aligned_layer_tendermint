package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"

	"github.com/consensys/gnark/frontend"
)

// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
// The witness has a
//   - secret part --> known to the prover only
//   - public part --> known to the prover and the verifier
type Circuit struct {
	X frontend.Variable `gnark:"x"`       // x  --> secret visibility (default)
	Y frontend.Variable `gnark:",public"` // Y  --> public visibility
}

// Define declares the circuit logic. The compiler then produces a list of constraints
// which must be satisfied (valid witness) in order to create a valid zk-SNARK
func (circuit *Circuit) Define(api frontend.API) error {
	// compute x**3 and store it in the local variable x3.
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)

	// compute x**3 + x + 5 and store it in the local variable res
	res := api.Add(x3, circuit.X, 5)

	// assert that the statement x**3 + x + 5 == y is true.
	api.AssertIsEqual(circuit.Y, res)
	return nil
}

// Defines the proof that will be proved.
func toProve() Circuit {
	return Circuit{
		X: 3,
		Y: 35,
	}
}

func main() {
	var myCircuit Circuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &myCircuit)

	scs := ccs.(*cs.SparseR1CS)
	kzgsrs, _ := test.NewKZGSRS(scs)
	pk, vk, _ := plonk.Setup(ccs, kzgsrs)

	circuit := toProve()
	fullWitness, _ := frontend.NewWitness(&circuit, ecc.BN254.ScalarField())
	publicWitness, _ := fullWitness.Public()

	proof, _ := plonk.Prove(ccs, pk, fullWitness)

	serialize(proof, "proof.base64")
	serialize(publicWitness, "public_witness.base64")
	serialize(vk, "verifying_key.base64")

}

func serialize[w io.WriterTo](src w, name string) {
	var buffer bytes.Buffer
	src.WriteTo(&buffer)

	inner := buffer.Bytes()

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(inner)))
	base64.StdEncoding.Encode(encoded, inner)

	os.WriteFile(name, encoded, 0644)
}
