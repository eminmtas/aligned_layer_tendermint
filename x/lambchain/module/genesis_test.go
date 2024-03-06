package lambchain_test

import (
	"testing"

	keepertest "lambchain/testutil/keeper"
	"lambchain/testutil/nullify"
	lambchain "lambchain/x/lambchain/module"
	"lambchain/x/lambchain/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LambchainKeeper(t)
	lambchain.InitGenesis(ctx, k, genesisState)
	got := lambchain.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
