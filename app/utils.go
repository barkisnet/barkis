//nolint
package app

import (
	"io"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/barkisnet/barkis/baseapp"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/staking"
)

var (
	genesisFile        string
	paramsFile         string
	exportParamsPath   string
	exportParamsHeight int
	exportStatePath    string
	exportStatsPath    string
	seed               int64
	initialBlockHeight int
	numBlocks          int
	blockSize          int
	enabled            bool
	verbose            bool
	lean               bool
	commit             bool
	period             int
	onOperation        bool // TODO Remove in favor of binary search for invariant violation
	allInvariants      bool
	genesisTime        int64
)

// DONTCOVER

// NewBarkisAppUNSAFE is used for debugging purposes only.
//
// NOTE: to not use this function with non-test code
func NewBarkisAppUNSAFE(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp),
) (bapp *BarkisApp, keyMain, keyStaking *sdk.KVStoreKey, stakingKeeper staking.Keeper) {

	bapp = NewBarkisApp(logger, db, traceStore, loadLatest, invCheckPeriod, baseAppOptions...)
	return bapp, bapp.keys[baseapp.MainStoreKey], bapp.keys[staking.StoreKey], bapp.stakingKeeper
}
