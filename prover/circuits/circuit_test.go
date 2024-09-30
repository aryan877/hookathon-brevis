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
// 		PoolId:                      sdk.ConstBytes32([]byte("test_pool_id")),
// 		CurrentTradingFee:           sdk.ConstUint248(3000), // 0.3%
// 		CurrentLPFee:                sdk.ConstUint248(1000), // 0.1%
// 		Token0Balance:               sdk.ConstUint248(1000000 * 1e18),
// 		Token1Balance:               sdk.ConstUint248(1000000 * 1e18),
// 		HistoricalVolumes:           [30]sdk.Uint248{},
// 		HistoricalVolatilities:      [30]sdk.Uint248{},
// 		TotalLiquidity:              sdk.ConstUint248(2000000 * 1e18),
// 		HistoricalLiquidities:       [30]sdk.Uint248{},
// 		ImpermanentLoss:             sdk.ConstUint248(50), // 0.5%
// 		HistoricalImpermanentLosses: [30]sdk.Uint248{},
// 		ExternalMarketTrend:         sdk.ConstUint248(100),
// 	}

// 	// Fill historical data
// 	for i := 0; i < 30; i++ {
// 		appCircuit.HistoricalVolumes[i] = sdk.ConstUint248(5000 * 1e18) // Example: 5000 tokens
// 		appCircuit.HistoricalVolatilities[i] = sdk.ConstUint248(100)    // Example value
// 		appCircuit.HistoricalLiquidities[i] = sdk.ConstUint248(2000000 * 1e18)
// 		appCircuit.HistoricalImpermanentLosses[i] = sdk.ConstUint248(50)
// 	}

// 	circuitInput, err := app.BuildCircuitInput(appCircuit)
// 	if err != nil {
// 		t.Fatalf("Failed to build circuit input: %v", err)
// 	}

// 	test.ProverSucceeded(t, appCircuit, appCircuit, circuitInput)
// }

// func TestCircuitEdgeCases(t *testing.T) {
// 	app, err := sdk.NewBrevisApp()
// 	if err != nil {
// 		t.Fatalf("Failed to create Brevis app: %v", err)
// 	}

// 	testCases := []struct {
// 		name    string
// 		circuit *AppCircuit
// 	}{
// 		{
// 			name: "Zero Balances",
// 			circuit: &AppCircuit{
// 				PoolId:                      sdk.ConstBytes32([]byte("test_pool_id")),
// 				CurrentTradingFee:           sdk.ConstUint248(3000),
// 				CurrentLPFee:                sdk.ConstUint248(1000),
// 				Token0Balance:               sdk.ConstUint248(0),
// 				Token1Balance:               sdk.ConstUint248(0),
// 				HistoricalVolumes:           [30]sdk.Uint248{},
// 				HistoricalVolatilities:      [30]sdk.Uint248{},
// 				TotalLiquidity:              sdk.ConstUint248(0),
// 				HistoricalLiquidities:       [30]sdk.Uint248{},
// 				ImpermanentLoss:             sdk.ConstUint248(0),
// 				HistoricalImpermanentLosses: [30]sdk.Uint248{},
// 				ExternalMarketTrend:         sdk.ConstUint248(0),
// 			},
// 		},
// 		{
// 			name: "Max Values",
// 			circuit: &AppCircuit{
// 				PoolId:                      sdk.ConstBytes32([]byte("test_pool_id")),
// 				CurrentTradingFee:           sdk.ConstUint248(10000),
// 				CurrentLPFee:                sdk.ConstUint248(10000),
// 				Token0Balance:               sdk.ConstUint248(1e18),
// 				Token1Balance:               sdk.ConstUint248(1e18),
// 				HistoricalVolumes:           [30]sdk.Uint248{},
// 				HistoricalVolatilities:      [30]sdk.Uint248{},
// 				TotalLiquidity:              sdk.ConstUint248(1e18),
// 				HistoricalLiquidities:       [30]sdk.Uint248{},
// 				ImpermanentLoss:             sdk.ConstUint248(1e18),
// 				HistoricalImpermanentLosses: [30]sdk.Uint248{},
// 				ExternalMarketTrend:         sdk.ConstUint248(10000),
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Fill historical data for edge cases
// 			for i := 0; i < 30; i++ {
// 				tc.circuit.HistoricalVolumes[i] = tc.circuit.Token0Balance
// 				tc.circuit.HistoricalVolatilities[i] = tc.circuit.CurrentTradingFee
// 				tc.circuit.HistoricalLiquidities[i] = tc.circuit.TotalLiquidity
// 				tc.circuit.HistoricalImpermanentLosses[i] = tc.circuit.ImpermanentLoss
// 			}

// 			circuitInput, err := app.BuildCircuitInput(tc.circuit)
// 			if err != nil {
// 				t.Fatalf("Failed to build circuit input: %v", err)
// 			}

// 			test.ProverSucceeded(t, tc.circuit, tc.circuit, circuitInput)
// 		})
// 	}
// }