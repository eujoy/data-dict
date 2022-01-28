package template

import (
    "html/template"
    "os"

    "github.com/eujoy/data-dict/internal/model/domain"
    "github.com/eujoy/data-dict/pkg"
)

// Engine describes the template engine service.
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

// GenerateHTML generates and prints the markdown template using the provided values.
func (eng *Engine) GenerateHTML(templateValues domain.TemplateValues) *pkg.Error {
    t, templateErr := template.New("template").Parse(dataDirectoryTemplateHTML)
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
