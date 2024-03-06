package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "lambchain/testutil/keeper"
	"lambchain/x/lambchain/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.LambchainKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
