package keeper

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"

	"alignedlayer/x/verification/types"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Verify(goCtx context.Context, msg *types.MsgVerify) (*types.MsgVerifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	result := verify(msg)
	event := sdk.NewEvent("verification_finished",
		sdk.NewAttribute("proof_verifies", strconv.FormatBool(result)))

	ctx.EventManager().EmitEvent(event)

	return &types.MsgVerifyResponse{}, nil
}

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

func get_proof() (plonk.Proof, plonk.VerifyingKey, constraint.ConstraintSystem, kzg.SRS) {
	var myCircuit Circuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &myCircuit)

	scs := ccs.(*cs.SparseR1CS)
	kzgsrs, _ := test.NewKZGSRS(scs)

	fmt.Println(scs.Coefficients)

	assignment := toProve()

	fullWitness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())

	pk, vk, _ := plonk.Setup(ccs, kzgsrs)

	var vk_buffer bytes.Buffer
	vk.WriteTo(&vk_buffer)

	proof, _ := plonk.Prove(ccs, pk, fullWitness)

	return proof, vk, ccs, kzgsrs
}

func verify(msg *types.MsgVerify) bool {
	proof := plonk.NewProof(ecc.BN254)
	deserialize(proof, msg.Proof)

	public_input, _ := witness.New(ecc.BN254.ScalarField())
	deserialize(public_input, msg.PublicInputs)

	verifying_key := plonk.NewVerifyingKey(ecc.BN254)
	deserialize(verifying_key, msg.ConstraintSystem)

	err := plonk.Verify(proof, verifying_key, public_input)
	if err != nil {
		fmt.Println("NO VERIFICA: ", err.Error())
	} else {
		fmt.Println("VERIFICA")
	}

	return err == nil
}

func deserialize[r io.ReaderFrom](dst r, encoded string) error {
	bytes_buffer, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytes_buffer)
	_, err = dst.ReadFrom(reader)

	return err
}
