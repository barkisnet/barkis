package keeper

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/barkisnet/barkis/client"
	"github.com/barkisnet/barkis/codec"
	sdk "github.com/barkisnet/barkis/types"
	assetTypes "github.com/barkisnet/barkis/x/asset/internal/types"
	"github.com/barkisnet/barkis/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		if !sdk.GlobalUpgradeMgr.IsUpgradeApplied(sdk.TokenIssueUpgrade) {
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("asset related query is not support until %d",
				sdk.GlobalUpgradeMgr.GetUpgradeHeight(sdk.TokenIssueUpgrade)))
		}
		switch path[0] {
		case assetTypes.QueryParams:
			return queryParams(ctx, path[1:], req, k)
		case assetTypes.GetToken:
			return queryToken(ctx, path[1:], req, k)
		case assetTypes.ListToken:
			return listToken(ctx, path[1:], req, k)
		case assetTypes.GetDelayedTranfer:
			return getDelayedTransfer(ctx, path[1:], req, k)
		case assetTypes.ListDelayedTranfer:
			return  listDelayedTransfer(ctx, path[1:], req, k)
		case assetTypes.ListDelayedTranferFrom:
			return  listDelayedTransferFrom(ctx, path[1:], req, k)
		case assetTypes.ListDelayedTranferTo:
			return  listDelayedTransferTo(ctx, path[1:], req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)
	if params == nil {
		return nil, sdk.ErrInternal("failed to get asset parameters")
	}
	bz, err := codec.MarshalJSONIndent(k.cdc, *params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryToken(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	if len(path) < 1 {
		return nil, sdk.ErrUnknownRequest("wrong query request")
	}
	tokenSymbol := path[0]
	token := k.GetToken(ctx, tokenSymbol)
	if token == nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("token %s is not exist", tokenSymbol))
	}
	bz, err := codec.MarshalJSONIndent(k.cdc, *token)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func listToken(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params assetTypes.QueryTokensParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	iter := k.ListToken(ctx)
	defer iter.Close()
	var tokens []*assetTypes.Token
	for ; iter.Valid(); iter.Next() {
		token := k.DecodeToken(iter.Value())
		tokens = append(tokens, token)
	}

	var queryResult []*assetTypes.Token
	start, end := client.Paginate(len(tokens), params.Page, params.Limit, assetTypes.DefaultQueryLimit)
	if start < 0 || end < 0 {
		queryResult = []*assetTypes.Token{}
	} else {
		queryResult = tokens[start:end]
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, queryResult)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func getDelayedTransfer(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	if len(path) < 1 {
		return nil, sdk.ErrUnknownRequest("wrong query request")
	}
	sequence, err := strconv.Atoi(path[0])
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("invalid parameter: %s", err.Error()))
	}

	delayedTransfer := k.GetDelayedTransfer(ctx, int64(sequence))
	if delayedTransfer == nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("delayedTransfer with sequence %d is not exist", sequence))
	}
	bz, err := codec.MarshalJSONIndent(k.cdc, *delayedTransfer)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func listDelayedTransfer(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params assetTypes.QueryDelayedTranferParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	iter := k.ListDelayedTransfer(ctx)
	defer iter.Close()
	var delayedTransferList []*assetTypes.DelayedTransfer
	for ; iter.Valid(); iter.Next() {
		delayedTransfer := k.DecodeDelayedTransfer(iter.Value())
		delayedTransferList = append(delayedTransferList, delayedTransfer)
	}

	var queryResult []*assetTypes.DelayedTransfer
	start, end := client.Paginate(len(delayedTransferList), params.Page, params.Limit, assetTypes.DefaultQueryLimit)
	if start < 0 || end < 0 {
		queryResult = []*assetTypes.DelayedTransfer{}
	} else {
		queryResult = delayedTransferList[start:end]
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, queryResult)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func listDelayedTransferFrom(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params assetTypes.QueryDelayedTranferParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	if len(path) < 1 {
		return nil, sdk.ErrUnknownRequest("wrong query request")
	}

	from, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdk.ErrInvalidAddress(err.Error())
	}

	iter := k.ListDelayedTransferFrom(ctx, from)
	defer iter.Close()
	var delayedTransferList []*assetTypes.DelayedTransfer
	for ; iter.Valid(); iter.Next() {
		sequenceBytes := iter.Value()
		sequence := int64(binary.BigEndian.Uint64(sequenceBytes))
		delayedTransfer := k.GetDelayedTransfer(ctx, sequence)
		delayedTransferList = append(delayedTransferList, delayedTransfer)
	}

	var queryResult []*assetTypes.DelayedTransfer
	start, end := client.Paginate(len(delayedTransferList), params.Page, params.Limit, assetTypes.DefaultQueryLimit)
	if start < 0 || end < 0 {
		queryResult = []*assetTypes.DelayedTransfer{}
	} else {
		queryResult = delayedTransferList[start:end]
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, queryResult)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func listDelayedTransferTo(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params assetTypes.QueryDelayedTranferParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	to, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdk.ErrInvalidAddress(err.Error())
	}

	iter := k.ListDelayedTransferTo(ctx, to)
	defer iter.Close()
	var delayedTransferList []*assetTypes.DelayedTransfer
	for ; iter.Valid(); iter.Next() {
		sequenceBytes := iter.Value()
		sequence := int64(binary.BigEndian.Uint64(sequenceBytes))
		delayedTransfer := k.GetDelayedTransfer(ctx, sequence)
		delayedTransferList = append(delayedTransferList, delayedTransfer)
	}

	var queryResult []*assetTypes.DelayedTransfer
	start, end := client.Paginate(len(delayedTransferList), params.Page, params.Limit, assetTypes.DefaultQueryLimit)
	if start < 0 || end < 0 {
		queryResult = []*assetTypes.DelayedTransfer{}
	} else {
		queryResult = delayedTransferList[start:end]
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, queryResult)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}
