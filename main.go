package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dgraph-io/badger"
	"github.com/iamolegga/enviper"
	"github.com/spf13/viper"
	abciclient "github.com/tendermint/tendermint/abci/client"
	abci "github.com/tendermint/tendermint/abci/types"
	tmcfg "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/service"
	tmnode "github.com/tendermint/tendermint/node"
)

var homeDir = "run"

func main() {
	db, err := badger.Open(badger.DefaultOptions(path.Join(homeDir, "data/badger")))
	if err != nil {
		panic(fmt.Sprintf("failed to open badger db: %v", err))
	}
	defer db.Close()
	app := NewKVStoreApplication(db)

	node, err := newTendermintNode(app)
	if err != nil {
		panic(fmt.Sprintf("failed to created tendermint app: %v", err))
	}

	node.Start()
	defer func() {
		node.Stop()
		node.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func newTendermintNode(app abci.Application) (service.Service, error) {
	tmConfig := tmcfg.DefaultConfig()

	// bind all the config keys before loading env variables
	// https://github.com/spf13/viper/issues/761#issuecomment-646644647
	ev := enviper.New(viper.GetViper())

	tmConfigFile := path.Join(homeDir, "config", "config.toml")
	ev.SetConfigFile(tmConfigFile)
	if err := ev.Unmarshal(tmConfig); err != nil {
		return nil, err
	}

	if err := ev.MergeInConfig(); err != nil {
		return nil, err
	}

	tmConfig.SetRoot(homeDir)
	tmcfg.EnsureRoot(tmConfig.RootDir)

	logger := tmlog.MustNewDefaultLogger(tmConfig.LogFormat, tmConfig.LogLevel, false)

	return tmnode.New(
		tmConfig,
		logger,
		abciclient.NewLocalCreator(app),
		nil,
	)
}
