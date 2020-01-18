package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
)

func TestSendKeeper(t *testing.T) {
	_, ctx, keeper, _, _, _, _ := SetupTestInput()

	addr1 := sdk.AccAddress(crypto.AddressHash([]byte("addr1")))

	params := keeper.GetParams(ctx)
	require.Equal(t, int8(10), params.MaxDecimal)
	require.True(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000))).IsEqual(params.IssueFee))
	require.True(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000))).IsEqual(params.MintFee))

	keeper.SetMaxDecimal(ctx, 8)
	keeper.SetIssueFee(ctx, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2000000000))))
	keeper.SetMintFee(ctx, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200000000))))
	params = keeper.GetParams(ctx)
	require.Equal(t, int8(8), params.MaxDecimal)
	require.True(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2000000000))).IsEqual(params.IssueFee))
	require.True(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200000000))).IsEqual(params.MintFee))

	iterator := keeper.ListToken(ctx)
	require.False(t, iterator.Valid())

	token := types.NewToken("btc_123", "bitcoin", 6, 21000000000000, false, "bitcoin on barkisnet", addr1)
	keeper.SetToken(ctx, token)

	iterator = keeper.ListToken(ctx)
	require.True(t, iterator.Valid())
	gettedToken := keeper.DecodeToToken(iterator.Value())
	require.Equal(t, "btc_123", gettedToken.Symbol)
	require.Equal(t, "bitcoin", gettedToken.Name)
	iterator.Next()
	require.False(t, iterator.Valid())

	require.True(t, keeper.IsTokenExist(ctx, "btc_123"))
	require.False(t, keeper.IsTokenExist(ctx, "btc_124"))
}
