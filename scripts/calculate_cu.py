import json
import math
import sys


def main() -> None:
    json_string = sys.stdin.read()
    stages = json.loads(json_string)
    maxlat = 0
    maxsize = 0.0
    maxthrp = 0.0

    for stage in stages:
        stage_name = stage["name"]

        if maxlat < stage["steps"][0]["step_summary"]["latencies"]["99th"] / 1000:
            maxlat = stage["steps"][0]["step_summary"]["latencies"]["99th"] / 1000
        if maxsize < stage["steps"][0]["step_summary"]["bytes_out"]["mean"]:
            maxsize = stage["steps"][0]["step_summary"]["bytes_out"]["mean"]

        for step in stage["steps"]:
            if step["step_summary"]["success"] == 1.0:
                if maxthrp < step["step_summary"]["throughput"]:
                    maxthrp = step["step_summary"]["throughput"]

    results = {}

    for stage in stages:
        stage_name = stage["name"]

        mthrpt = 0.0
        for step in stage["steps"]:
            if mthrpt < step["step_summary"]["throughput"]:
                mthrpt = step["step_summary"]["throughput"]

        thrp = mthrpt / maxthrp
        lats = float(stage["steps"][0]["step_summary"]["latencies"]["99th"] / 1000) / float(maxlat)
        szss = stage["steps"][0]["step_summary"]["bytes_out"]["mean"] / maxsize

        score = 0.3 * (1 - thrp) + 0.5 * lats + 0.2 * szss
        cost = math.floor(score * 100) + 10

        results[stage_name] = cost
    

    print(json.dumps(results, indent=4))


if __name__ == "__main__":
    main()