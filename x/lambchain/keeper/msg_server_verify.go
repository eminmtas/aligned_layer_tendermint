package keeper

import (
	"context"

	"lambchain/x/lambchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Verify(goCtx context.Context, msg *types.MsgVerify) (*types.MsgVerifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgVerifyResponse{}, nil
}
