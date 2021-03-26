# What is Barkis?

`barkis` is the name of the barkis application for the BarkisNet. It comes with 2 main entrypoints:

- `barkisd`: The Barkis Daemon, runs a full-node of the `barkis` application.
- `barkiscli`: The Barkis command-line interface, which enables interaction with a Barkis full-node.

`barkis` is built on the barkis using the following modules:

- `x/auth`: Accounts and signatures.
- `x/bank`: Token transfers.
- `x/staking`: Staking logic.
- `x/mint`: Inflation logic.
- `x/distribution`: Fee distribution logic.
- `x/slashing`: Slashing logic.
- `x/gov`: Governance logic.
- `x/ibc`: Inter-blockchain transfers.
- `x/params`: Handles app-level parameters.

Next, learn how to [install Barkis](./installation.md).
