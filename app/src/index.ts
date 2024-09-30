import { Brevis, ProofRequest, Prover, asBytes32, ReceiptData, Field } from 'brevis-sdk-typescript';
import { ethers } from 'ethers';
import { abi as HOOK_ABI } from './DynamicFeeAdjustmentHook.json';
import dotenv from 'dotenv';
import axios from 'axios';

dotenv.config();

const HOOK_ADDRESS = process.env.HOOK_ADDRESS || '';
const RPC_URL = process.env.RPC_URL || 'https://bsc-testnet.public.blastapi.io';
const PRIVATE_KEY = process.env.PRIVATE_KEY || '';
const BREVIS_ENDPOINT = process.env.BREVIS_ENDPOINT || 'appsdkv2.brevis.network:9094';
const PROVER_ENDPOINT = process.env.PROVER_ENDPOINT || 'localhost:33247';
const BSC_SCAN_API_KEY = process.env.BSC_SCAN_API_KEY || '';

async function updatePoolFees(poolId: string) {
    const prover = new Prover(PROVER_ENDPOINT);
    const brevis = new Brevis(BREVIS_ENDPOINT);
    const proofReq = new ProofRequest();

    const provider = new ethers.providers.JsonRpcProvider(RPC_URL);
    const signer = new ethers.Wallet(PRIVATE_KEY, provider);
    const hookContract = new ethers.Contract(HOOK_ADDRESS, HOOK_ABI, signer);

    try {
        console.log('Fetching historical events...');
        const apiUrl = `https://api-testnet.bscscan.com/api?module=logs&action=getLogs&fromBlock=0&address=${HOOK_ADDRESS}&topic0=${ethers.utils.id(
            'PoolDataUpdated(bytes32,uint256,uint256,uint256,uint256)',
        )}&topic0_1_opr=and&topic1=${ethers.utils.hexZeroPad(poolId, 32)}&apikey=${BSC_SCAN_API_KEY}`;

        const response = await axios.get(apiUrl);
        const events = response.data.result;

        console.log(`Fetched ${events.length} events`);

        if (events.length === 0) {
            console.log('No events found. Skipping proof generation.');
            return;
        }

        console.log('Preparing proof request...');
        proofReq.setCustomInput({
            PoolId: asBytes32(poolId),
        });

        for (const event of events) {
            const receiptData = new ReceiptData({
                block_num: parseInt(event.blockNumber, 16),
                tx_hash: event.transactionHash,
                fields: [
                    new Field({
                        contract: HOOK_ADDRESS,
                        log_index: parseInt(event.logIndex, 16),
                        event_id: event.topics[0],
                        value: event.data,
                        is_topic: false,
                        field_index: 0,
                    }),
                ],
            });
            proofReq.addReceipt(receiptData);
        }

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

// Run the update every 1h
const poolId = process.argv[2];
if (!poolId) {
    console.error('Please provide a pool ID as a command line argument.');
    process.exit(1);
}

setInterval(() => updatePoolFees(poolId), 60 * 60 * 1000);
updatePoolFees(poolId);
