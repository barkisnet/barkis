package asset

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/keeper"
	"github.com/barkisnet/barkis/x/asset/internal/types"
)

func TestSendKeeper(t *testing.T) {
	_, ctx, assetKeeper, _, _, supplyKeeper, _ := keeper.SetupTestInput()

	addr1 := sdk.AccAddress(crypto.AddressHash([]byte("addr1")))
	addr2 := sdk.AccAddress(crypto.AddressHash([]byte("addr2")))

	handler := NewHandler(assetKeeper)

	issueMsg := types.NewIssueMsg(addr1, "bitcoin", "btc", 21000000000000, false, 6, "bitcoin on barkisnet")

	ctx = ctx.WithTxBytes([]byte("123"))
	result := handler(ctx, issueMsg)
	require.Equal(t, sdk.CodeOK, result.Code, result.Log)
	btcSymbol := string(result.Data)

	mintMsg := types.NewMintMsg(addr1, "btc_123", 1000)
	result = handler(ctx, mintMsg)
	require.Equal(t, types.CodeInvalidTokenSymbol, result.Code, result.Log)

	mintMsg = types.NewMintMsg(addr1, btcSymbol, 1000)
	result = handler(ctx, mintMsg)
	require.Equal(t, types.CodeNotMintableToken, result.Code, result.Log)

	issueMsg = types.NewIssueMsg(addr1, "ethereum", "eth", 100000000000000, true, 6, "ethereum on barkisnet")
	ctx = ctx.WithTxBytes([]byte("123456789"))
	result = handler(ctx, issueMsg)
	require.Equal(t, sdk.CodeOK, result.Code, result.Log)
	ethSymbol := string(result.Data)

	mintMsg = types.NewMintMsg(addr2, ethSymbol, 10000)
	result = handler(ctx, mintMsg)
	require.Equal(t, types.CodeUnauthorizedMint, result.Code, result.Log)


	mintMsg = types.NewMintMsg(addr1, ethSymbol, types.MaxTotalSupply)
	result = handler(ctx, mintMsg)
	require.Equal(t, types.CodeInvalidMintAmount, result.Code, result.Log)

	mintMsg = types.NewMintMsg(addr1, ethSymbol, 100000000000000)
	result = handler(ctx, mintMsg)
	require.Equal(t, sdk.CodeOK, result.Code, result.Log)

	expectTotalSupply := sdk.Coins{sdk.NewCoin(btcSymbol, sdk.NewInt(21000000000000)), sdk.NewCoin(ethSymbol, sdk.NewInt(200000000000000)), sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20000000000))}
	require.True(t, expectTotalSupply.IsEqual(supplyKeeper.GetSupply(ctx).GetTotal()), expectTotalSupply.String())
}
