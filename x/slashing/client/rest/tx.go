package rest

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/barkisnet/barkis/client/context"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/types/rest"
	"github.com/barkisnet/barkis/x/auth/client/utils"
	"github.com/barkisnet/barkis/x/slashing/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/slashing/validators/{validatorAddr}/unjail",
		unjailRequestHandlerFn(cliCtx),
	).Methods("POST")
}

// Unjail TX body
type UnjailReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

func unjailRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		bech32validator := vars["validatorAddr"]

		var req UnjailReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		valAddr, err := sdk.ValAddressFromBech32(bech32validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, valAddr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own validator address")
			return
		}

		msg := types.NewMsgUnjail(valAddr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// derive the from account address and name from the Keybase
		fromAddress, fromName, err := context.GetFromFieldsFromAddr(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithFromName(fromName).WithFromAddress(fromAddress).WithBroadcastMode(req.BaseReq.BroadcastMode)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
