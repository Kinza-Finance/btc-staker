package main

import (
	"fmt"

	babylonApp "github.com/babylonchain/babylon/app"
	"github.com/babylonchain/btc-staker/stakercfg"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/go-bip39"
	"github.com/jessevdk/go-flags"
	"github.com/urfave/cli"
)

var adminCommands = []cli.Command{
	{
		Name:      "admin",
		ShortName: "ad",
		Usage:     "Differnt utility and admin commands",
		Category:  "Admin",
		Subcommands: []cli.Command{
			dumpCfgCommand,
			createCosmosKeyringCommand,
		},
	},
}

const (
	configFileDirFlag = "config-file-dir"
)

var (
	defaultConfigPath = stakercfg.DefaultConfigFile
)

var dumpCfgCommand = cli.Command{
	Name:      "dump-config",
	ShortName: "dc",
	Usage:     "dump default configuration file",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  configFileDirFlag,
			Usage: "path to the directory where the config file will be dumped",
			Value: defaultConfigPath,
		},
	},
	Action: dumpCfg,
}

func dumpCfg(c *cli.Context) error {
	configPath := c.String(configFileDirFlag)

	if stakercfg.FileExists(configPath) {
		return cli.NewExitError(
			fmt.Sprintf("config already exists under provided path: %s", configPath),
			1,
		)
	}

	defaultConfig := stakercfg.DefaultConfig()
	fileParser := flags.NewParser(&defaultConfig, flags.Default)

	err := flags.NewIniParser(fileParser).WriteFile(configPath, flags.IniIncludeComments|flags.IniIncludeDefaults)

	if err != nil {
		return err
	}

	return nil
}

const (
	mnemonicEntropySize = 256
	secp256k1Type       = "secp256k1"

	chainIdFlag        = "chain-id"
	keyringBackendFlag = "keyring-backend"
	keyNameFlag        = "key-name"
	keyringDir         = "keyring-dir"
)

var (
	defaultBBNconfig = stakercfg.DefaultBBNConfig()
	defaultChainID   = defaultBBNconfig.ChainID
	defaultBackend   = defaultBBNconfig.KeyringBackend
	defaultKeyName   = defaultBBNconfig.Key
	defaultKeyDir    = defaultBBNconfig.KeyDirectory
)

func createKey(name string, kr keyring.Keyring) (*keyring.Record, error) {
	keyringAlgos, _ := kr.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(secp256k1Type, keyringAlgos)
	if err != nil {
		return nil, err
	}

	// read entropy seed straight from tmcrypto.Rand and convert to mnemonic
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return nil, err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return nil, err
	}

	record, err := kr.NewAccount(name, mnemonic, "", "", algo)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createKeyRing(c *cli.Context) error {
	keyringOptions := []keyring.Option{}
	keyringOptions = append(keyringOptions, func(options *keyring.Options) {
		options.SupportedAlgos = keyring.SigningAlgoList{hd.Secp256k1}
		options.SupportedAlgosLedger = keyring.SigningAlgoList{hd.Secp256k1}
	})

	encConf := babylonApp.GetEncodingConfig()
	encConf.InterfaceRegistry.RegisterImplementations(
		(*cryptotypes.PubKey)(nil),
		&secp256k1.PubKey{},
	)

	chainId := c.String(chainIdFlag)
	backend := c.String(keyringBackendFlag)
	keyName := c.String(keyNameFlag)
	keyDir := c.String(keyringDir)

	kb, err := keyring.New(
		chainId,
		backend,
		keyDir,
		nil,
		encConf.Marshaler,
		keyringOptions...)

	if err != nil {
		return err
	}

	_, err = createKey(keyName, kb)

	if err != nil {
		return err
	}

	list, err := kb.List()

	if err != nil {
		return err
	}

	fmt.Println("Accounts in keyring:")
	for _, r := range list {
		fmt.Println(r.Name)
	}

	return nil
}

var createCosmosKeyringCommand = cli.Command{
	Name:      "create-keyring",
	ShortName: "ck",
	Usage: "Creates cosmos keyring with secp256k1 key with and account with provided name." +
		"If account already existis in the keyring, new address will be created for given key.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  keyNameFlag,
			Usage: "Name of the key account to be created",
			Value: defaultKeyName,
		},
		cli.StringFlag{
			Name:  keyringBackendFlag,
			Usage: "Backend for keyring",
			Value: defaultBackend,
		},
		cli.StringFlag{
			Name:  chainIdFlag,
			Usage: "Chain ID for which account are created",
			Value: defaultChainID,
		},
		cli.StringFlag{
			Name:  keyringDir,
			Usage: "Directory in which keyring should be created",
			Value: defaultKeyDir,
		},
	},
	Action: createKeyRing,
}
