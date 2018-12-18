import psycopg2
import random

def getDBConn(db):
    conn = psycopg2.connect(dbname=db, user="root", host="10.224.45.158", port="26257")
    # conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT) 
    return conn, conn.cursor()

def getUsersFromApp(app):
    conn, cur = getDBConn(app)
    cur.execute("SELECT DISTINCT(c_id) FROM customer")
    return [x[0] for x in cur.fetchall()]

def giveAccessToUser(accessGiver, accessHavers):
    tables = ["order", "history", "new_order", "orderline"]


if __name__ == "__main__":
    app = "app1"
    users = getUsersFromApp(app)
    random.shuffle(users)
    accessGivers = users[0 : int(0.3 * len(users))]
    accessHavers = list(set(users) - set(accessGivers))
    for accessGiver in accessGivers: giveAccessToUser(accessGiver, accessHavers)
    print len(accessGivers), len(accessHavers), len(users)
