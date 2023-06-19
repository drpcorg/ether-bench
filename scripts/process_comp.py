import json
import csv
import sys

json_string = sys.stdin.read()
data = json.loads(json_string)

# Create the table header
# table = "| Name | success | mean | p50 | p90 |\n| --- | --- | --- | --- | --- |\n"

# Iterate over each stage and add a row to the table

writer = csv.writer(sys.stdout)
writer.writerow(['Name', 'success', 'mean', 'p50', 'p90'])

for stage in data:
    name = stage['name']
    success = stage['stage_summary']['success']
    mean = stage['stage_summary']['latencies']['mean'] / 1000000
    p50 = stage['stage_summary']['latencies']['50th'] / 1000000
    p90 = stage['stage_summary']['latencies']['90th'] / 1000000
    writer.writerow([name, success, mean, p50, p90])

