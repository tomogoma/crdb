# crdb

Database initialization Utility methods for cockroachDB and Postgres.

## usage

go get the package:

```bash
go get -u github.com/tomogoma/crdb
```

Use to initialize your cockroach/postgres database,
godoc [here](https://godoc.org/github.com/tomogoma/crdb):

```go

// 1. Define a DSN to connect to the DB.
// See https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
// for all DSN options.

// Config can be read from a JSON or YAML file.
config := crdb.Config{
    User: "my_username",
    Password: "my_strong_password",
    DBName: "my_db_name",
    SSLMode: "disable",
}
// config#FormatDSN() yields a string like
//    "user='my_username' password='my_strong_password' dbname='my_db_name' sslmode='disable'"
// You can use crdb.Config or format DSN on your own.
DSN := config.FormatDSN()

// 2. Connect to and ping the DBMS.

db, err := crdb.DBConn(DSN)
if err != nil {
    log.Fatalf("Error establishing connection to DB: %v", err)
}

// 3. Instantiate your database and its tables (if not already instantiated)

tableDescs := []string{
    "CREATE TABLE IF NOT EXISTS foos (name VARCHAR(25))",
    "CREATE TABLE IF NOT EXISTS bars (name VARCHAR(25))",
}
err = crdb.InstantiateDB(db, config.DBName, tableDescs...)
if err != nil {
    log.Fatalf("Error instantiating database: %v", err)
}
```
