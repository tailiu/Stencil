import psycopg2

def getDBConn(db):
    conn = psycopg2.connect(dbname=db, user="cow", password="123456", host="10.230.12.86", port="5432")
    cursor = conn.cursor()
    return conn, cursor

def runQueries(dbname):
    conn, cur = getDBConn(dbname)
    conn.autocommit = True

    query_set = {
        "stencil": [
            "truncate table data_bags CASCADE; truncate table date_test CASCADE; truncate table deletion_hold CASCADE; truncate table display_flags CASCADE; truncate table display_registration CASCADE; truncate table error_log CASCADE; truncate table evaluation CASCADE; truncate table identity_table CASCADE; truncate table attribute_changes CASCADE; truncate table migration_registration CASCADE; truncate table reference_table CASCADE; truncate table reference_table_v2 CASCADE; truncate table resolved_references CASCADE; truncate table txn_logs CASCADE; truncate table user_table CASCADE;"],
        "mastodon": [
            "select pg_terminate_backend(pg_stat_activity.pid) from pg_stat_activity where pg_stat_activity.datname in ('mastodon_test', 'mastodon_template') and pid <> pg_backend_pid();",
            "drop database mastodon_test;",
            "create database mastodon_test template mastodon_template owner cow;",],
        "twitter":[
            "select pg_terminate_backend(pg_stat_activity.pid) from pg_stat_activity where pg_stat_activity.datname in ('twitter_template', 'twitter_test') and pid <> pg_backend_pid();",
            "drop database twitter_test;",
            "create database twitter_test with template twitter_template owner cow;",],
        "gnusocial":[
            "select pg_terminate_backend(pg_stat_activity.pid) from pg_stat_activity where pg_stat_activity.datname in ('gnusocial_test', 'gnusocial_template') and pid <> pg_backend_pid();",
            "drop database gnusocial_test;",
            "create database gnusocial_test with template gnusocial_template owner cow;",],
        "diaspora":[
            "select pg_terminate_backend(pg_stat_activity.pid) from pg_stat_activity where pg_stat_activity.datname in ('diaspora_test') and pid <> pg_backend_pid(); ",
            "drop database diaspora_test; ",
            "select pg_terminate_backend(pg_stat_activity.pid) from pg_stat_activity where pg_stat_activity.datname in ('diaspora_100000') and pid <> pg_backend_pid();",
            "create database diaspora_test with template diaspora_100000 owner cow;"]
    }

    for items in query_set.items():
        db = items[0]
        queries = items[1]
        print("***** Resetting %s *****"%db)
        for query in queries:
            print(query)
            cur.execute(query)

if __name__ == "__main__":
    
    dbName  = "stencil_test"
    runQueries(dbName)