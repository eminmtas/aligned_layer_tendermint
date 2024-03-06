package keeper

import (
	"lambchain/x/lambchain/types"
)

var _ types.QueryServer = Keeper{}
