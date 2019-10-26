package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tm-db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/barkisnet/barkis/codec"
	"github.com/barkisnet/barkis/simapp"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestBarkisdExport(t *testing.T) {
	db := db.NewMemDB()
	gapp := NewBarkisApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	setGenesis(gapp)

	// Making a new app object with the db, so that initchain hasn't been called
	newGapp := NewBarkisApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	_, _, err := newGapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func setGenesis(gapp *BarkisApp) error {

	genesisState := simapp.NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(gapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	gapp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	gapp.Commit()
	return nil
}
