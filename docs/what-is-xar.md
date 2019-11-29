# What is Xar?

`xar` is the name of the Fantom application for the Cosmos Hub. It comes with 2 main entrypoints:

- `xard`: The Xar Daemon, runs a full-node of the `xar` application.
- `xarcli`: The Xar command-line interface, which enables interaction with a Xar full-node.

`xar` is built on the Cosmos & Fantom SDK using the following modules:

- `x/auth`: Accounts and signatures.
- `x/bank`: Token transfers.
- `x/staking`: Staking logic.
- `x/mint`: Inflation logic.
- `x/distribution`: Fee distribution logic.
- `x/slashing`: Slashing logic.
- `x/gov`: Governance logic.
- `x/ibc`: Inter-blockchain transfers.
- `x/params`: Handles app-level parameters.

Next, learn how to [install Xar](./installation.md).
