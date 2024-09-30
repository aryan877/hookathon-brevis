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
    ExternalMarketTrend       sdk.Uint248  
}

var _ sdk.AppCircuit = &AppCircuit{}

func (c *AppCircuit) Allocate() (maxReceipts, maxSlots, maxTransactions int) {
    return 0, 0, 0
}

func (c *AppCircuit) Define(api *sdk.CircuitAPI, input sdk.DataInput) error {
    u248 := api.Uint248

    // Calculate advanced metrics
    volumeTrend := calculateTrend(api, c.HistoricalVolumes[:])
    volatilityTrend := calculateTrend(api, c.HistoricalVolatilities[:])
    liquidityTrend := calculateTrend(api, c.HistoricalLiquidities[:])
    impermanentLossTrend := calculateTrend(api, c.HistoricalImpermanentLosses[:])

    // Calculate liquidity utilization
    utilization := calculateUtilization(api, c.Token0Balance, c.Token1Balance, c.TotalLiquidity)

    // Define base fees and adjustment factors
    baseTradingFee := sdk.ConstUint248(2000) // 0.2%
    baseLPFee := sdk.ConstUint248(1000) // 0.1%
    maxAdjustment := sdk.ConstUint248(500) // 0.05%

    // Calculate fee adjustments
    tradingFeeAdjustment := calculateFeeAdjustment(api, volumeTrend, volatilityTrend, utilization, c.ExternalMarketTrend)
    lpFeeAdjustment := calculateFeeAdjustment(api, liquidityTrend, impermanentLossTrend, utilization, c.ExternalMarketTrend)

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

func calculateTrend(api *sdk.CircuitAPI, values []sdk.Uint248) sdk.Uint248 {
    u248 := api.Uint248
    recentAvg := average(api, values[20:])
    oldAvg := average(api, values[:20])    
    return u248.Sub(recentAvg, oldAvg)
}

func calculateUtilization(api *sdk.CircuitAPI, token0Balance, token1Balance, totalLiquidity sdk.Uint248) sdk.Uint248 {
    u248 := api.Uint248
    usedLiquidity := u248.Add(token0Balance, token1Balance)
    
    isZero := u248.IsZero(totalLiquidity)
    
    return u248.Select(
        isZero,
        sdk.ConstUint248(0),
        func() sdk.Uint248 {
            quotient, _ := u248.Div(usedLiquidity, totalLiquidity)
            return quotient
        }(),
    )
}

func calculateFeeAdjustment(api *sdk.CircuitAPI, trend1, trend2, utilization, externalTrend sdk.Uint248) sdk.Uint248 {
    u248 := api.Uint248
    internalFactor := u248.Add(trend1, trend2)
    return u248.Add(u248.Mul(internalFactor, sdk.ConstUint248(3)), u248.Add(utilization, externalTrend))
}

func average(api *sdk.CircuitAPI, values []sdk.Uint248) sdk.Uint248 {
    u248 := api.Uint248
    sum := sdk.ConstUint248(0)
    for _, value := range values {
        sum = u248.Add(sum, value)
    }
    
    length := sdk.ConstUint248(uint64(len(values)))
    isZero := u248.IsZero(length)
    
    return u248.Select(
        isZero,
        sdk.ConstUint248(0),
        func() sdk.Uint248 {
            quotient, _ := u248.Div(sum, length)
            return quotient
        }(),
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