# Ether Bench

This project is designed to benchmark Ethereum clients and RPCs to determine their performance and reliability. The goal is to provide developers with accurate and up-to-date information on the performance of different Ethereum clients and RPCs, allowing them to make informed decisions when choosing which client or RPC to use for their project.

## How it Works
The benchmarking process involves running a series of tests on different Ethereum clients and RPCs, measuring their performance in terms of latency, throuhput etc. The data for tests are providing by [ethspam](https://github.com/p2p-org/ethspam) tool, to provide an accurate representation of how the clients and RPCs perform in a production environment.

## Getting Started

To get started with this project, simply clone the repository and run the following command:
```bash
make build
./ether-bench -t http://10.0.0.1:8545
```

## Arguments

```
-t, --target=   target eth host
-s, --source=   source eth host (default: https://eth.drpc.org/)
-c, --scenario= scenario file (default: stages.json)  
-r, --result=   result file (default: result.json)
```

## Stages config

The file contains an array of objects, each representing a different stage of a benchmarking process. Each stage has a name, a profile, and a series of steps. The profile appears to be a set of key-value pairs, with the key being the name of an Ethereum function and the value being the number of times that function should be called during the stage. The steps appear to be a series of load tests, with each step specifying a duration and a rate. The duration is the length of time the load test should run, and the rate is the number of requests per second that should be sent during the test.

```json
[
  {
    "name": "eth_call",
    "profile": {
      "eth_call": 1
    },
    "steps": [
      {
        "duration": 30,
        "rate": 10
      }
    ]
  },
  {
    "name": "eth_getTransactionReceipt",
    "profile": {
      "eth_getTransactionReceipt": 1
    },
    "steps": [
      {
        "duration": 30,
        "rate": 10
      }
    ]
  }
]
```

## Scripts

### generate_ramp_methods.py
Script generates example stages file for ramping eth methods of node or rpc provider. Can be used for generation of test data.

### calculate_cu.py
Process result json and apply simple model to estimate relative CU's of each method

### chainlist.py
Script generates example stages file for testing public rpc endpoints listed at chainlist for latency.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.