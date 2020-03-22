import json
import sys

if len(sys.argv) > 2:
	param, limit =  sys.argv[1], int(sys.argv[2])
	with open("/Users/zaintq/Downloads/diaspora100KCounter", "r") as file:
		counter = []
		for row in file:
			counter.append(json.loads(row))
		counter = sorted(counter, key=lambda x: x[param])
		for c in counter: 
			if c[param] > limit: print(c)
else:
	print("Args: (['nodes', 'edges'], [int])")
	exit(0)