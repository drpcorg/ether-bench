import json

methods = [
    "eth_call",
    "eth_getTransactionReceipt",
    "eth_getBalance",
	"eth_getBlockByNumber",
	"eth_getBlockByNumber#full",
	"eth_getTransactionCount",
	"eth_blockNumber",
	"eth_getTransactionByHash",
	"eth_getLogs",
	"eth_getCode",
	"eth_estimateGas",
	"eth_getBlockByHash",
	"eth_getBlockByHash#full",
	"eth_getTransactionByBlockNumberAndIndex",
	"net_version",
	"eth_gasPrice",
	"net_listening",
	"net_peerCount",
	"eth_syncing",
	"eth_getStorageAt",
	"eth_accounts",
	"eth_chainId",
	"eth_protocolVersion",
	"eth_feeHistory",
	"eth_maxPriorityFeePerGas",
	"eth_getTransactionByBlockHashAndIndex",
	"eth_getBlockTransactionCountByHash",
	"eth_getBlockTransactionCountByNumber",
	"eth_getBlockReceipts",
	"trace_block",
	"trace_transaction",
	"trace_replayTransaction",
	"trace_ReplayBlockTransactions",
	"debug_TraceTransaction",
	"debug_TraceBlockByNumber",
	"debug_TraceBlockByHash",
	"eth_createAccessList",
	"eth_getProof"
]

steps = []
for i in [10,100,1000,10000]:
    steps.append({
        "duration": 30,
        "rate": i,
    })

result = []
for i in methods:
    stage = {
        "name": f"{i}",
        "profile": {i: 1},
        "steps": steps,
    }
    result.append(stage)

with open("stages.json", "w") as f:
    json.dump(result, f)
