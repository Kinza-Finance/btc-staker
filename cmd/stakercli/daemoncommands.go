package main

import (
	"context"
	"strconv"

	scfg "github.com/babylonchain/btc-staker/stakercfg"
	dc "github.com/babylonchain/btc-staker/stakerservice/client"
	"github.com/urfave/cli"
)

var daemonCommands = []cli.Command{
	{
		Name:      "daemon",
		ShortName: "dn",
		Usage:     "More advanced commands which require staker daemon to be running.",
		Category:  "Daemon commands",
		Subcommands: []cli.Command{
			checkDaemonHealthCmd,
			listOutputsCmd,
			babylonFinalityProvidersCmd,
			getStakeOutputCmd,
			stakeCmd,
			unstakeCmd,
			stakingDetailsCmd,
			listStakingTransactionsCmd,
			withdrawableTransactionsCmd,
			unbondCmd,
		},
	},
}

const (
	stakingDaemonAddressFlag   = "daemon-address"
	offsetFlag                 = "offset"
	limitFlag                  = "limit"
	fpPksFlag                  = "finality-providers-pks"
	stakingTimeBlocksFlag      = "staking-time"
	stakingTransactionHashFlag = "staking-transaction-hash"
	feeRateFlag                = "fee-rate"
	stakerPubKeyFlag           = "staker-pubkey"
)

var (
	defaultStakingDaemonAddress = "tcp://127.0.0.1:" + strconv.Itoa(scfg.DefaultRPCPort)
)

var checkDaemonHealthCmd = cli.Command{
	Name:      "check-health",
	ShortName: "ch",
	Usage:     "Check if staker daemon is running.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "Full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
	},
	Action: checkHealth,
}

var listOutputsCmd = cli.Command{
	Name:      "list-outputs",
	ShortName: "lo",
	Usage:     "List unspent outputs in connected wallet.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "Full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
	},
	Action: listOutputs,
}

var babylonFinalityProvidersCmd = cli.Command{
	Name:      "babylon-finality-providers",
	ShortName: "bfp",
	Usage:     "List current BTC finality providers on Babylon chain",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
		cli.IntFlag{
			Name:  offsetFlag,
			Usage: "offset of the first finality provider to return",
			Value: 0,
		},
		cli.IntFlag{
			Name:  limitFlag,
			Usage: "maximum number of finality providers to return",
			Value: 100,
		},
	},
	Action: babylonFinalityProviders,
}

var getStakeOutputCmd = cli.Command{
	Name:      "getStakeOutput",
	ShortName: "gsto",
	Usage:     "Generate the output address of the staking position",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
		cli.StringFlag{
			Name:     stakerPubKeyFlag,
			Usage:    "BTC public key of the staker",
			Required: true,
		},
		cli.Int64Flag{
			Name:     stakingAmountFlag,
			Usage:    "Staking amount in satoshis",
			Required: true,
		},
		cli.StringSliceFlag{
			Name:     fpPksFlag,
			Usage:    "BTC public keys of the finality providers in hex",
			Required: true,
		},
		cli.Int64Flag{
			Name:     stakingTimeBlocksFlag,
			Usage:    "Staking time in BTC blocks",
			Required: true,
		},
	},
	Action: getStakeOutput,
}

var stakeCmd = cli.Command{
	Name:      "stake",
	ShortName: "st",
	Usage:     "Stake an amount of BTC to Babylon",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
		cli.StringFlag{
			Name:     stakerAddressFlag,
			Usage:    "BTC address of the staker in hex",
			Required: true,
		},
		cli.Int64Flag{
			Name:     stakingAmountFlag,
			Usage:    "Staking amount in satoshis",
			Required: true,
		},
		cli.StringSliceFlag{
			Name:     fpPksFlag,
			Usage:    "BTC public keys of the finality providers in hex",
			Required: true,
		},
		cli.Int64Flag{
			Name:     stakingTimeBlocksFlag,
			Usage:    "Staking time in BTC blocks",
			Required: true,
		},
	},
	Action: stake,
}

var unstakeCmd = cli.Command{
	Name:      "unstake",
	ShortName: "ust",
	Usage:     "Spends staking transaction and sends funds back to staker; this can only be done after timelock of staking transaction expires",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
		cli.StringFlag{
			Name:     stakingTransactionHashFlag,
			Usage:    "Hash of original staking transaction in bitcoin hex format",
			Required: true,
		},
	},
	Action: unstake,
}

var unbondCmd = cli.Command{
	Name:      "unbond",
	ShortName: "ubd",
	Usage:     "initiates unbonding flow: build unbonding tx, send to babylon, wait for signatures, and send unbonding tx to bitcoin",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
		cli.StringFlag{
			Name:     stakingTransactionHashFlag,
			Usage:    "Hash of original staking transaction in bitcoin hex format",
			Required: true,
		},
		cli.IntFlag{
			Name:  feeRateFlag,
			Usage: "fee rate to pay for unbonding tx in sats/kb",
		},
	},
	Action: unbond,
}

var stakingDetailsCmd = cli.Command{
	Name:      "staking-details",
	ShortName: "sds",
	Usage:     "Displays details of staking transaction with given hash",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
		cli.StringFlag{
			Name:     stakingTransactionHashFlag,
			Usage:    "Hash of original staking transaction in bitcoin hex format",
			Required: true,
		},
	},
	Action: stakingDetails,
}

var listStakingTransactionsCmd = cli.Command{
	Name:      "list-staking-transactions",
	ShortName: "lst",
	Usage:     "List current staking transactions in db",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
		cli.IntFlag{
			Name:  offsetFlag,
			Usage: "offset of the first transactions to return",
			Value: 0,
		},
		cli.IntFlag{
			Name:  limitFlag,
			Usage: "maximum number of transactions to return",
			Value: 100,
		},
	},
	Action: listStakingTransactions,
}

var withdrawableTransactionsCmd = cli.Command{
	Name:      "withdrawable-transactions",
	ShortName: "wt",
	Usage:     "List current tranactions that can be withdrawn i.e funds can be transferred back to staker address",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  stakingDaemonAddressFlag,
			Usage: "full address of the staker daemon in format tcp:://<host>:<port>",
			Value: defaultStakingDaemonAddress,
		},
		cli.IntFlag{
			Name:  offsetFlag,
			Usage: "offset of the first transactions to return",
			Value: 0,
		},
		cli.IntFlag{
			Name:  limitFlag,
			Usage: "maximum number of transactions to return",
			Value: 100,
		},
	},
	Action: withdrawableTransactions,
}

func checkHealth(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}

	sctx := context.Background()

	health, err := client.Health(sctx)

	if err != nil {
		return err
	}

	printRespJSON(health)

	return nil
}

func listOutputs(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}

	sctx := context.Background()

	outputs, err := client.ListOutputs(sctx)

	if err != nil {
		return err
	}

	printRespJSON(outputs)

	return nil
}

func babylonFinalityProviders(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}

	sctx := context.Background()

	offset := ctx.Int(offsetFlag)

	if offset < 0 {
		return cli.NewExitError("Offset must be non-negative", 1)
	}

	limit := ctx.Int(limitFlag)

	if limit < 0 {
		return cli.NewExitError("Limit must be non-negative", 1)
	}

	finalityProviders, err := client.BabylonFinalityProviders(sctx, &offset, &limit)

	if err != nil {
		return err
	}

	printRespJSON(finalityProviders)

	return nil
}

func getStakeOutput(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)

	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}
	sctx := context.Background()

	stakerKey := ctx.String(stakerPubKeyFlag)
	stakingAmount := ctx.Int64(stakingAmountFlag)
	fpPks := ctx.StringSlice(fpPksFlag)
	stakingTimeBlocks := ctx.Int64(stakingTimeBlocksFlag)

	results, err := client.GetStakeOutput(sctx, stakerKey, stakingAmount, fpPks, stakingTimeBlocks)
	if err != nil {
		return err
	}
	printRespJSON(results)

	return nil
}

func stake(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}

	sctx := context.Background()

	stakerAddress := ctx.String(stakerAddressFlag)
	stakingAmount := ctx.Int64(stakingAmountFlag)
	fpPks := ctx.StringSlice(fpPksFlag)
	stakingTimeBlocks := ctx.Int64(stakingTimeFlag)

	results, err := client.Stake(sctx, stakerAddress, stakingAmount, fpPks, stakingTimeBlocks)
	if err != nil {
		return err
	}

	printRespJSON(results)

	return nil
}

func unstake(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}

	sctx := context.Background()

	stakingTransactionHash := ctx.String(stakingTransactionHashFlag)

	result, err := client.SpendStakingTransaction(sctx, stakingTransactionHash)
	if err != nil {
		return err
	}

	printRespJSON(result)

	return nil
}

func unbond(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}

	sctx := context.Background()

	stakingTransactionHash := ctx.String(stakingTransactionHashFlag)

	feeRate := ctx.Int(feeRateFlag)

	if feeRate < 0 {
		return cli.NewExitError("Fee rate must be non-negative", 1)
	}

	var fr *int = nil
	if feeRate > 0 {
		fr = &feeRate
	}

	result, err := client.UnbondStaking(sctx, stakingTransactionHash, fr)
	if err != nil {
		return err
	}

	printRespJSON(result)

	return nil
}

func stakingDetails(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}

	sctx := context.Background()

	stakingTransactionHash := ctx.String(stakingTransactionHashFlag)

	result, err := client.StakingDetails(sctx, stakingTransactionHash)
	if err != nil {
		return err
	}

	printRespJSON(result)

	return nil
}

func listStakingTransactions(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}

	sctx := context.Background()

	offset := ctx.Int(offsetFlag)

	if offset < 0 {
		return cli.NewExitError("Offset must be non-negative", 1)
	}

	limit := ctx.Int(limitFlag)

	if limit < 0 {
		return cli.NewExitError("Limit must be non-negative", 1)
	}

	transactions, err := client.ListStakingTransactions(sctx, &offset, &limit)

	if err != nil {
		return err
	}

	printRespJSON(transactions)

	return nil
}

func withdrawableTransactions(ctx *cli.Context) error {
	daemonAddress := ctx.String(stakingDaemonAddressFlag)
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return err
	}

	sctx := context.Background()

	offset := ctx.Int(offsetFlag)

	if offset < 0 {
		return cli.NewExitError("Offset must be non-negative", 1)
	}

	limit := ctx.Int(limitFlag)

	if limit < 0 {
		return cli.NewExitError("Limit must be non-negative", 1)
	}

	transactions, err := client.WithdrawableTransactions(sctx, &offset, &limit)

	if err != nil {
		return err
	}

	printRespJSON(transactions)

	return nil
}
