import json
import csv
import sys

json_string = sys.stdin.read()
data = json.loads(json_string)

writer = csv.writer(sys.stdout)
writer.writerow(['Name', 'success', 'mean', 'p50', 'p90', 'errors'])

for stage in data:
    name = stage['name']
    success = stage['stage_summary']['success']
    mean = stage['stage_summary']['latencies']['mean'] / 1000000
    p50 = stage['stage_summary']['latencies']['50th'] / 1000000
    p90 = stage['stage_summary']['latencies']['90th'] / 1000000
    errors = ', '.join(stage['stage_summary']['errors'])
    writer.writerow([name, success, mean, p50, p90, errors])
