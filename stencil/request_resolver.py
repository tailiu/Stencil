#!/usr/bin/python
import MySQLdb

db = MySQLdb.connect(host="127.0.0.1",    # your host, usually localhost
                     port= 3300,
                     user="root",         # your username
                     passwd="",  # your password
                     db="stencil_storage")        # name of the data base

# you must create a Cursor object. It will let
#  you execute all the queries you need
cur = db.cursor()

# Use all the SQL you like
cur.execute("SELECT * FROM apps")

# print all the first cell of all the rows
for row in cur.fetchall():
    print row

db.close()