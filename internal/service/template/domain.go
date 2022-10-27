package template

var (
    dataDirectoryTemplateMermaid = `erDiagram
	{{- range .TableList }}
	{{ .TableName }} {
	{{- range .ColumnList }}
		{{ .DataType }} {{ .Name }} "{{ if .PK }}PK{{ end }}{{ if .FK }}{{ if .PK }}_{{ end }}FK{{ end }}"
	{{- end }}
	}
	{{ end }}
    %% ----- Relationships ----
    {{ range .TableList }}
    {{- $tableName := .TableName }}
    {{- range .ConstraintsList }}
    {{- if .ReferencesTable }}
    %% {{ $tableName }} }o--o{ {{ .ReferencesTable }} : "{{ $tableName }}.{{ .Column }} relates to {{ .ReferencesTable }}.{{ .ReferencesColumn }}"
    {{ $tableName }} }o--o{ {{ .ReferencesTable }} : "{{ .Column }} to {{ .ReferencesColumn }}"
    {{- end }}
    {{- end }}
    {{- end }}
`

    dataDirectoryTemplateERDiagram = `# ER Diagram Definition

title {label:"{{ .DatabaseName }}"}

# Definition of tables.{{print "\n"}}

{{- range .TableList }}
[{{ .TableName }}]
{{- range .ColumnList }}
	{{ if .PK }}*{{ end }}{{ if .FK }}+{{ end }}{{ .Name }} {label:"{{ .DataType }}"}
{{- end }}
{{ end }}
# -----

# Definition of foreign keys.

{{ range .TableList }}
{{- $tableName := .TableName }}
{{- range .ConstraintsList }}
{{- if .ReferencesTable }}{{ $tableName }} *--* {{ .ReferencesTable }} {label:"{{ $tableName }}.{{ .Column }} relates to {{ .ReferencesTable }}.{{ .ReferencesColumn }}"}{{print "\n"}}{{- end }}
{{- end }}
{{- end }}
`

    dataDirectoryTemplateMarkdown = `# Data Directory

Database: {{ .DatabaseName }}

Table of contents
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
| {{ .Ordinal }} | {{ .Name }} | {{ .DataType }} | {{ if .PK }}:heavy_check_mark:{{ end }} | {{ if .FK }}:heavy_check_mark:{{ end }} | {{ if .UQ }}:heavy_check_mark:{{ end }} | {{ if .NotNull }}:heavy_check_mark:{{ end }} | {{ .DefaultValue }} | {{ .Comment }} |
{{- end }}

### Constraints: {{ .TableName }}

| Name | Type | Column(s) | References |
| :--- | :--- | :-------- | :--------- |
{{- range .ConstraintsList }}
| {{ .Name }} | {{ .Type }} | {{ .Column }} | {{ if .ReferencesTable }}[{{ .ReferencesTable }}.{{ .ReferencesColumn }}](#table-{{ .ReferencesTable }}){{ end }} |
{{- end }}

[Top :top:](#data-directory)
{{- end }}
`

    dataDirectoryTemplateHTML = `<html>
    <head>
        <title>Database: {{ .DatabaseName }}</title>
        <style>
            .styled-table {
                border-collapse: collapse;
                margin: 25px 0;
                font-size: 0.9em;
                font-family: sans-serif;
                min-width: 400px;
                box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
            }
    
            .styled-table thead tr {
                background-color: #009879;
                color: #ffffff;
                text-align: left;
            }
    
            .styled-table th,
            .styled-table td {
                padding: 12px 15px;
            }
    
            .styled-table tbody tr {
                border-bottom: 1px solid #dddddd;
            }
            
            .styled-table tbody tr:nth-of-type(even) {
                background-color: #f3f3f3;
            }
            
            .styled-table tbody tr:last-of-type {
                border-bottom: 2px solid #009879;
            }

            .color-with-pseudo {
                list-style: none;
                list-style-position: inside;
            }
            .color-with-pseudo li::before {
                content: "•";
                font-size: 130%;
                line-height: 0;
                margin: 0 0.3rem 0 -0.25rem;
                position: relative;
                top: 0.08rem;
                color: #009879;
            }
        </style>
    </head>
    <body>
        <h1>Data Directory</h1>
        
        <h2>Database: {{ .DatabaseName }}</h2>
        
        <h2 id="top">Table of contents</h2>
        <ul class="color-with-pseudo">
        {{- range .TableList }}
            <li><a href="#table-{{ .TableName }}">Table: {{ .TableName }}</a></li>
                <ul class="color-with-pseudo">
                    <li><a href="#field-details-{{ .TableName }}">Field Details</a></li>
                    <li><a href="#constraints-{{ .TableName }}">Constraints</a></li>
                </ul>
            </li>
        {{- end }}
        </ul>
        
        <br/>
        
        {{- range .TableList }}
        
        <h2 id="table-{{ .TableName }}">Table: {{ .TableName }}</h2>
        
        <h3 id="field-details-{{ .TableName }}">Field Details: {{ .TableName }}</h3>
        
        <table class="styled-table">
            <thead>
                <tr>
                    <th>#</th>
                    <th>Name</th>
                    <th>Data Type</th>
                    <th>PK</th>
                    <th>FK</th>
                    <th>UQ</th>
                    <th>Not null</th>
                    <th>Default Value</th>
                    <th>Description</th>
                </tr>
            </thead>
        {{- range .ColumnList }}
            <tbody>
                <tr>
                    <td style="text-align:center">{{ .Ordinal }}</td>
                    <td style="text-align:left">{{ .Name }}</td>
                    <td style="text-align:left">{{ .DataType }}</td>
                    <td style="text-align:center">{{ if .PK }}&#x2714;{{ end }}</td>
                    <td style="text-align:center">{{ if .FK }}&#x2714;{{ end }}</td>
                    <td style="text-align:center">{{ if .UQ }}&#x2714;{{ end }}</td>
                    <td style="text-align:center">{{ if .NotNull }}&#x2714;{{ end }}</td>
                    <td style="text-align:left">{{ .DefaultValue }}</td>
                    <td style="text-align:left">{{ .Comment }}</td>
                </tr>
            </tbody>
        {{- end }}
        </table>
        
        <h3 id="constraints-{{ .TableName }}">Constraints: {{ .TableName }}</h3>
        
        <table class="styled-table">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Type</th>
                    <th>Column(s)</th>
                    <th>References</th>
                </tr>
            </thead>
        {{- range .ConstraintsList }}
            <tbody>
                <tr>
                    <td style="text-align:left">{{ .Name }}</td>
                    <td style="text-align:left">{{ .Type }}</td>
                    <td style="text-align:left">{{ .Column }}</td>
                    <td style="text-align:left">{{ if .ReferencesTable }}<a href="#table-{{ .ReferencesTable }}">{{ .ReferencesTable }}.{{ .ReferencesColumn }}</a>{{ end }}</td>
                </tr>
            </tbody>
        {{- end }}
        </table>
        
        <a href="#top">[Top &#x21a5;]</a>
        {{- end }}
    </body>
</html>
`
)
