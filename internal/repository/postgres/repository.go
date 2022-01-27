package postgres

import (
    "github.com/eujoy/data-dict/internal/model/database"
    "github.com/eujoy/data-dict/pkg"
    "github.com/gocraft/dbr/v2"
)

var (
    queryStmtFetchPKConstraints = `
    SELECT
        kcu.column_name AS column_name,
        tco.constraint_name AS constraint_name
    FROM
        information_schema.table_constraints AS tco
        JOIN information_schema.key_column_usage AS kcu
            ON kcu.constraint_name = tco.constraint_name
            AND kcu.constraint_schema = tco.constraint_schema
            AND kcu.constraint_name = tco.constraint_name
    WHERE
        tco.constraint_type = 'PRIMARY KEY'
        AND kcu.table_name = ?`

    queryStmtFetchFKConstraints = `
    SELECT
        tc.constraint_name AS constraint_name,
        tc.table_name AS source_table_name,
        kcu.column_name AS source_column_name,
        ccu.table_name AS foreign_table_name,
        ccu.column_name AS foreign_column_name
    FROM
        information_schema.table_constraints AS tc
        JOIN information_schema.key_column_usage AS kcu
            ON tc.constraint_name = kcu.constraint_name
            AND tc.table_schema = kcu.table_schema
        JOIN information_schema.constraint_column_usage AS ccu
            ON ccu.constraint_name = tc.constraint_name
            AND ccu.table_schema = tc.table_schema
    WHERE
        tc.constraint_type = 'FOREIGN KEY'
        AND tc.table_name = ?`

    queryStmtGenericConstraints = `
    SELECT
        kcu.column_name AS column_name,
        tco.constraint_name AS constraint_name,
        tco.constraint_type AS constraint_type
    FROM
        information_schema.table_constraints AS tco
        JOIN information_schema.key_column_usage AS kcu
            ON kcu.constraint_name = tco.constraint_name
            AND kcu.constraint_schema = tco.constraint_schema
            AND kcu.constraint_name = tco.constraint_name
    WHERE
        tco.constraint_type != 'FOREIGN KEY'
        AND tco.constraint_type != 'PRIMARY KEY'
        AND tco.table_name = ?`
)

// Repo describes the repository structure for the postgres client.
type Repo struct {
    dbName string
    session *dbr.Session
}

// New creates and returns a new repository structure.
func New(dbName string, session *dbr.Session) *Repo {
    return &Repo{
        dbName: dbName,
        session: session,
    }
}

// GetTables retrieves and returns the tables of the database.
func (r *Repo) GetTables() ([]database.TableDef, *pkg.Error) {
    var tableDefList []database.TableDef
    _, execErr := r.session.Select("table_name").
        From("information_schema.tables").
        Where("table_catalog = ?", r.dbName).
        Where("table_schema = ?", "public").
        Where("table_type = ?", "BASE TABLE").
        Load(&tableDefList)
    if execErr != nil {
        err := &pkg.Error{Err: execErr}
        return []database.TableDef{}, err
    }

    return tableDefList, nil
}

// GetColumnsOfTable retrieves and returns tha column details of a table.
func (r *Repo) GetColumnsOfTable(tableName string) ([]database.ColumnDef, *pkg.Error) {
    var columnDefList []database.ColumnDef
    _, execErr := r.session.Select("*").
        From("information_schema.columns").
        Where("table_name = ?", tableName).
        Load(&columnDefList)
    if execErr != nil {
        err := &pkg.Error{Err: execErr}
        return []database.ColumnDef{}, err
    }

    return columnDefList, nil
}

// GetPrimaryKeysOfTable retrieves and returns tha primary key details of a table.
func (r *Repo) GetPrimaryKeysOfTable(tableName string) ([]database.PKConstraintDef, *pkg.Error) {
    var pkConstraintList []database.PKConstraintDef
    _, execErr := r.session.SelectBySql(queryStmtFetchPKConstraints, tableName).
        Load(&pkConstraintList)
    if execErr != nil {
        err := &pkg.Error{Err: execErr}
        return []database.PKConstraintDef{}, err
    }

    return pkConstraintList, nil
}

// GetForeignKeysOfTable retrieves and returns tha foreign key details of a table.
func (r *Repo) GetForeignKeysOfTable(tableName string) ([]database.FKConstraintDef, *pkg.Error) {
    var fkConstraintList []database.FKConstraintDef
    _, execErr := r.session.SelectBySql(queryStmtFetchFKConstraints, tableName).
        Load(&fkConstraintList)
    if execErr != nil {
        err := &pkg.Error{Err: execErr}
        return []database.FKConstraintDef{}, err
    }

    return fkConstraintList, nil
}

// GetGenericConstraintsOfTable retrieves and returns tha generic constraints details of a table.
func (r *Repo) GetGenericConstraintsOfTable(tableName string) ([]database.GenericConstraintDef, *pkg.Error) {
    var genConstraintList []database.GenericConstraintDef
    _, execErr := r.session.SelectBySql(queryStmtGenericConstraints, tableName).
        Load(&genConstraintList)
    if execErr != nil {
        err := &pkg.Error{Err: execErr}
        return []database.GenericConstraintDef{}, err
    }

    return genConstraintList, nil
}
