package dynamicfee

import (
	"github.com/brevis-network/brevis-sdk/sdk"
)

type AppCircuit struct {
	PoolId sdk.Bytes32
}

var _ sdk.AppCircuit = &AppCircuit{}

func (c *AppCircuit) Allocate() (maxReceipts, maxSlots, maxTransactions int) {
	return 100, 0, 0 // Increased to 100 receipts for more historical data
}

func (c *AppCircuit) Define(api *sdk.CircuitAPI, input sdk.DataInput) error {
	u248 := api.Uint248
	receipts := sdk.NewDataStream(api, input.Receipts)

	// Constants
	hookAddress := api.ToBytes32(sdk.ConstBytes32([]byte{0xCb, 0x38, 0xF6, 0x97, 0x00, 0x54, 0xD3, 0x26, 0xEc, 0xc8, 0x9e, 0xf2, 0x48, 0x62, 0x5b, 0x52, 0x8f, 0xfC, 0xAa, 0x5f}))
	poolDataUpdatedEventID := api.ToUint248(sdk.ConstUint248("0xd0f41fd5b4d393ea3222f2ecd77d99386e8ad292339ad0bbc6e3e5530e5e059e"))

	// Assert that all receipts are from the HOOK_ADDRESS and have the correct event ID
	sdk.AssertEach(receipts, func(receipt sdk.Receipt) sdk.Uint248 {
		return u248.And(
			api.Bytes32.IsEqual(api.ToBytes32(receipt.Fields[0].Contract), hookAddress),
			u248.IsEqual(receipt.Fields[0].EventID, poolDataUpdatedEventID),
		)
	})

	// Extract data from receipts
	historicalVolumes := sdk.Map(receipts, func(receipt sdk.Receipt) sdk.Uint248 {
		return api.ToUint248(receipt.Fields[0].Value)
	})
	historicalVolatilities := sdk.Map(receipts, func(receipt sdk.Receipt) sdk.Uint248 {
		return api.ToUint248(receipt.Fields[1].Value)
	})
	historicalLiquidities := sdk.Map(receipts, func(receipt sdk.Receipt) sdk.Uint248 {
		return api.ToUint248(receipt.Fields[2].Value)
	})
	historicalImpermanentLosses := sdk.Map(receipts, func(receipt sdk.Receipt) sdk.Uint248 {
		return api.ToUint248(receipt.Fields[3].Value)
	})

	// Calculate metrics
	volumeTrend := calculateTrend(api, historicalVolumes)
	volatilityTrend := calculateTrend(api, historicalVolatilities)
	liquidityTrend := calculateTrend(api, historicalLiquidities)
	impermanentLossTrend := calculateTrend(api, historicalImpermanentLosses)

	// Calculate utilization
	utilization := calculateUtilization(api, historicalLiquidities)

	// Define base fees and adjustment factors
	baseTradingFee := sdk.ConstUint248(2000) // 0.2%
	baseLPFee := sdk.ConstUint248(1000) // 0.1%
	maxAdjustment := sdk.ConstUint248(500) // 0.05%

	// Calculate fee adjustments
	externalMarketTrend := sdk.ConstUint248(5000) // Default value
	tradingFeeAdjustment := calculateFeeAdjustment(api, volumeTrend, volatilityTrend, utilization, externalMarketTrend)
	lpFeeAdjustment := calculateFeeAdjustment(api, liquidityTrend, impermanentLossTrend, utilization, externalMarketTrend)

	// Apply adjustments
	newTradingFee := u248.Add(baseTradingFee, u248.Mul(tradingFeeAdjustment, maxAdjustment))
	newLPFee := u248.Add(baseLPFee, u248.Mul(lpFeeAdjustment, maxAdjustment))

	// Ensure fees are within acceptable ranges
	newTradingFee = clampFee(api, newTradingFee)
	newLPFee = clampFee(api, newLPFee)

	// Output results
	api.OutputBytes32(c.PoolId)
	api.OutputUint(24, newTradingFee)
	api.OutputUint(24, newLPFee)

	return nil
}

func calculateTrend(api *sdk.CircuitAPI, values *sdk.DataStream[sdk.Uint248]) sdk.Uint248 {
	u248 := api.Uint248
	recentAvg := sdk.Mean(sdk.RangeUnderlying(values, 0, 10))
	oldAvg := sdk.Mean(sdk.RangeUnderlying(values, 90, 100))
	return u248.Sub(recentAvg, oldAvg)
}

func calculateUtilization(api *sdk.CircuitAPI, liquidities *sdk.DataStream[sdk.Uint248]) sdk.Uint248 {
	u248 := api.Uint248
	currentLiquidity := sdk.GetUnderlying(liquidities, 0)
	maxLiquidity := sdk.Max(liquidities)
	
	isZero := u248.IsZero(maxLiquidity)
	
	quotient, _ := u248.Div(currentLiquidity, maxLiquidity)
	return u248.Select(
		isZero,
		sdk.ConstUint248(0),
		quotient,
	)
}

func calculateFeeAdjustment(api *sdk.CircuitAPI, trend1, trend2, utilization, externalTrend sdk.Uint248) sdk.Uint248 {
	u248 := api.Uint248
	internalFactor := u248.Add(trend1, trend2)
	return u248.Add(u248.Mul(internalFactor, sdk.ConstUint248(3)), u248.Add(utilization, externalTrend))
}

func clampFee(api *sdk.CircuitAPI, fee sdk.Uint248) sdk.Uint248 {
	u248 := api.Uint248
	return u248.Select(
		u248.IsLessThan(fee, sdk.ConstUint248(100)),
		sdk.ConstUint248(100),
		u248.Select(
			u248.IsGreaterThan(fee, sdk.ConstUint248(10000)),
			sdk.ConstUint248(10000),
			fee,
		),
	)
}