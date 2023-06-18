import json

hosts = [
    "https://endpoints.omniatech.io/v1/eth/mainnet/public",
    "https://rpc.ankr.com/eth",
    "https://eth-mainnet.nodereal.io/v1/1659dfb40aa24bbb8153a677b98064d7",
    "https://ethereum.publicnode.com/",
    "https://1rpc.io/eth",
    "https://rpc.builder0x69.io/",
    "https://rpc.mevblocker.io",
    "https://rpc.flashbots.net/",
    "https://eth.rpc.blxrbdn.com/",
    "https://cloudflare-eth.com/",
    "https://eth-mainnet.public.blastapi.io/",
    "https://api.securerpc.com/v1",
    "https://api.bitstack.com/v1/wNFxbiJyQsSeLrX8RRCHi7NpRxrlErZk/DjShIqLishPCTB9HiMkPHXjUM9CNM9Na/ETH/mainnet",
    "https://eth-rpc.gateway.pokt.network",
    "https://eth-mainnet-public.unifra.io",
    "https://ethereum.blockpi.network/v1/rpc/public",
    "https://rpc.payload.de",
    "https://api.zmok.io/mainnet/oaen6dy8ff6hju9k",
    "https://eth-mainnet.g.alchemy.com/v2/demo",
    "https://eth.api.onfinality.io/public",
    "https://core.gashawk.io/rpc",
    "https://rpc.eth.gateway.fm/",
    "https://rpc.chain49.com/ethereum?api_key=14d1a8b86d8a4b4797938332394203dc",
    "https://eth.meowrpc.com/",
    "https://eth.drpc.org/",
    "https://mainnet.gateway.tenderly.co"
]

profiles = [
    {"name": "getBalance", "profile": {"eth_getBalance": 1}},
    {"name": "getBlockByNumber", "profile": {"eth_getBlockByNumber": 1}},
]

steps = []
for i in [10,100]:
    steps.append({
        "duration": 30,
        "rate": i,
    })

stages = []

for host in hosts:
    for profile in profiles:
        stage = {
            "name": f"{host} {profile['name']}",
            "profile": profile["profile"],
            "steps": steps,
            "target": host
        }
        stages.append(stage)

with open("stages.json", "w") as f:
    json.dump(stages, f, indent=4)

