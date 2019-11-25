package client

import (
	"github.com/gorilla/mux"

	"github.com/barkisnet/barkis/client/context"
)

// Register routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	RegisteKeysRoutes(r, true)
	RegisterRPCRoutes(cliCtx, r)
}
