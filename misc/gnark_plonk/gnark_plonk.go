package main

import (
	"fmt"

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
	X frontend.Variable `gnark:",public"` // x  --> public visibility
	Y frontend.Variable `gnark:",public"` // Y  --> public visibility

	E frontend.Variable
}

// Define declares the circuit logic. The compiler then produces a list of constraints
// y == x**e
func (circuit *Circuit) Define(api frontend.API) error {
	// number of bits of exponent
	const bitSize = 4000

	// specify constraints
	output := frontend.Variable(1)
	bits := api.ToBinary(circuit.E, bitSize)

	for i := 0; i < len(bits); i++ {
		// api.Println(fmt.Sprintf("e[%d]", i), bits[i]) // we may print a variable for testing and / or debugging purposes

		if i != 0 {
			output = api.Mul(output, output)
		}
		multiply := api.Mul(output, circuit.X)
		output = api.Select(bits[len(bits)-1-i], multiply, output)

	}

	api.AssertIsEqual(circuit.Y, output)

	return nil
}

func main() {
	var myCircuit Circuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &myCircuit)

	scs := ccs.(*cs.SparseR1CS)
	kzgsrs, _ := test.NewKZGSRS(scs)

	// Witnesses instantiation. Witness is known only by the prover,
	// while public w is a public data known by the verifier.
	var w Circuit
	w.X = 2
	w.E = 2
	w.Y = 4

	witnessFull, _ := frontend.NewWitness(&w, ecc.BN254.ScalarField())
	witnessPublic, _ := frontend.NewWitness(&w, ecc.BN254.ScalarField(), frontend.PublicOnly())

	// public data consists of the polynomials describing the constants involved
	// in the constraints, the polynomial describing the permutation ("grand
	// product argument"), and the FFT domains.

	pk, vk, _ := plonk.Setup(ccs, kzgsrs)

	proof, _ := plonk.Prove(ccs, pk, witnessFull)

	err := plonk.Verify(proof, vk, witnessPublic)

	if err != nil {
		fmt.Print("Invalid Proof")
		return
	}
	fmt.Print("ValidProof")

}
