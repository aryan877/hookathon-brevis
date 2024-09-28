package dynamicfee

import (
	"github.com/brevis-network/brevis-sdk/sdk"
)

type AppCircuit struct {
    PoolId                    sdk.Bytes32
    CurrentTradingFee         sdk.Uint248
    CurrentLPFee              sdk.Uint248
    Token0Balance             sdk.Uint248
    Token1Balance             sdk.Uint248
    HistoricalVolumes         [30]sdk.Uint248
    HistoricalVolatilities    [30]sdk.Uint248
    TotalLiquidity            sdk.Uint248
    HistoricalLiquidities     [30]sdk.Uint248
    ImpermanentLoss           sdk.Uint248
    HistoricalImpermanentLosses [30]sdk.Uint248
}

var _ sdk.AppCircuit = &AppCircuit{}

func (c *AppCircuit) Allocate() (maxReceipts, maxSlots, maxTransactions int) {
    return 0, 0, 0
}

func (c *AppCircuit) Define(api *sdk.CircuitAPI, input sdk.DataInput) error {
    u248 := api.Uint248

    // Calculate averages
    averageVolume := calculateAverage(api, c.HistoricalVolumes[:])
    averageVolatility := calculateAverage(api, c.HistoricalVolatilities[:])
    averageLiquidity := calculateAverage(api, c.HistoricalLiquidities[:])
    averageImpermanentLoss := calculateAverage(api, c.HistoricalImpermanentLosses[:])

    // Define thresholds and base fees
    baseTradingFee := sdk.ConstUint248(2000) // 0.2%
    baseLPFee := sdk.ConstUint248(1000) // 0.1%
    lowVolumeThreshold := sdk.ConstUint248(1000 * 1e18)
    highVolumeThreshold := sdk.ConstUint248(10000 * 1e18)
    lowLiquidityThreshold := sdk.ConstUint248(100000 * 1e18)
    highLiquidityThreshold := sdk.ConstUint248(1000000 * 1e18)
    lowVolatilityThreshold := sdk.ConstUint248(100)  // 1%
    highVolatilityThreshold := sdk.ConstUint248(500) // 5%
    highImpermanentLossThreshold := sdk.ConstUint248(50) // 0.5%

    // Calculate fee adjustments for trading fee
    tradingFeeVolumeAdjustment := calculateAdjustment(api, averageVolume, lowVolumeThreshold, highVolumeThreshold, sdk.ConstUint248(200), sdk.ConstUint248(-100))
    tradingFeeLiquidityAdjustment := calculateAdjustment(api, averageLiquidity, lowLiquidityThreshold, highLiquidityThreshold, sdk.ConstUint248(100), sdk.ConstUint248(-50))
    tradingFeeVolatilityAdjustment := calculateAdjustment(api, averageVolatility, lowVolatilityThreshold, highVolatilityThreshold, sdk.ConstUint248(-100), sdk.ConstUint248(200))

    // Calculate fee adjustments for LP fee
    lpFeeVolumeAdjustment := calculateAdjustment(api, averageVolume, lowVolumeThreshold, highVolumeThreshold, sdk.ConstUint248(-50), sdk.ConstUint248(100))
    lpFeeLiquidityAdjustment := calculateAdjustment(api, averageLiquidity, lowLiquidityThreshold, highLiquidityThreshold, sdk.ConstUint248(100), sdk.ConstUint248(-50))
    lpFeeImpermanentLossAdjustment := u248.Select(
        u248.IsGreaterThan(averageImpermanentLoss, highImpermanentLossThreshold),
        sdk.ConstUint248(200),
        sdk.ConstUint248(0),
    )

    // Calculate new fees
    newTradingFee := u248.Add(baseTradingFee, u248.Add(tradingFeeVolumeAdjustment, u248.Add(tradingFeeLiquidityAdjustment, tradingFeeVolatilityAdjustment)))
    newLPFee := u248.Add(baseLPFee, u248.Add(lpFeeVolumeAdjustment, u248.Add(lpFeeLiquidityAdjustment, lpFeeImpermanentLossAdjustment)))

    // Ensure fees are within acceptable ranges (0.01% to 1% for each)
    newTradingFee = clampFee(api, newTradingFee)
    newLPFee = clampFee(api, newLPFee)

    // Output results
    api.OutputBytes32(c.PoolId)
    api.OutputUint(24, newTradingFee)
    api.OutputUint(24, newLPFee)

    return nil
}

func calculateAverage(api *sdk.CircuitAPI, values []sdk.Uint248) sdk.Uint248 {
    u248 := api.Uint248
    sum := sdk.ConstUint248(0)
    for _, value := range values {
        sum = u248.Add(sum, value)
    }
    result, _ := u248.Div(sum, sdk.ConstUint248(uint64(len(values))))
    return result
}

func calculateAdjustment(api *sdk.CircuitAPI, value, lowThreshold, highThreshold, lowAdjustment, highAdjustment sdk.Uint248) sdk.Uint248 {
    u248 := api.Uint248
    return u248.Select(
        u248.IsLessThan(value, lowThreshold),
        lowAdjustment,
        u248.Select(
            u248.IsGreaterThan(value, highThreshold),
            highAdjustment,
            sdk.ConstUint248(0),
        ),
    )
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