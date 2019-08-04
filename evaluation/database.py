import psycopg2

def connDB(db):
    connection = psycopg2.connect(
        user = "cow",
        password = "123456",
        host = "10.230.12.86",
        port = "5432",
        database = db)
    print "You are connected to - {} \n".format(db)
    return connection, connection.cursor()

def closeDB(connection):
    if(connection):
        connection.close()

def getDataFromDatabase(cursor, query):
    cursor.execute(query)
    return cursor.fetchall()

def updateOrInsertDataToDatabase(cursor, query):
    cursor.execute(query)
    connection.commit()

def getColsOfTable(cursor, table):
    query = "Select * FROM {} LIMIT 1".format(table)
    cursor.execute(query)
    return  [desc[0] for desc in cursor.description]