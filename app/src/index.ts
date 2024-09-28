import { Brevis, ErrCode, ProofRequest, Prover, asBytes32, asUint248 } from 'brevis-sdk-typescript';
import { ethers } from 'ethers';
import { abi as HOOK_ABI } from './DynamicFeeAdjustmentHook.json';

const HOOK_ADDRESS = 'YOUR_HOOK_CONTRACT_ADDRESS';
const RPC_URL = 'https://bsc-testnet.public.blastapi.io';
const PRIVATE_KEY = 'YOUR_PRIVATE_KEY';
const BREVIS_ENDPOINT = 'appsdkv2.brevis.network:9094';
const PROVER_ENDPOINT = 'localhost:33247';

async function updatePoolFees(poolId: string) {
    const prover = new Prover(PROVER_ENDPOINT);
    const brevis = new Brevis(BREVIS_ENDPOINT);

    const proofReq = new ProofRequest();

    const provider = new ethers.providers.JsonRpcProvider(RPC_URL);
    const signer = new ethers.Wallet(PRIVATE_KEY, provider);
    const hookContract = new ethers.Contract(HOOK_ADDRESS, HOOK_ABI, signer);

    try {
        // Fetch current pool data
        const poolData = await hookContract.getPoolData(poolId);

        // Prepare proof request
        proofReq.setCustomInput({
            PoolId: asBytes32(poolId),
            CurrentTradingFee: asUint248(poolData.currentTradingFee.toString()),
            CurrentLPFee: asUint248(poolData.currentLPFee.toString()),
            Token0Balance: asUint248(poolData.token0Balance.toString()),
            Token1Balance: asUint248(poolData.token1Balance.toString()),
            HistoricalVolumes: poolData.historicalVolumes.map((v: any) => asUint248(v.toString())),
            HistoricalVolatilities: poolData.historicalVolatilities.map((v: any) => asUint248(v.toString())),
            TotalLiquidity: asUint248(poolData.totalLiquidity.toString()),
            HistoricalLiquidities: poolData.historicalLiquidities.map((v: any) => asUint248(v.toString())),
            ImpermanentLoss: asUint248(poolData.impermanentLoss.toString()),
            HistoricalImpermanentLosses: poolData.historicalImpermanentLosses.map((v: any) => asUint248(v.toString())),
        });

        // Generate proof
        console.log(`Sending prove request for pool ${poolId}`);
        const proofRes = await prover.prove(proofReq);

        // Handle errors
        if (proofRes.has_err) {
            const err = proofRes.err;
            switch (err.code) {
                case ErrCode.ERROR_INVALID_INPUT:
                    console.error('Invalid input:', err.msg);
                    break;
                case ErrCode.ERROR_INVALID_CUSTOM_INPUT:
                    console.error('Invalid custom input:', err.msg);
                    break;
                case ErrCode.ERROR_FAILED_TO_PROVE:
                    console.error('Failed to prove:', err.msg);
                    break;
            }
            return;
        }

        console.log('Proof generated successfully');

        // Submit proof to Brevis
        try {
            const brevisRes = await brevis.submit(proofReq, proofRes, 97, 97, 0, '', HOOK_ADDRESS);
            console.log('Brevis submission result:', brevisRes);

            // Wait for the proof to be processed
            await brevis.wait(brevisRes.queryKey, 97);

            console.log(`Fee update submitted for pool ${poolId}`);
        } catch (err) {
            console.error('Error submitting proof to Brevis:', err);
        }
    } catch (err) {
        console.error(`Error processing pool ${poolId}:`, err);
    }
}

// Get the pool ID from command line arguments
const poolId = process.argv[2];

if (!poolId) {
    console.error('Please provide a pool ID as a command line argument.');
    process.exit(1);
}

// Run the update every hour
setInterval(() => updatePoolFees(poolId), 60 * 60 * 1000);

// Initial run
updatePoolFees(poolId);
