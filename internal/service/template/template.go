package template

import (
    "bytes"
    "fmt"
    "html/template"

    "github.com/eujoy/data-dict/internal/model/domain"
    "github.com/eujoy/data-dict/pkg"
)

const (
    erDiagram = "er"
    html      = "html"
    markdown  = "md"
    mermaid   = "mermaid"
)

// Engine describes the template engine service.
type Engine struct{}

// New creates and returns a new Engine service instance.
func New() *Engine {
    return &Engine{}
}

// Generate and print the respective template.
func (eng *Engine) Generate(outputType string, templateValues domain.TemplateValues) (string, *pkg.Error) {
    switch outputType {
    case erDiagram:
        return eng.generateType(dataDirectoryTemplateERDiagram, templateValues)
    case html:
        return eng.generateType(dataDirectoryTemplateHTML, templateValues)
    case markdown:
        return eng.generateType(dataDirectoryTemplateMarkdown, templateValues)
    case mermaid:
        return eng.generateType(dataDirectoryTemplateMermaid, templateValues)
    default:
        return "", &pkg.Error{Err: fmt.Errorf("invalid output type provided: %v", outputType)}
    }
}

// generateType prepares, generates and print the respective template requested.
func (eng *Engine) generateType(typeTemplate string, templateValues domain.TemplateValues) (string, *pkg.Error) {
    t, templateErr := template.New("template").Parse(typeTemplate)
    if templateErr != nil {
        err := &pkg.Error{Err: templateErr}
        err.LogError()
        return "", err
    }

    var tplData bytes.Buffer
    tmplExecErr := t.Execute(&tplData, templateValues)
    if tmplExecErr != nil {
        err := &pkg.Error{Err: tmplExecErr}
        err.LogError()
        return "", err
    }

    return tplData.String(), nil
}
