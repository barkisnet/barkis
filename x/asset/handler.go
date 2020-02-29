package asset

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

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

		case DelayedTransferMsg:
			return handleDelayedTransferMsg(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Called every block, handler matured delayed transfer
func EndBlocker(ctx sdk.Context, k Keeper) []abci.ValidatorUpdate {
	iterator := k.ListDelayedTransferMaturedTime(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		if len(key) != 17 {
			continue
		}

		matureTime := int64(binary.BigEndian.Uint64(key[1:9]))
		if matureTime > ctx.BlockTime().Unix() {
			break
		}

		sequenceBytes := iterator.Value()
		sequence := int64(binary.BigEndian.Uint64(sequenceBytes))

		delayedTransfer := k.GetDelayedTransfer(ctx, sequence)
		err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delayedTransfer.To, delayedTransfer.Amount)
		if err != nil {
			ctx.Logger().Error("failed to process matured delayed transfer", "from", delayedTransfer.From.String(), "to", delayedTransfer.To.String(),
				"amount", delayedTransfer.Amount.String(), "maturedTime", delayedTransfer.MaturedTime)
			continue
		}

		k.DeleteDelayedTransfer(ctx, delayedTransfer)
	}

	return nil
}

func handleIssueMsg(ctx sdk.Context, k Keeper, msg IssueMsg) sdk.Result {
	maxDecimal := k.GetMaxDecimal(ctx)
	if msg.Decimal > maxDecimal {
		return types.ErrInvalidDecimal(types.DefaultCodespace, fmt.Sprintf("token decimal should not greater than %d", maxDecimal)).Result()
	}
	if k.IsTokenExist(ctx, strings.ToLower(msg.Symbol)) {
		return types.ErrInvalidTokenSymbol(types.DefaultCodespace, fmt.Sprintf("duplicated token symbol: %s", strings.ToLower(msg.Symbol))).Result()
	}

	token := NewToken(strings.ToLower(msg.Symbol), msg.Name, msg.Decimal, msg.TotalSupply, msg.Mintable, msg.Description, msg.From)
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

func handleDelayedTransferMsg(ctx sdk.Context, k Keeper, msg DelayedTransferMsg) sdk.Result {
	err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.From, ModuleName, msg.Amount)
	if err != nil {
		return err.Result()
	}

	maturedTime := ctx.BlockTime().Unix() + msg.DelayedPeriod
	sequence := k.GetSequence(ctx)
	delayedTransfer := NewDelayedTransfer(msg.From, msg.To, msg.Amount, maturedTime, sequence)
	k.InsertDelayedTransfer(ctx, delayedTransfer)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.EventTypeDelayedTransfer, msg.Amount.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}
