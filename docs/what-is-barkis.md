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

>About the BarkisNet: The BarkisNet is the first Hub to be launched in the BarkisNet. The role of a Hub is to facilitate transfers between blockchains. If a blockchain connects to a Hub via IBC, it automatically gains access to all the other blockchains that are connected to it. The BarkisNet is a public Proof-of-Stake chain. Its staking token is called the Atom.

Next, learn how to [install Barkis](./installation.md).
