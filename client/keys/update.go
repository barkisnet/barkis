package keys

import (
	"bufio"
	"encoding/json"
	"github.com/barkisnet/barkis/crypto/keys"
	"github.com/barkisnet/barkis/crypto/keys/keyerror"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/barkisnet/barkis/client/input"
)

func updateKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Change the password used to protect private key",
		RunE:  runUpdateCmd,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func runUpdateCmd(cmd *cobra.Command, args []string) error {
	name := args[0]

	buf := bufio.NewReader(cmd.InOrStdin())
	kb, err := NewKeyBaseFromHomeFlag()
	if err != nil {
		return err
	}
	oldpass, err := input.GetPassword("Enter the current passphrase:", buf)
	if err != nil {
		return err
	}

	getNewpass := func() (string, error) {
		return input.GetCheckPassword(
			"Enter the new passphrase:",
			"Repeat the new passphrase:", buf)
	}
	if err := kb.Update(name, oldpass, getNewpass); err != nil {
		return err
	}

	cmd.PrintErrln("Password successfully updated!")
	return nil
}

///////////////////////
// REST

// update key request REST body
type UpdateKeyBody struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

// update key REST handler
func UpdateKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var kb keys.Keybase
	var m UpdateKeyBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	kb, err = NewKeyBaseFromHomeFlag()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	getNewpass := func() (string, error) { return m.NewPassword, nil }

	err = kb.Update(name, m.OldPassword, getNewpass)
	if keyerror.IsErrKeyNotFound(err) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
		return
	} else if keyerror.IsErrWrongPassword(err) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(err.Error()))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
