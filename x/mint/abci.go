package mint

import (
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/mint/internal/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)
	if ctx.BlockHeight() == 1 {
		mintedCoins := sdk.NewCoins(sdk.NewCoin(params.MintDenom, sdk.NewIntWithDecimal(125, 12)))
		err := k.MintCoins(ctx, mintedCoins)
		if err != nil {
			panic(err)
		}
		minter.RemainedTokens = mintedCoins
	}
	var unfreezenTokens sdk.Coins
	if sdk.GlobalUpgradeMgr.IsUpgradeApplied(sdk.RewardUpgrade) {
		updatedParams := k.GetUpdatedParams(ctx)
		unfreezenTokens = sdk.NewCoins(sdk.NewCoin(updatedParams.MintDenom, sdk.NewInt(updatedParams.UnfreezeAmountPerBlock)))
	} else {
		unfreezenTokens = sdk.NewCoins(sdk.NewCoin(params.MintDenom, sdk.NewIntWithDecimal(5567, 2)))
	}
	if minter.RemainedTokens.IsAllGTE(unfreezenTokens) {
		// send the minted coins to the fee collector account
		err := k.AddCollectedFees(ctx, unfreezenTokens)
		if err != nil {
			panic(err)
		}
		minter.RemainedTokens = minter.RemainedTokens.Sub(unfreezenTokens)
		k.SetMinter(ctx, minter)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMint,
				sdk.NewAttribute(types.AttributeKeyRemainedTokens, minter.RemainedTokens.String()),
				sdk.NewAttribute(types.AttributeKeyUnfreezenTokens, unfreezenTokens.String()),
			),
		)
	}
}
