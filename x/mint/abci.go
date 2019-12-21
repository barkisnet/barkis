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
/*
	// recalculate inflation rate
	totalStakingSupply := k.StakingTokenSupply(ctx)
	bondedRatio := k.BondedRatio(ctx)
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	k.SetMinter(ctx, minter)

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)
*/
	if ctx.BlockHeight() == 1 {
		mintedCoins := sdk.NewCoins(sdk.NewCoin(params.MintDenom, sdk.NewIntWithDecimal(125, 12)))
		err := k.MintCoins(ctx, mintedCoins)
		if err != nil {
			panic(err)
		}
		minter.RemainedTokens = mintedCoins
	}
	var unfreezenTokens sdk.Coins
	if sdk.GlobalUpgradeMgr.IsUpgradeHeight(sdk.ChangeReward) {
		unfreezenTokens = sdk.NewCoins(sdk.NewCoin(params.MintDenom, sdk.NewIntWithDecimal(37, 4))) // 0.37barkis
	} else {
		unfreezenTokens = sdk.NewCoins(sdk.NewCoin(params.MintDenom, sdk.NewIntWithDecimal(5567, 2))) // 0.5567barkis
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
