package rest

import (
	"github.com/gorilla/mux"

	"github.com/barkisnet/barkis/client/context"
	"github.com/barkisnet/barkis/x/asset/internal/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/asset/issue", IssueRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/asset/mint", MintRequestHandlerFn(cliCtx)).Methods("POST")

	r.HandleFunc("/asset/get/{symbol}", getHandlerFn(cliCtx, types.QuerierRoute)).Methods("GET")
	r.HandleFunc("/asset/list", listHandlerFn(cliCtx, types.QuerierRoute)).Methods("GET")
	r.HandleFunc("/asset/params", paramsHandlerFn(cliCtx, types.QuerierRoute)).Methods("GET")
}
