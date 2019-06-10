import os
import json
import psycopg2
import re

def getDB():
    conn = psycopg2.connect(
        database='diaspora',
        user='zain',
        sslmode='disable',
        port=26257,
        host='10.230.12.75'
    )
    return conn, conn.cursor()

def replaceSymbolsWithStrings(args_str):
    regex = r"\[:([^,]*),"
    for symbol in re.findall(regex, args_str):
        symbol_str = '"'+symbol+'"'
        args_str = args_str.replace(":"+symbol, symbol_str)
    return args_str


conn, cur = getDB() 

# log_dir = "/home/user/Downloads/diaspora/log/test_logs/"
log_dir = "/home/user/Project/logs/"

# for log_file in os.listdir(log_dir):
for log_file in ['qg_reshare_post.log']:
    print "## PARSING =>", log_file
    with open(log_dir+log_file, "r") as fh:
        log_data = fh.read()
        log_lines = log_data.split("\n")
        for row in log_lines:
            sql, args = "", ""
            identifiers = ["ActiveRecord::Base:", "DEBUG"]
            if all(idf in row for idf in identifiers):
                stmt = row.split("ms)")[1].split("[[")
                sql = stmt[0].strip()
                if len(stmt) > 1:
                    args_str = "[["+stmt[1]
                    try:
                        args = json.loads(args_str)
                    except Exception as e:
                        args_str = replaceSymbolsWithStrings(args_str)
                        args = json.loads(args_str)
                    args = [ arg[1] for arg in args]
                if any(idf in sql for idf in ["UPDATE", "INSERT INTO", "DELETE", "BEGIN", "COMMIT", "ABORT"]):
                    print sql, args
    print "------------------------------------------------------------\n\n\n"
    # break