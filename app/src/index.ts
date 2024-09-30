import { Brevis, ProofRequest, Prover, asBytes32, asUint248 } from 'brevis-sdk-typescript';
import { ethers } from 'ethers';
import { abi as HOOK_ABI } from './DynamicFeeAdjustmentHook.json';
import axios from 'axios';
import dotenv from 'dotenv';

dotenv.config();

const HOOK_ADDRESS = process.env.HOOK_ADDRESS || '';
const RPC_URL = process.env.RPC_URL || 'https://bsc-testnet.public.blastapi.io';
const PRIVATE_KEY = process.env.PRIVATE_KEY || '';
const BREVIS_ENDPOINT = process.env.BREVIS_ENDPOINT || 'appsdkv2.brevis.network:9094';
const PROVER_ENDPOINT = process.env.PROVER_ENDPOINT || 'localhost:33247';

async function updatePoolFees(poolId: string) {
    const prover = new Prover(PROVER_ENDPOINT);
    const brevis = new Brevis(BREVIS_ENDPOINT);
    const proofReq = new ProofRequest();

    const provider = new ethers.providers.JsonRpcProvider(RPC_URL);
    const signer = new ethers.Wallet(PRIVATE_KEY, provider);
    const hookContract = new ethers.Contract(HOOK_ADDRESS, HOOK_ABI, signer);

    try {
        console.log('Fetching pool data...');
        const poolData = await hookContract.getPoolData(poolId);
        console.log('Pool data fetched successfully');
        console.log('Pool Data:', {
            currentTradingFee: poolData.currentTradingFee.toString(),
            currentLPFee: poolData.currentLPFee.toString(),
            token0Balance: poolData.token0Balance.toString(),
            token1Balance: poolData.token1Balance.toString(),
            totalLiquidity: poolData.totalLiquidity.toString(),
            impermanentLoss: poolData.impermanentLoss.toString(),
            historicalVolumes: poolData.historicalVolumes.map((v: any) => v.toString()),
            historicalVolatilities: poolData.historicalVolatilities.map((v: any) => v.toString()),
            historicalLiquidities: poolData.historicalLiquidities.map((v: any) => v.toString()),
            historicalImpermanentLosses: poolData.historicalImpermanentLosses.map((v: any) => v.toString()),
        });

        console.log('Fetching external market trend...');
        const externalMarketTrend = await fetchExternalMarketTrend();
        console.log('External market trend fetched successfully:', externalMarketTrend);

        const ensureNonZero = (value: ethers.BigNumber) => (value.eq(0) ? ethers.utils.parseUnits('1', 'wei') : value);

        console.log('Preparing proof request...');
        proofReq.setCustomInput({
            PoolId: asBytes32(poolId),
            CurrentTradingFee: asUint248(poolData.currentTradingFee.toString()),
            CurrentLPFee: asUint248(poolData.currentLPFee.toString()),
            Token0Balance: asUint248(ensureNonZero(poolData.token0Balance).toString()),
            Token1Balance: asUint248(ensureNonZero(poolData.token1Balance).toString()),
            HistoricalVolumes: poolData.historicalVolumes.map((v: any) => asUint248(ensureNonZero(v).toString())),
            HistoricalVolatilities: poolData.historicalVolatilities.map((v: any) =>
                asUint248(ensureNonZero(v).toString()),
            ),
            TotalLiquidity: asUint248(ensureNonZero(poolData.totalLiquidity).toString()),
            HistoricalLiquidities: poolData.historicalLiquidities.map((v: any) =>
                asUint248(ensureNonZero(v).toString()),
            ),
            ImpermanentLoss: asUint248(ensureNonZero(poolData.impermanentLoss).toString()),
            HistoricalImpermanentLosses: poolData.historicalImpermanentLosses.map((v: any) =>
                asUint248(ensureNonZero(v).toString()),
            ),
            ExternalMarketTrend: asUint248(externalMarketTrend.toString()),
        });
        console.log('Proof request prepared successfully');
        console.log(`Generating proof for pool ${poolId}`);
        const proofRes = await prover.prove(proofReq);
        console.log('Proof generated successfully', proofRes.proof);

        if (proofRes.has_err) {
            console.error('Proof generation failed:', proofRes.err);
            return;
        }

        console.log('Proof generated successfully');

        console.log('Preparing to submit to Brevis...');
        const brevis_partner_key = process.argv[3] || '';
        console.log(`Using Brevis partner key: ${brevis_partner_key}`);
        console.log(`HOOK_ADDRESS: ${HOOK_ADDRESS}`);

        const brevisRes = await brevis.submit(proofReq, proofRes, 97, 97, 0, brevis_partner_key, HOOK_ADDRESS);
        console.log('Brevis submission result:', brevisRes);

        await brevis.wait(brevisRes.queryKey, 97);

        console.log(`Fee update submitted for pool ${poolId}`);
    } catch (err) {
        console.error(`Error processing pool ${poolId}:`, err);
        if (err instanceof Error) {
            console.error('Error details:', err.message);
            console.error('Stack trace:', err.stack);
        }
    }
}

async function fetchExternalMarketTrend(): Promise<number> {
    try {
        const response = await axios.get('https://api.coingecko.com/api/v3/global');
        const marketData = response.data.data;

        // Calculate trends
        const marketCapChange = marketData.market_cap_change_percentage_24h_usd;

        // Normalize the trend
        return Math.max(0, Math.min(10000, Math.floor((marketCapChange + 10) * 500)));
    } catch (error) {
        console.error('Error fetching external market data:', error);
        return 5000;
    }
}

// Run the update every 1h
const poolId = process.argv[2];
if (!poolId) {
    console.error('Please provide a pool ID as a command line argument.');
    process.exit(1);
}

setInterval(() => updatePoolFees(poolId), 60 * 60 * 1000);
updatePoolFees(poolId);
