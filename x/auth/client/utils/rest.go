package utils

import (
	"log"
	"net/http"

	"github.com/barkisnet/barkis/client/context"
	"github.com/barkisnet/barkis/client/flags"
	"github.com/barkisnet/barkis/crypto/keys/keyerror"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/types/rest"
	"github.com/barkisnet/barkis/x/auth/types"
)

// WriteGenerateStdTxResponse writes response for the generate only mode.
func WriteGenerateStdTxResponse(w http.ResponseWriter, cliCtx context.CLIContext, br rest.BaseReq, msgs []sdk.Msg) {
	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gas, err := flags.ParseGas(br.Gas)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := types.NewTxBuilder(
		GetTxEncoder(cliCtx.Codec), br.AccountNumber, br.Sequence, gas, gasAdj,
		br.Simulate, br.ChainID, br.Memo, br.Fees, br.GasPrices,
	)

	if br.Simulate || simAndExec {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errInvalidGasAdjustment.Error())
			return
		}

		txBldr, err = EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if br.Simulate {
			rest.WriteSimulationResponse(w, cliCtx.Codec, txBldr.Gas())
			return
		}
	}
	if br.GenerateOnly {
		stdMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		output, err := cliCtx.Codec.MarshalJSON(types.NewStdTx(stdMsg.Msgs, stdMsg.Fee, nil, stdMsg.Memo))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(output); err != nil {
			log.Printf("could not write response: %v", err)
		}
		return
	}

	txBytes, err := txBldr.BuildAndSign(cliCtx.GetFromName(), br.Password, msgs)
	if keyerror.IsErrKeyNotFound(err) {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	} else if keyerror.IsErrWrongPassword(err) {
		rest.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	} else if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	rest.PostProcessResponse(w, cliCtx, res)
}
