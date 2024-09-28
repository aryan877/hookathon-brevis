package dynamicfee

// import (
// 	"testing"

// 	"github.com/brevis-network/brevis-sdk/sdk"
// 	"github.com/brevis-network/brevis-sdk/test"
// )

// func TestCircuit(t *testing.T) {
// 	app, err := sdk.NewBrevisApp()
// 	if err != nil {
// 		t.Fatalf("Failed to create Brevis app: %v", err)
// 	}

// 	appCircuit := &AppCircuit{
// 		PoolId:                    sdk.ConstBytes32([]byte("test_pool_id")),
// 		CurrentTradingFee:         sdk.ConstUint248(3000), // 0.3%
// 		CurrentLPFee:              sdk.ConstUint248(1000), // 0.1%
// 		Token0Balance:             sdk.ConstUint248(1000000 * 1e18), // 1,000,000 tokens
// 		Token1Balance:             sdk.ConstUint248(1000000 * 1e18), // 1,000,000 tokens
// 		HistoricalVolumes:         [30]sdk.Uint248{sdk.ConstUint248(5000 * 1e18)}, // Example: 5000 tokens
// 		HistoricalVolatilities:    [30]sdk.Uint248{sdk.ConstUint248(100)},        // Example value
// 		TotalLiquidity:            sdk.ConstUint248(2000000 * 1e18), // 2,000,000 tokens
// 		HistoricalLiquidities:     [30]sdk.Uint248{sdk.ConstUint248(2000000 * 1e18)},
// 		ImpermanentLoss:           sdk.ConstUint248(50), // 0.5%
// 		HistoricalImpermanentLosses: [30]sdk.Uint248{sdk.ConstUint248(50)},
// 	}

// 	appCircuitAssignment := &AppCircuit{
// 		PoolId:                    appCircuit.PoolId,
// 		CurrentTradingFee:         appCircuit.CurrentTradingFee,
// 		CurrentLPFee:              appCircuit.CurrentLPFee,
// 		Token0Balance:             appCircuit.Token0Balance,
// 		Token1Balance:             appCircuit.Token1Balance,
// 		HistoricalVolumes:         appCircuit.HistoricalVolumes,
// 		HistoricalVolatilities:    appCircuit.HistoricalVolatilities,
// 		TotalLiquidity:            appCircuit.TotalLiquidity,
// 		HistoricalLiquidities:     appCircuit.HistoricalLiquidities,
// 		ImpermanentLoss:           appCircuit.ImpermanentLoss,
// 		HistoricalImpermanentLosses: appCircuit.HistoricalImpermanentLosses,
// 	}

// 	circuitInput, err := app.BuildCircuitInput(appCircuitAssignment)
// 	if err != nil {
// 		t.Fatalf("Failed to build circuit input: %v", err)
// 	}

// 	test.ProverSucceeded(t, appCircuit, appCircuitAssignment, circuitInput)
// }

// func TestCircuitEdgeCases(t *testing.T) {
// 	app, err := sdk.NewBrevisApp()
// 	if err != nil {
// 		t.Fatalf("Failed to create Brevis app: %v", err)
// 	}

// 	testCases := []struct {
// 		name   string
// 		circuit *AppCircuit
// 	}{
// 		{
// 			name: "Zero Balances",
// 			circuit: &AppCircuit{
// 				PoolId:                    sdk.ConstBytes32([]byte("test_pool_id")),
// 				CurrentTradingFee:         sdk.ConstUint248(3000),
// 				CurrentLPFee:              sdk.ConstUint248(1000),
// 				Token0Balance:             sdk.ConstUint248(0),
// 				Token1Balance:             sdk.ConstUint248(0),
// 				HistoricalVolumes:         [30]sdk.Uint248{sdk.ConstUint248(0)},
// 				HistoricalVolatilities:    [30]sdk.Uint248{sdk.ConstUint248(0)},
// 				TotalLiquidity:            sdk.ConstUint248(0),
// 				HistoricalLiquidities:     [30]sdk.Uint248{sdk.ConstUint248(0)},
// 				ImpermanentLoss:           sdk.ConstUint248(0),
// 				HistoricalImpermanentLosses: [30]sdk.Uint248{sdk.ConstUint248(0)},
// 			},
// 		},
// 		{
// 			name: "Max Values",
// 			circuit: &AppCircuit{
// 				PoolId:                    sdk.ConstBytes32([]byte("test_pool_id")),
// 				CurrentTradingFee:         sdk.ConstUint248(10000), // 1% (max allowed)
// 				CurrentLPFee:              sdk.ConstUint248(10000), // 1% (max allowed)
// 				Token0Balance:             sdk.ConstUint248((1 << 248) - 1),
// 				Token1Balance:             sdk.ConstUint248((1 << 248) - 1),
// 				HistoricalVolumes:         [30]sdk.Uint248{sdk.ConstUint248((1 << 248) - 1)},
// 				HistoricalVolatilities:    [30]sdk.Uint248{sdk.ConstUint248((1 << 248) - 1)},
// 				TotalLiquidity:            sdk.ConstUint248((1 << 248) - 1),
// 				HistoricalLiquidities:     [30]sdk.Uint248{sdk.ConstUint248((1 << 248) - 1)},
// 				ImpermanentLoss:           sdk.ConstUint248((1 << 248) - 1),
// 				HistoricalImpermanentLosses: [30]sdk.Uint248{sdk.ConstUint248((1 << 248) - 1)},
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			circuitInput, err := app.BuildCircuitInput(tc.circuit)
// 			if err != nil {
// 				t.Fatalf("Failed to build circuit input: %v", err)
// 			}

// 			test.ProverSucceeded(t, tc.circuit, tc.circuit, circuitInput)
// 		})
// 	}
// }