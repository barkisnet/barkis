package config

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/barkisnet/barkis/types"
)

const (
	defaultMinGasPrices = ""
)

type BarkisContext struct {
	*ServerContext
	*viper.Viper
	*AppConfig
}

func NewDefaultContext() *BarkisContext {
	return &BarkisContext{
		NewServerContext(cfg.DefaultConfig(), log.NewTMLogger(log.NewSyncWriter(os.Stdout))),
		viper.New(),
		DefaultAppConfig(),
	}
}

func (context *BarkisContext) ParseAppConfigInPlace() error {
	homeDir := viper.GetString(cli.HomeFlag)
	context.Viper.SetConfigName("app")
	context.Viper.AddConfigPath(homeDir)
	context.Viper.AddConfigPath(filepath.Join(homeDir, "config"))

	if err := context.Viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// ignore not found error, return other errors
		return err
	}

	err := context.Viper.Unmarshal(context.AppConfig)
	if err != nil {
		return err
	}
	return nil
}

type ServerContext struct {
	Config *cfg.Config
	Logger log.Logger
}

func NewServerContext(config *cfg.Config, logger log.Logger) *ServerContext {
	return &ServerContext{config, logger}
}

func NewDefaultServerContext() *ServerContext {
	return &ServerContext{
		cfg.DefaultConfig(),
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	}
}

// Config defines the server's top level configuration
type AppConfig struct {
	BaseConfig    `mapstructure:"base"`
	UpgradeConfig `mapstructure:"upgrade"`
}

// BaseConfig defines the server's basic configuration
type BaseConfig struct {
	// The minimum gas prices a validator is willing to accept for processing a
	// transaction. A transaction's fees must meet the minimum of any denomination
	// specified in this config (e.g. 0.25token1;0.0001token2).
	MinGasPrices string `mapstructure:"minimum-gas-prices"`

	// HaltHeight contains a non-zero height at which a node will gracefully halt
	// and shutdown that can be used to assist upgrades and testing.
	HaltHeight uint64 `mapstructure:"halt-height"`
}

type UpgradeConfig struct {
	RewardUpgrade                int64 `mapstructure:"RewardUpgrade"`
	TokenIssueHeight             int64 `mapstructure:"TokenIssueHeight"`
	UpdateVotingPeriodHeight     int64 `mapstructure:"UpdateVotingPeriodHeight"`
	UpdateTokenSymbolRulesHeight int64 `mapstructure:"UpdateTokenSymbolRulesHeight"`
	DelayedTransferUpgradeHeight int64 `mapstructure:"DelayedTransferUpgradeHeight"`
}

// SetMinGasPrices sets the validator's minimum gas prices.
func (c *AppConfig) SetMinGasPrices(gasPrices sdk.DecCoins) {
	c.MinGasPrices = gasPrices.String()
}

// GetMinGasPrices returns the validator's minimum gas prices based on the set
// configuration.
func (c *AppConfig) GetMinGasPrices() sdk.DecCoins {
	if c.MinGasPrices == "" {
		return sdk.DecCoins{}
	}

	gasPricesStr := strings.Split(c.MinGasPrices, ";")
	gasPrices := make(sdk.DecCoins, len(gasPricesStr))

	for i, s := range gasPricesStr {
		gasPrice, err := sdk.ParseDecCoin(s)
		if err != nil {
			panic(fmt.Errorf("failed to parse minimum gas price coin (%s): %s", s, err))
		}

		gasPrices[i] = gasPrice
	}

	return gasPrices
}

// DefaultAppConfig returns server's default configuration.
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		BaseConfig: BaseConfig{
			MinGasPrices: defaultMinGasPrices,
			HaltHeight:   0,
		},
		UpgradeConfig: UpgradeConfig{
			RewardUpgrade:                math.MaxInt64,
			TokenIssueHeight:             math.MaxInt64,
			UpdateVotingPeriodHeight:     math.MaxInt64,
			UpdateTokenSymbolRulesHeight: math.MaxInt64,
			DelayedTransferUpgradeHeight: math.MaxInt64,
		},
	}
}
