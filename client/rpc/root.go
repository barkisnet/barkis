package rpc

import (
	"github.com/gorilla/mux"

	"github.com/barkisnet/barkis/client/context"
)

// Register REST endpoints
func RegisterRPCRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/node_info", NodeInfoRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/syncing", NodeSyncingRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks/latest", LatestBlockRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks/{height}", BlockRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks_results/latest", LatestBlockResultsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks_results/{height}", BlockResultsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/validatorsets/latest", LatestValidatorSetRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/validatorsets/{height}", ValidatorSetRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/simulate/{txBytes}", TxSimulateRequestHandlerFn(cliCtx)).Methods("GET")
}
