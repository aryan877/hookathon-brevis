package dynamicfee

import (
	"math/big"
	"testing"

	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/brevis-network/brevis-sdk/test"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestAppCircuit(t *testing.T) {
	app, err := sdk.NewBrevisApp()
	if err != nil {
		t.Fatalf("Failed to create BrevisApp: %v", err)
	}

	// Constants
	hookAddress := common.HexToAddress("0xCb38F69700054D326Ecc89ef2486255b528ffCAa5f")
	poolDataUpdatedEventID := common.HexToHash("0xd0f41fd5b4d393ea3222f2ecd77d99386e8ad292339ad0bbc6e3e5530e5e059e")
	poolId := sdk.Bytes32{Val: [2]frontend.Variable{
		frontend.Variable("0x0000000000000000000000000000000000000000"),
		frontend.Variable("0x000000000000000000000001"),
	}}

	// Mock data
	volume := big.NewInt(1000000)
	volatility := big.NewInt(100)
	liquidity := big.NewInt(10000000)
	impermanentLoss := big.NewInt(50000)

	txHash := common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234")

	// Add transaction
	app.AddTransaction(sdk.TransactionData{
		Hash:                txHash,
		ChainId:             big.NewInt(1),
		BlockNum:            big.NewInt(12345678),
		Nonce:               100,
		GasTipCapOrGasPrice: common.Big0,
		GasFeeCap:           common.Big0,
		Value:               common.Big0,
		From:                common.Address{},
		To:                  hookAddress,
		GasLimit:            100000,
	})

	// Add receipt
	app.AddReceipt(sdk.ReceiptData{
		BlockNum: big.NewInt(12345678),
		TxHash:   txHash,
		Fields: [sdk.NumMaxLogFields]sdk.LogFieldData{
			{
				Contract:   hookAddress,
				LogIndex:   0,
				EventID:    poolDataUpdatedEventID,
				IsTopic:    false,
				FieldIndex: 0,
				Value:      common.HexToHash(hexutil.EncodeBig(volume)),
			},
			{
				Contract:   hookAddress,
				LogIndex:   0,
				EventID:    poolDataUpdatedEventID,
				IsTopic:    false,
				FieldIndex: 1,
				Value:      common.HexToHash(hexutil.EncodeBig(volatility)),
			},
			{
				Contract:   hookAddress,
				LogIndex:   0,
				EventID:    poolDataUpdatedEventID,
				IsTopic:    false,
				FieldIndex: 2,
				Value:      common.HexToHash(hexutil.EncodeBig(liquidity)),
			},
			{
				Contract:   hookAddress,
				LogIndex:   0,
				EventID:    poolDataUpdatedEventID,
				IsTopic:    false,
				FieldIndex: 3,
				Value:      common.HexToHash(hexutil.EncodeBig(impermanentLoss)),
			},
		},
	})

	// Initialize AppCircuit
	appCircuit := &AppCircuit{
		PoolId: poolId,
	}
	appCircuitAssignment := &AppCircuit{
		PoolId: poolId,
	}

	in, err := app.BuildCircuitInput(appCircuit)
	if err != nil {
		t.Fatalf("Failed to build circuit input: %v", err)
	}

	test.ProverSucceeded(t, appCircuit, appCircuitAssignment, in)
}

func TestE2E(t *testing.T) {
	app, err := sdk.NewBrevisApp()
	if err != nil {
		t.Fatalf("Failed to create BrevisApp: %v", err)
	}

	poolId := sdk.Bytes32{Val: [2]frontend.Variable{
		frontend.Variable("0x0000000000000000000000000000000000000000"),
		frontend.Variable("0x000000000000000000000001"),
	}}

	// Initialize AppCircuit
	appCircuit := &AppCircuit{
		PoolId: poolId,
	}
	appCircuitAssignment := &AppCircuit{
		PoolId: poolId,
	}

	in, err := app.BuildCircuitInput(appCircuit)
	if err != nil {
		t.Fatalf("Failed to build circuit input: %v", err)
	}

	// Test prover
	test.ProverSucceeded(t, appCircuit, appCircuitAssignment, in)

	// Compile circuit
	outDir := t.TempDir()
	srsDir := t.TempDir()
	compiledCircuit, pk, vk, err := sdk.Compile(appCircuit, outDir, srsDir)
	if err != nil {
		t.Fatalf("Failed to compile circuit: %v", err)
	}

	// Generate proof
	witness, publicWitness, err := sdk.NewFullWitness(appCircuitAssignment, in)
	if err != nil {
		t.Fatalf("Failed to generate witness: %v", err)
	}
	proof, err := sdk.Prove(compiledCircuit, pk, witness)
	if err != nil {
		t.Fatalf("Failed to generate proof: %v", err)
	}

	// Verify proof
	err = sdk.Verify(vk, publicWitness, proof)
	if err != nil {
		t.Fatalf("Failed to verify proof: %v", err)
	}
}