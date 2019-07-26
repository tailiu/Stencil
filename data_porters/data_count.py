import psycopg2, datetime, time

conn = psycopg2.connect(dbname="stencil", user="cow", host="10.230.12.86", port="5432", password="123456")
cursor = conn.cursor()
while True:
    q = "SELECT count(*) FROM owned_data"
    cursor.execute(q)
    print datetime.datetime.now(), " -- Current rows:", cursor.fetchone()[0]
    time.sleep(10)
    # break