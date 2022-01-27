package database

// TableDef describes the table related info required as they are retrieved from information_schema.tables.
type TableDef struct {
    TableName string `db:"table_name"`
}

// ColumnDef describes the column related info as they are retrieved from information_schema.columns.
type ColumnDef struct {
    OrdinalPosition int     `db:"ordinal_position"`
    ColumnName      string  `db:"column_name"`
    Default         *string `db:"column_default"`
    IsNullable      string  `db:"is_nullable"` // Can be "YES" or "NO"
    DataType        string  `db:"data_type"`
    UDataType       string  `db:"udt_name"`
    Comment         string  `db:"comment"` // for new this field is empty as we are not fetching it properly.
}

// PKConstraintDef describes the columns that are part of the primary key of a table.
type PKConstraintDef struct {
    ConstraintName string `db:"constraint_name"`
    ColumnName     string `db:"column_name"`
}

// FKConstraintDef describes the definition of the foreign key constraints for a table.
type FKConstraintDef struct {
    ConstraintName    string `db:"constraint_name"`
    SourceTableName   string `db:"source_table_name"`
    SourceColumnName  string `db:"source_column_name"`
    ForeignTableName  string `db:"foreign_table_name"`
    ForeignColumnName string `db:"foreign_column_name"`
}

// GenericConstraintDef describes the generic constraints definition of a table.
type GenericConstraintDef struct {
    ConstraintName string `db:"constraint_name"`
    ColumnName     string `db:"column_name"`
    ConstraintType string `db:"constraint_type"`
}
