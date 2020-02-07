package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/barkisnet/barkis/client/context"
	"github.com/barkisnet/barkis/types/rest"
	"github.com/barkisnet/barkis/x/asset/internal/types"
)

// HTTP request handler to query the asset params values
func paramsHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParams)
		bz, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var params types.Params
		if err := cliCtx.Codec.UnmarshalJSON(bz, &params); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, params)
	}
}

// HTTP request handler to query a specified token information
func getHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := mux.Vars(r)["symbol"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		resp, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.GetToken, symbol))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var token types.Token
		if err := cliCtx.Codec.UnmarshalJSON(resp, &token); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

// HTTP request handler to list all tokens information
func listHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryTokensParams{
			Page:  page,
			Limit: limit,
		}

		paramsBytes, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.ListToken), paramsBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var tokenList types.TokenList
		if err := cliCtx.Codec.UnmarshalJSON(resp, &tokenList); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, tokenList)
	}
}
