package asset

import (
	"bytes"
	"fmt"
	"strings"

	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
	"github.com/barkisnet/barkis/x/auth"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case IssueMsg:
			return handleIssueMsg(ctx, k, msg)

		case MintMsg:
			return handleMintMsg(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleIssueMsg(ctx sdk.Context, k Keeper, msg IssueMsg) sdk.Result {
	maxDecimal := k.GetMaxDecimal(ctx)
	if msg.Decimal > maxDecimal {
		return types.ErrInvalidDecimal(types.DefaultCodespace, fmt.Sprintf("token decimal should not greater than %d", maxDecimal)).Result()
	}
	if k.IsTokenExist(ctx, strings.ToLower(msg.Symbol)) {
		return types.ErrInvalidTokenSymbol(types.DefaultCodespace, fmt.Sprintf("duplicated token symbol: %s", strings.ToLower(msg.Symbol))).Result()
	}

	token := types.NewToken(strings.ToLower(msg.Symbol), msg.Name, msg.Decimal, msg.TotalSupply, msg.Mintable, msg.Description, msg.From)
	k.SetToken(ctx, token)

	issueFee := k.GetIssueFee(ctx)
	err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.From, auth.FeeCollectorName, issueFee)
	if err != nil {
		return err.Result()
	}

	mintedToken := sdk.Coins{sdk.NewCoin(token.Symbol, sdk.NewInt(token.TotalSupply))}

	err = k.SupplyKeeper.MintCoins(ctx, types.ModuleName, mintedToken)
	if err != nil {
		return err.Result()
	}

	err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, token.Owner, mintedToken)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.EventTypeIssueToken, mintedToken.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMintMsg(ctx sdk.Context, k Keeper, msg MintMsg) sdk.Result {
	token := k.GetToken(ctx, msg.Symbol)
	if token == nil {
		return types.ErrInvalidTokenSymbol(types.DefaultCodespace, fmt.Sprintf("token %s is not exist", msg.Symbol)).Result()
	}
	if !token.Mintable {
		return types.ErrNotMintableToken(types.DefaultCodespace, fmt.Sprintf("token %s is not mintable", token.Symbol)).Result()
	}
	if !bytes.Equal(token.Owner, msg.From) {
		return types.ErrUnauthorizedMint(types.DefaultCodespace, fmt.Sprintf("only %s is authorized to mint token %s", token.Owner.String(), token.Symbol)).Result()
	}
	possibleMintAmount := types.MaxTotalSupply - token.TotalSupply
	if msg.Amount > possibleMintAmount {
		return types.ErrInvalidMintAmount(types.DefaultCodespace, fmt.Sprintf("minted too many token, maximum possible minted amount %d, actual minted amount %d", possibleMintAmount, msg.Amount)).Result()
	}

	mintFee := k.GetMintFee(ctx)
	err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.From, auth.FeeCollectorName, mintFee)
	if err != nil {
		return err.Result()
	}

	token.TotalSupply = msg.Amount + token.TotalSupply
	k.UpdateToken(ctx, token)

	mintedToken := sdk.Coins{sdk.NewCoin(token.Symbol, sdk.NewInt(msg.Amount))}
	err = k.SupplyKeeper.MintCoins(ctx, types.ModuleName, mintedToken)
	if err != nil {
		return err.Result()
	}

	err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, token.Owner, mintedToken)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.EventTypeMintToken, mintedToken.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}
