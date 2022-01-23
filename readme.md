# Tendermint Built-in KVStore app

## Install & build

Pull all vendor dependencies:

```bash
go mod vendor
```

Build app:

```bash
go build -o kvapp
```

## Run

Configure node first, all config and data files will go into the `run/` folder.

```bash
mkdir -p run && cp -r sample/* run/
```

Then start the node:

```bash
./kvapp
```

You can now send transactions through the Tendermint API:

```bash
curl -s 'localhost:26657/broadcast_tx_commit?tx="tendermint=rocks"'
```

This one should fail, as the validate function expect transaction content to be of the form "key=value".

```bash
curl -s 'localhost:26657/broadcast_tx_commit?tx="tendermint"'
```

Replaying an old transaction should also fail:

```bash
curl -s 'localhost:26657/broadcast_tx_commit?tx="tendermint=rocks"'
```
