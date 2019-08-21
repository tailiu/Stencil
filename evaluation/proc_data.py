import json

def convertToFloat(data):
    for i, val in enumerate(data):
        data[i] = float(val)
    return data

def convertToJSON(data):
    return json.loads(data)