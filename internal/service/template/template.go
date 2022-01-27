package template

import (
    "html/template"
    "os"

    "github.com/eujoy/data-dict/internal/model/domain"
    "github.com/eujoy/data-dict/pkg"
)

var (
    dataDirectoryTemplateMarkdown = `# Data Directory

Database: {{ .DatabaseName }}

## Table of contents

Tables:
{{ range .TableList }}
* [Table: {{ .TableName }}](#table-{{ .TableName }})
  * [Field Details](#field-details-{{ .TableName }})
  * [Constraints](#constraints-{{ .TableName }})
{{- end }}

----

{{- range .TableList }}

## Table: {{ .TableName }}

### Field Details: {{ .TableName }}

| #   | Name | Data Type | PK  | FK  | UQ  | Not null | Default Value | Description |
| :-: | :--- | :-------- | :-: | :-: | :-: | :------: | :------------ | :---------- |
{{- range .ColumnList }}
| {{ .Ordinal }} | {{ .Name }} | {{ .DataType }} | {{ .PK }} | {{ .FK }} | {{ .UQ }} | {{ .NotNull }} | {{ .DefaultValue }} | {{ .Comment }} |
{{- end }}

### Constraints: {{ .TableName }}

| Name | Type | Column(s) | References |
| :--- | :--- | :-------- | :--------- |
{{- range .ConstraintsList }}
| {{ .Name }} | {{ .Type }} | {{ .Columns }} | {{ .References }} |
{{- end }}

[Top :top:](#table-of-contents)
{{- end }}
`
)

type Engine struct {}

// New creates and returns a new Engine service instance.
func New() *Engine {
    return &Engine{}
}

// GenerateMarkdown generates and prints the markdown template using the provided values.
func (eng *Engine) GenerateMarkdown(templateValues domain.TemplateValues) *pkg.Error {
    t, templateErr := template.New("template").Parse(dataDirectoryTemplateMarkdown)
    if templateErr != nil {
        err := &pkg.Error{Err: templateErr}
        err.LogError()
        return err
    }

    tmplExecErr := t.Execute(os.Stdout, templateValues)
    if tmplExecErr != nil {
        err := &pkg.Error{Err: tmplExecErr}
        err.LogError()
        return err
    }

    return nil
}
