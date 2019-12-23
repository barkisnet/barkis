package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpgrade(t *testing.T) {
	type testCase struct {
		upgradeName string
		msgName     string
		storeName   string
		config      UpgradeConfig
		blockHeight int64

		upgradeResult bool
		msgCheck      bool
		storeCheck    bool
	}

	testCases := []testCase{
		{
			upgradeName: "tokenIssue",
			msgName:     "issueToken",
			storeName:   "token",
			config: UpgradeConfig{
				UpgradeHeight:  map[string]int64{"tokenIssue": 10000, "bugfix": 20000},
				NewStoreHeight: map[string]int64{"token": 10000},
				NewMsgHeight:   map[string]int64{"issueToken": 10000, "mintToken": 10000},
			},
			blockHeight: 10000,

			upgradeResult: true,
			msgCheck:      true,
			storeCheck:    true,
		},
		{
			upgradeName: "tokenIssue",
			msgName:     "mintToken",
			storeName:   "token",
			config: UpgradeConfig{
				UpgradeHeight:  map[string]int64{"tokenIssue": 10000, "bugfix": 20000},
				NewStoreHeight: map[string]int64{"token": 10000},
				NewMsgHeight:   map[string]int64{"issueToken": 10000, "mintToken": 10000},
			},
			blockHeight:   9999,
			upgradeResult: false,
			msgCheck:      false,
			storeCheck:    false,
		},
		{
			upgradeName: "tokenIssue",
			msgName:     "mintToken",
			storeName:   "token",
			config: UpgradeConfig{
				UpgradeHeight:  map[string]int64{"tokenIssue": 10000, "bugfix": 20000},
				NewStoreHeight: map[string]int64{"token": 10000},
				NewMsgHeight:   map[string]int64{"issueToken": 10000, "mintToken": 10000},
			},
			blockHeight:   10001,
			upgradeResult: true,
			msgCheck:      true,
			storeCheck:    true,
		},
		{
			upgradeName: "tokenIssue",
			msgName:     "mintToken",
			storeName:   "token1",
			config: UpgradeConfig{
				UpgradeHeight:  map[string]int64{"tokenIssue": 10000, "bugfix": 20000},
				NewStoreHeight: map[string]int64{"token": 10000},
				NewMsgHeight:   map[string]int64{"issueToken": 10000, "mintToken": 10000},
			},
			blockHeight: 10000,

			upgradeResult: true,
			msgCheck:      true,
			storeCheck:    true,
		},
		{
			upgradeName: "bugfix",
			msgName:     "mintToken",
			storeName:   "token",
			config: UpgradeConfig{
				UpgradeHeight:  map[string]int64{"tokenIssue": 10000, "bugfix": 20000},
				NewStoreHeight: map[string]int64{"token": 10000},
				NewMsgHeight:   map[string]int64{"issueToken": 10000, "mintToken": 10000},
			},
			blockHeight:   20000,
			upgradeResult: true,
			msgCheck:      true,
			storeCheck:    true,
		},
		{
			upgradeName: "bugfix",
			msgName:     "mintToken",
			storeName:   "token",
			config: UpgradeConfig{
				UpgradeHeight:  map[string]int64{"tokenIssue": 10000, "bugfix": 20000},
				NewStoreHeight: map[string]int64{"token": 10000},
				NewMsgHeight:   map[string]int64{"issueToken": 10000, "mintToken": 10000},
			},
			blockHeight:   19999,
			upgradeResult: false,
			msgCheck:      true,
			storeCheck:    true,
		},
		{
			upgradeName: "bugfix",
			msgName:     "mintToken",
			storeName:   "token",
			config: UpgradeConfig{
				UpgradeHeight:  map[string]int64{"tokenIssue": 10000, "bugfix": 20000},
				NewStoreHeight: map[string]int64{"token": 10000},
				NewMsgHeight:   map[string]int64{"issueToken": 10000, "mintToken": 10000},
			},
			blockHeight:   20001,
			upgradeResult: true,
			msgCheck:      true,
			storeCheck:    true,
		},
	}

	for index, tc := range testCases {
		GlobalUpgradeMgr.Config = tc.config
		GlobalUpgradeMgr.SetBlockHeight(tc.blockHeight)
		require.Equal(t, tc.upgradeResult, GlobalUpgradeMgr.IsUpgradeApplied(tc.upgradeName), fmt.Sprintf("upgrade height test case failed, index: %d", index))
		require.Equal(t, tc.msgCheck, GlobalUpgradeMgr.MsgCheck(tc.msgName), fmt.Sprintf("new msg test case failed, index: %d", index))
		require.Equal(t, tc.storeCheck, GlobalUpgradeMgr.StoreCheck(tc.storeName), fmt.Sprintf("new store test case failed, index: %d", index))
	}
}
