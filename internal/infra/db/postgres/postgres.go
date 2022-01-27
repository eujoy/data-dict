package postgres

import (
    "fmt"
    "time"

    "github.com/eujoy/data-dict/pkg"
    "github.com/gocraft/dbr/v2"
    _ "github.com/lib/pq"
)

// New creates a new connection poll for postgres.
func New(dbHost string, dbPort int, dbName string, dbUser string, dbPass string) (*dbr.Connection, *pkg.Error) {
    psqlInfo := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        dbHost,
        dbPort,
        dbUser,
        dbPass,
        dbName,
    )

    dbConn, err := dbr.Open("postgres", psqlInfo, nil)
    if err != nil {
        return nil, &pkg.Error{Err: err}
    }

    err = dbConn.Ping()
    if err != nil {
        return nil, &pkg.Error{Err: err}
    }

    dbConn.SetMaxOpenConns(5)
    dbConn.SetMaxIdleConns(10)
    dbConn.SetConnMaxLifetime(10 * time.Second)

    return dbConn, nil
}
