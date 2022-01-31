package domain

// TemplateValues describes the details required for the respective values required for the template.
type TemplateValues struct {
    DatabaseName string
    TableList    []TableTmplValue
}

// TableTmplValue describes the table related values for the template.
type TableTmplValue struct {
    TableName       string
    ColumnList      []ColumnTmplValue
    ConstraintsList []ConstraintTmplValue
}

// ColumnTmplValue describes the column related values for the template.
type ColumnTmplValue struct {
    Ordinal      int
    Name         string
    DataType     string
    PK           bool
    FK           bool
    UQ           bool
    NotNull      bool
    DefaultValue string
    Comment      string
}

// ConstraintTmplValue describes the constraint values for the template.
type ConstraintTmplValue struct {
    Name             string
    Type             string
    Column           string
    References       string
    ReferencesTable  string
    ReferencesColumn string
}
