package config

import (
	"bytes"
	"text/template"

	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### main base config options #####

# The minimum gas prices a validator is willing to accept for processing a
# transaction. A transaction's fees must meet the minimum of any denomination
# specified in this config (e.g. 0.25token1;0.0001token2).
minimum-gas-prices = "{{ .BaseConfig.MinGasPrices }}"

# HaltHeight contains a non-zero height at which a node will gracefully halt
# and shutdown that can be used to assist upgrades and testing.
halt-height = {{ .BaseConfig.HaltHeight }}

[upgrade]
# Upgrade to change reward rules
RewardUpgrade = {{ .UpgradeConfig.RewardUpgrade }}

# Upgrade to change reward rules
TokenIssueHeight = {{ .UpgradeConfig.TokenIssueHeight }}

# Upgrade to update voting period
UpdateVotingPeriodHeight = {{ .UpgradeConfig.UpdateVotingPeriodHeight }}

# Upgrade to update token symbol rules
UpdateTokenSymbolRulesHeight = {{ .UpgradeConfig.UpdateTokenSymbolRulesHeight }}

# Upgrade to change token description length limitation
TokenDesLenLimitUpgradeHeight = {{ .UpgradeConfig.TokenDesLenLimitUpgradeHeight }}
`

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("appConfigFileTemplate")
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

// ParseConfig retrieves the default environment configuration for the
// application.
func ParseConfig() (*AppConfig, error) {
	conf := DefaultAppConfig()
	err := viper.Unmarshal(conf)
	return conf, err
}

// WriteConfigFile renders config using the template and writes it to
// configFilePath.
func WriteConfigFile(configFilePath string, config *AppConfig) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	cmn.MustWriteFile(configFilePath, buffer.Bytes(), 0644)
}
