package keeper

import (
	"bytes"
	"context"
	"encoding/base64"
	"strconv"

	"alignedlayer/x/verification/types"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Verify(goCtx context.Context, msg *types.MsgVerify) (*types.MsgVerifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	result := verify(msg.Proof, msg.PublicInputs)
	event := sdk.NewEvent("verification_finished",
		sdk.NewAttribute("proof_verifies", strconv.FormatBool(result)))

	ctx.EventManager().EmitEvent(event)

	return &types.MsgVerifyResponse{}, nil
}

func verify(proof string, pw string) bool {
	// Deserialize Proof
	proof_des := plonk.NewProof(ecc.BN254)
	proof_base64 := []byte(proof)
	proof_bytes := make([]byte, base64.StdEncoding.DecodedLen(len(proof_base64)))
	_, err := base64.StdEncoding.Decode(proof_bytes, proof_base64)
	if err != nil {
		return false
	}
	proof_reader := bytes.NewReader(proof_bytes)
	proof_des.ReadFrom(proof_reader)

	// Deserialize VK
	vk_base64 := []byte("AAAAAAAAEAAwYUgt+gOND7W0wLImGUBHomFlCfUx1Po6zbd0lsEAAQkx1ZbeL9EPAd3Qc/1akKl28WnHbwObuRxHdXIAQtQ6AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABdy0mQ14gV1EStfBPtsFFkmaiOGAOwcSwaFd+Di7k/ge4BkEEQGkzMkdmSU8debcrr5rnmuDk1vIPnuzsVkyEdrEfzxjAU6hwqhcxpSgTJF79SWIvgjeSCObDMic4ucD++sh1cF1ULuZjYmFwDKqt0+rGyTKhdOdpE5e9Bx0PPCE2s8bGKvkU4piZ7MHZ3RQ3e3UwelnapAX2Ge7iDPFcRuMjm84yvut2F6lpCdbT66ntUv2/kAthcnk+qiNSn48ndOeqCVr8CnLg0PhgZOKnu7IXiy/lTylD1WwQz/U+lhzr9I8cahl6BNzAwzthw17RvdXc/IqgxjYLqsLJ9mhf68AAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGZjpOTkg1IOnJgv7cx+10l8apJMzWp5xKX5IW3rvMSwhgA3u8SHx52QmoAZl5cRHlnQyLU917a3UbevVzZkvbtgjDx4kyAPcVBYE85V+nvc9OSStU9OT1YIAG+poYKjSIJWlw4dDYZlOEYwQYY+QxDnjlJNkiO2zGP9mRt+yRUVwAAAAA=")
	vk_des := plonk.NewVerifyingKey(ecc.BN254)
	vk_bytes := make([]byte, base64.StdEncoding.DecodedLen(len(vk_base64)))
	base64.StdEncoding.Decode(vk_bytes, vk_base64)
	vk_reader := bytes.NewReader(vk_bytes)
	vk_des.ReadFrom(vk_reader)

	// Deserialize PublicWitness
	pw_des, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return false
	}
	pw_base64 := []byte(pw)
	pw_bytes := make([]byte, base64.StdEncoding.DecodedLen(len(pw_base64)))
	base64.StdEncoding.Decode(pw_bytes, pw_base64)
	pw_des.UnmarshalBinary(pw_bytes)

	err = plonk.Verify(proof_des, vk_des, pw_des)

	return err == nil
}
