package client

import (
	govclient "github.com/barkisnet/barkis/x/gov/client"
	"github.com/barkisnet/barkis/x/params/client/cli"
	"github.com/barkisnet/barkis/x/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
