package decorator

import (
    "sort"
    "strings"

    "github.com/eujoy/data-dict/internal/model/database"
    "github.com/eujoy/data-dict/internal/model/domain"
    "github.com/eujoy/data-dict/pkg"
)

type repo interface {
    GetTables() ([]database.TableDef, *pkg.Error)
    GetColumnsOfTable(tableName string) ([]database.ColumnDef, *pkg.Error)
    GetPrimaryKeysOfTable(tableName string) ([]database.PKConstraintDef, *pkg.Error)
    GetForeignKeysOfTable(tableName string) ([]database.FKConstraintDef, *pkg.Error)
    GetGenericConstraintsOfTable(tableName string) ([]database.GenericConstraintDef, *pkg.Error)
}

// Service describes the decorator service for preparing and generating the template values.
type Service struct {
    repo repo
    err  *pkg.Error

    databaseName            string
    tableDefList            []database.TableDef
    columnDefMap            map[string][]database.ColumnDef
    primaryKeyDefMap        map[string][]database.PKConstraintDef
    foreignKeyDefMap        map[string][]database.FKConstraintDef
    genericConstraintDefMap map[string][]database.GenericConstraintDef
}

// New creates and returns a new decorator service.
func New(repo repo, databaseName string) *Service {
    columnDefMap := make(map[string][]database.ColumnDef)
    primaryKeyDefMap := make(map[string][]database.PKConstraintDef)
    foreignKeyDefMap := make(map[string][]database.FKConstraintDef)
    genericConstraintDefMap := make(map[string][]database.GenericConstraintDef)

    return &Service{
        repo:                    repo,
        databaseName:            databaseName,
        tableDefList:            []database.TableDef{},
        columnDefMap:            columnDefMap,
        primaryKeyDefMap:        primaryKeyDefMap,
        foreignKeyDefMap:        foreignKeyDefMap,
        genericConstraintDefMap: genericConstraintDefMap,
    }
}

// GetTables retrieves and returns the tables of the database.
func (s *Service) GetTables() *Service {
    if s.err != nil {
        return s
    }

    tableDefList, err := s.repo.GetTables()
    if err != nil {
        s.err = err
        return s
    }

    s.tableDefList = tableDefList
    return s
}

// GetColumnsOfAllTables retrieves all the columns for all the tables that have been already retrieved.
func (s *Service) GetColumnsOfAllTables() *Service {
    if s.err != nil {
        return s
    }

    for _, tb := range s.tableDefList {
        columnDefList, err := s.repo.GetColumnsOfTable(tb.TableName)
        if err != nil {
            s.err = err
            return s
        }

        s.columnDefMap[tb.TableName] = columnDefList
    }

    return s
}

// GetPrimaryKeyOfAllTables retrieves all the primary key details for all the tables that have been already retrieved.
func (s *Service) GetPrimaryKeyOfAllTables() *Service {
    if s.err != nil {
        return s
    }

    for _, tb := range s.tableDefList {
        primaryKeyDefList, err := s.repo.GetPrimaryKeysOfTable(tb.TableName)
        if err != nil {
            s.err = err
            return s
        }

        s.primaryKeyDefMap[tb.TableName] = primaryKeyDefList
    }

    return s
}

// GetForeignKeyOfAllTables retrieves all the foreign key details for all the tables that have been already retrieved.
func (s *Service) GetForeignKeyOfAllTables() *Service {
    if s.err != nil {
        return s
    }

    for _, tb := range s.tableDefList {
        foreignKeyDefList, err := s.repo.GetForeignKeysOfTable(tb.TableName)
        if err != nil {
            s.err = err
            return s
        }

        s.foreignKeyDefMap[tb.TableName] = foreignKeyDefList
    }

    return s
}

// GetGenericConstraintsOfAllTables retrieves all the generic constraints details for all the tables that have been already retrieved.
func (s *Service) GetGenericConstraintsOfAllTables() *Service {
    if s.err != nil {
        return s
    }

    for _, tb := range s.tableDefList {
        genericConstraintsDefList, err := s.repo.GetGenericConstraintsOfTable(tb.TableName)
        if err != nil {
            s.err = err
            return s
        }

        s.genericConstraintDefMap[tb.TableName] = genericConstraintsDefList
    }

    return s
}

// PrepareTemplateValues prepares and returns the template values based on the fetched information.
func (s *Service) PrepareTemplateValues() (domain.TemplateValues, *pkg.Error) {
    if s.err != nil {
        return domain.TemplateValues{}, s.err
    }

    templateValues := domain.TemplateValues{DatabaseName: s.databaseName}
    for _, tb := range s.tableDefList {
        var constraintsList []domain.ConstraintTmplValue

        for _, pk := range s.primaryKeyDefMap[tb.TableName] {
            constr := domain.ConstraintTmplValue{
                Name:             pk.ConstraintName,
                Type:             "PRIMARY KEY",
                Column:           pk.ColumnName,
                ReferencesTable:  "",
                ReferencesColumn: "",
            }

            constraintsList = append(constraintsList, constr)
        }

        for _, fk := range s.foreignKeyDefMap[tb.TableName] {
            constr := domain.ConstraintTmplValue{
                Name:             fk.ConstraintName,
                Type:             "FOREIGN KEY",
                Column:           fk.SourceColumnName,
                ReferencesTable:  fk.ForeignTableName,
                ReferencesColumn: fk.ForeignColumnName,
            }

            constraintsList = append(constraintsList, constr)
        }

        for _, gen := range s.genericConstraintDefMap[tb.TableName] {
            constr := domain.ConstraintTmplValue{
                Name:             gen.ConstraintName,
                Type:             gen.ConstraintType,
                Column:           gen.ColumnName,
                ReferencesTable:  "",
                ReferencesColumn: "",
            }

            constraintsList = append(constraintsList, constr)
        }

        var columnList []domain.ColumnTmplValue
        for _, col := range s.columnDefMap[tb.TableName] {
            defaultVal := ""
            if col.Default != nil {
                defaultVal = *col.Default
            }

            commentVal := ""
            if col.Comment != nil {
                commentVal = *col.Comment
            }

            colTmplVal := domain.ColumnTmplValue{
                Ordinal:      col.OrdinalPosition,
                Name:         col.ColumnName,
                DataType:     strings.Replace(col.UDataType, "_", "", -1),
                PK:           s.getPKValueForColumn(tb.TableName, col.ColumnName),
                FK:           s.getFKValueForColumn(tb.TableName, col.ColumnName),
                UQ:           s.getUQValueForColumn(tb.TableName, col.ColumnName),
                NotNull:      col.IsNullable == "YES",
                DefaultValue: defaultVal,
                Comment:      commentVal,
            }

            columnList = append(columnList, colTmplVal)
        }

        sort.Slice(columnList, func(i int, j int) bool {
            return columnList[i].Ordinal < columnList[j].Ordinal
        })

        sort.Slice(constraintsList, func(i int, j int) bool {
            return constraintsList[i].Name < constraintsList[j].Name
        })

        templateValues.TableList = append(templateValues.TableList, domain.TableTmplValue{
            TableName:       tb.TableName,
            ColumnList:      columnList,
            ConstraintsList: constraintsList,
        })
    }

    sort.Slice(templateValues.TableList, func(i int, j int) bool {
        return templateValues.TableList[i].TableName < templateValues.TableList[j].TableName
    })

    return templateValues, nil
}

func (s *Service) getPKValueForColumn(tableName string, columnName string) bool {
    for _, pk := range s.primaryKeyDefMap[tableName] {
        if columnName == pk.ColumnName {
            return true
        }
    }

    return false
}

func (s *Service) getFKValueForColumn(tableName string, columnName string) bool {
    for _, fk := range s.foreignKeyDefMap[tableName] {
        if columnName == fk.SourceColumnName {
            return true
        }
    }

    return false
}

func (s *Service) getUQValueForColumn(tableName string, columnName string) bool {
    for _, gen := range s.genericConstraintDefMap[tableName] {
        if columnName == gen.ColumnName && gen.ConstraintType == "UNIQUE" {
            return true
        }
    }

    return false
}
