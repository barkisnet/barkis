package asset

import (
	"fmt"

	sdk "github.com/barkisnet/barkis/types"
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

func handleIssueMsg(ctx sdk.Context, k Keeper, msg sdk.Msg) sdk.Result {
	return sdk.Result{}
}

func handleMintMsg(ctx sdk.Context, k Keeper, msg sdk.Msg) sdk.Result {
	return sdk.Result{}
}
