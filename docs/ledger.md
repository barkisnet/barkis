# Ledger Nano Support

Using a hardware wallet to store your keys greatly improves the security of your crypto assets. The Ledger device acts as an enclave of the seed and private keys, and the process of signing transaction takes place within it. No private information ever leaves the Ledger device. The following is a short tutorial on using the Barkis Ledger app with the Barkis CLI or the [Lunie.io](https://lunie.io/#/) web wallet.

At the core of a Ledger device there is a mnemonic seed phrase that is used to generate private keys. This phrase is generated when you initialize you Ledger. The mnemonic is compatible with Barkis and can be used to seed new accounts.

::: danger
Do not lose or share your 24 words with anyone. To prevent theft or loss of funds, it is best to keep multiple copies of your mnemonic stored in safe, secure places. If someone is able to gain access to your mnemonic, they will fully control the accounts associated with them.
:::

## Barkis CLI + Ledger Nano

The tool used to generate addresses and transactions on the BarkisNet is `barkiscli`. Here is how to get started. If using a CLI tool is unfamiliar to you, scroll down and follow instructions for using the Lunie.io web wallet instead.

### Before you Begin

- [Install the Barkis app onto your Ledger](https://github.com/barkis/ledger-barkis/blob/master/README.md#installing)
- [Install Golang](https://golang.org/doc/install)
- [Install Barkis](https://barkis.network/docs/barkisnet/installation.html)

Verify that barkiscli is installed correctly with the following command

```bash
barkiscli version --long

➜ barkiscli: 0.34.3
git commit: 67ab0b1e1d1e5b898c8cbdede35ad5196dba01b2
vendor hash: 0341b356ad7168074391ca7507f40b050e667722
build tags: netgo ledger
go version go1.11.5 darwin/amd64

```

### Add your Ledger key

- Connect and unlock your Ledger device.
- Open the Barkis app on your Ledger.
- Create an account in barkiscli from your ledger key.

::: tip
Be sure to change the _keyName_ parameter to be a meaningful name. The `ledger` flag tells `barkiscli` to use your Ledger to seed the account.
:::

```bash
barkiscli keys add <keyName> --ledger

➜ NAME: TYPE: ADDRESS:     PUBKEY:
<keyName> ledger barkis1... barkispub1...
```

Barkis uses [HD Wallets](./hd-wallets.md). This means you can setup many accounts using the same Ledger seed. To create another account from your Ledger device, run;

```bash
barkiscli keys add <secondKeyName> --ledger
```

### Confirm your address

Run this command to display your address on the device. Use the `keyName` you gave your ledger key. The `-d` flag is supported in version `1.5.0` and higher.

```bash
barkiscli keys show <keyName> -d
```

Confirm that the address displayed on the device matches that displayed when you added the key.

### Connect to a full node

Next, you need to configure barkiscli with the URL of a Barkis full node and the appropriate `chain_id`. In this example we connect to the public load balanced full node operated by Chorus One on the `barkisnet-2` chain. But you can point your `barkiscli` to any Barkis full node. Be sure that the `chain_id` is set to the same chain as the full node.

```bash
barkiscli config node https://barkis.chorus.one:26657
barkiscli config chain_id barkisnet-2
```

Test your connection with a query such as:

``` bash
`barkiscli query staking validators`
```

::: tip
To run your own full node locally [read more here.](https://barkis.network/docs/barkisnet/join-mainnet.html#setting-up-a-new-node).
:::

### Sign a transaction

You are now ready to start signing and sending transactions. Send a transaction with barkiscli using the `tx send` command.

``` bash
barkiscli tx send --help # to see all available options.
```

::: tip
Be sure to unlock your device with the PIN and open the Barkis app before trying to run these commands
:::

Use the `keyName` you set for your Ledger key and barkis will connect with the Barkis Ledger app to then sign your transaction.

```bash
barkiscli tx send <keyName> <destinationAddress> <amount><denomination>
```

When prompted with `confirm transaction before signing`, Answer `Y`.

Next you will be prompted to review and approve the transaction on your Ledger device. Be sure to inspect the transaction JSON displayed on the screen. You can scroll through each field and each message. Scroll down to read more about the data fields of a standard transaction object.

Now, you are all set to start [sending transactions on the network](./delegator-guide-cli.md#sending-transactions).

### Receive funds

To receive funds to the Barkis account on your Ledger device, retrieve the address for your Ledger account (the ones with `TYPE ledger`) with this command:

```bash
barkiscli keys list

➜ NAME: TYPE: ADDRESS:     PUBKEY:
<keyName> ledger barkis1... barkispub1...
```

### Further documentation

Not sure what `barkiscli` can do? Simply run the command without arguments to output documentation for the commands in supports.

::: tip
The `barkiscli` help commands are nested. So `$ barkiscli` will output docs for the top level commands (status, config, query, and tx). You can access documentation for sub commands with further help commands.

For example, to print the `query` commands:

```bash
barkiscli query --help
```

Or to print the `tx` (transaction) commands:

```bash
barkiscli tx --help
```
:::

# The Barkis Standard Transaction

Transactions in Barkis embed the [Standard Transaction type](https://godoc.org/github.com/barkis/barkis/x/auth#StdTx) from the Barkis. The Ledger device displays a serialized JSON representation of this object for you to review before signing the transaction. Here are the fields and what they mean:

- `chain-id`: The chain to which you are broadcasting the tx, such as the `barkis-13003` testnet or `barkisnet-2`: mainnet.
- `account_number`: The global id of the sending account assigned when the account receives funds for the first time.
- `sequence`: The nonce for this account, incremented with each transaction.
- `fee`: JSON object describing the transaction fee, its gas amount and coin denomination
- `memo`: optional text field used in various ways to tag transactions.
- `msgs_<index>/<field>`: The array of messages included in the transaction. Double click to drill down into nested fields of the JSON.
