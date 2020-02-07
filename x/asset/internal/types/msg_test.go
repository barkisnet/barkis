package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/barkisnet/barkis/types"
	"github.com/tendermint/tendermint/crypto"
)

func TestIssueMsgValidation(t *testing.T) {
	var emptyAddr sdk.AccAddress
	issuer := sdk.AccAddress(crypto.AddressHash([]byte("issuer")))

	cases := []struct {
		valid   bool
		errCode CodeType
		tx      IssueMsg
	}{
		{true, 0, NewIssueMsg(issuer, "bitcoin", "btc", 21000000000000, false, 6, "bitcoin on barkisnet")},
		{false, sdk.CodeInvalidAddress, NewIssueMsg(emptyAddr, "bitcoin", "btc", 21000000000000, false, 6, "bitcoin on barkisnet")},

		{false, CodeInvalidTokenSymbol, NewIssueMsg(issuer, "bitcoin", "ubarkis", 21000000000000, false, 6, "bitcoin on barkisnet")},
		{false, CodeInvalidTokenSymbol, NewIssueMsg(issuer, "bitcoin", "barkis", 21000000000000, false, 6, "bitcoin on barkisnet")},
		{false, CodeInvalidTokenSymbol, NewIssueMsg(issuer, "bitcoin", "bt1", 21000000000000, false, 6, "bitcoin on barkisnet")},
		{false, CodeInvalidTokenSymbol, NewIssueMsg(issuer, "bitcoin", "btc_", 21000000000000, false, 6, "bitcoin on barkisnet")},
		{false, CodeInvalidTokenSymbol, NewIssueMsg(issuer, "bitcoin", "btcbtcbtcbtc", 21000000000000, false, 6, "bitcoin on barkisnet")},

		{false, CodeInvalidTokenName, NewIssueMsg(issuer, "ubarkis", "btc", 21000000000000, false, 6, "bitcoin on barkisnet")},
		{false, CodeInvalidTokenName, NewIssueMsg(issuer, "barkis", "btc", 21000000000000, false, 6, "bitcoin on barkisnet")},
		{false, CodeInvalidTokenName, NewIssueMsg(issuer, "bitcoinbitcoinbitcoinbitcoinbitcoin", "btc", 21000000000000, false, 6, "bitcoin on barkisnet")},

		{false, CodeInvalidTotalSupply, NewIssueMsg(issuer, "bitcoin", "btc", 9000000000000000001, false, 6, "bitcoin on barkisnet")},
		{false, CodeInvalidDecimal, NewIssueMsg(issuer, "bitcoin", "btc", 21000000000000, false, -1, "bitcoin on barkisnet")},
		{false, CodeInvalidTokenDescription, NewIssueMsg(issuer, "bitcoin", "btc", 21000000000000, false, 6, "bitcoin on barkisnetbitcoin on barkisnetbitcoin on barkisnetbitcoin on barkisnetbitcoin on barkisnetbitcoin on barkisnetbitcoin on barkisnetbitcoin on barkisnet")},
	}

	for index, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
			require.Equal(t, tc.errCode, err.Code(), fmt.Sprintf("index: %d, errMsg: %s", index, err.Error()))
		}
	}
}

func TestMintMsgValidation(t *testing.T) {
	var emptyAddr sdk.AccAddress
	minter := sdk.AccAddress(crypto.AddressHash([]byte("minter")))

	cases := []struct {
		valid   bool
		errCode CodeType
		tx      MintMsg
	}{
		{true, 0, NewMintMsg(minter, "btc_123", 10000)},

		{false, sdk.CodeInvalidAddress, NewMintMsg(emptyAddr, "btc_123", 10000)},

		{false, CodeInvalidTokenSymbol, NewMintMsg(minter, "btc", 10000)},
		{false, CodeInvalidTokenSymbol, NewMintMsg(minter, "btc_", 10000)},
		{false, CodeInvalidTokenSymbol, NewMintMsg(minter, "ubarkis", 10000)},
		{false, CodeInvalidTokenSymbol, NewMintMsg(minter, "barkis", 10000)},

		{false, CodeInvalidMintAmount, NewMintMsg(minter, "btc_123", 9000000000000000001)},
	}

	for index, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
			require.Equal(t, tc.errCode, err.Code(), fmt.Sprintf("index: %d, errMsg: %s", index, err.Error()))
		}
	}
}
