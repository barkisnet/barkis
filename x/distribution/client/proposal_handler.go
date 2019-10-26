package client

import (
	"github.com/barkisnet/barkis/x/distribution/client/cli"
	"github.com/barkisnet/barkis/x/distribution/client/rest"
	govclient "github.com/barkisnet/barkis/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
