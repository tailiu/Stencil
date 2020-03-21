import json
with open("/Users/zaintq/Downloads/diaspora100KCounter", "r") as file:
	param =  'edges' #'nodes'
	counter = []
	for row in file:
		counter.append(json.loads(row))
	counter = sorted(counter, key=lambda x: x[param])
	for c in counter: 
		if c[param] > 1000:
			print(c)
