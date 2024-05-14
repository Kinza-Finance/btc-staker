package service

import (
	"context"
	"errors"

	service "github.com/babylonchain/btc-staker/stakerservice"
	dc "github.com/babylonchain/btc-staker/stakerservice/client"
)

func Stake(daemonAddress string, stakerAddress string, stakingAmount int64, fpPks []string, stakingTimeBlocks int64) (*service.ResultStake, error) {
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return nil, err
	}

	sctx := context.Background()

	results, err := client.Stake(sctx, stakerAddress, stakingAmount, fpPks, stakingTimeBlocks)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func Unbond(daemonAddress string, stakingTransactionHash string, feeRate int) (*service.UnbondingResponse, error) {
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return nil, err
	}

	sctx := context.Background()

	if feeRate < 0 {
		return nil, errors.New("fee rate must be non-negative")
	}

	var fr *int = nil
	if feeRate > 0 {
		fr = &feeRate
	}

	result, err := client.UnbondStaking(sctx, stakingTransactionHash, fr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Unstake(daemonAddress string, stakingTransactionHash string) (*service.SpendTxDetails, error) {
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return nil, err
	}

	sctx := context.Background()

	result, err := client.SpendStakingTransaction(sctx, stakingTransactionHash)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetStakeOutput(daemonAddress string, stakerKey string, stakingAmount int64, fpPks []string, stakingTimeBlocks int64) (*service.ResultStakeOutput, error) {
	client, err := dc.NewStakerServiceJsonRpcClient(daemonAddress)
	if err != nil {
		return nil, err
	}
	sctx := context.Background()

	results, err := client.GetStakeOutput(sctx, stakerKey, stakingAmount, fpPks, stakingTimeBlocks)
	if err != nil {
		return nil, err
	}

	return results, nil
}
