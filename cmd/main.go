package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/eujoy/data-dict/internal/config"
    "github.com/eujoy/data-dict/internal/infra/db/postgres"
    postgresRepository "github.com/eujoy/data-dict/internal/repository/postgres"
    "github.com/eujoy/data-dict/internal/service/decorator"
    "github.com/eujoy/data-dict/internal/service/template"
    "github.com/eujoy/data-dict/pkg"
    "github.com/gocraft/dbr/v2"
    "github.com/urfave/cli/v2"
)

const (
    configurationFileName = "configuration.yaml"
)

func main() {
    cfg, err := config.New(configurationFileName)
    if err != nil {
        err.LogError()
    }

    var app = cli.NewApp()
    info(app, cfg)

    var output, outputType, outputFile string
    var dbHost, dbName, dbUser, dbPass, dbSchema string
    var dbPort int

    var dbConn *dbr.Connection
    tmplEngine := template.New()

    app.Commands = []*cli.Command{
        {
            Name:    "generate",
            Aliases: []string{"gen"},
            Usage:   "Generate the data dictionary / model representation from the database.",
            Flags: []cli.Flag{
                &cli.StringFlag{
                    Name:        "outputType",
                    Aliases:     []string{"t", "T"},
                    Usage:       "Define the output type. Allowed values: ['er', 'html', 'md', 'mermaid']",
                    Required:    false,
                    Value:       "mermaid",
                    Destination: &outputType,
                },
                &cli.StringFlag{
                    Name:        "output",
                    Aliases:     []string{"o", "O"},
                    Usage:       "Define the output of the generated data. Allowed values: ['std', 'file']",
                    Required:    false,
                    Value:       "std",
                    Destination: &output,
                },
                &cli.StringFlag{
                    Name:        "outputFile",
                    Aliases:     []string{"f", "F"},
                    Usage:       "Define the output file to publish the data to. This value will be used only in combination when [--output file] is provided.",
                    Required:    false,
                    Value:       "std",
                    Destination: &outputFile,
                },
                &cli.StringFlag{
                    Name:        "dbHost",
                    Aliases:     []string{"l", "L"},
                    Usage:       "Define the host of the database.",
                    Required:    true,
                    Destination: &dbHost,
                },
                &cli.IntFlag{
                    Name:        "dbPort",
                    Aliases:     []string{"p", "P"},
                    Usage:       "Define the port of the database.",
                    Required:    true,
                    Destination: &dbPort,
                },
                &cli.StringFlag{
                    Name:        "dbName",
                    Aliases:     []string{"n", "N"},
                    Usage:       "Define the name of the database.",
                    Required:    true,
                    Destination: &dbName,
                },
                &cli.StringFlag{
                    Name:        "dbUser",
                    Aliases:     []string{"u", "U"},
                    Usage:       "Define the user of the database.",
                    Required:    true,
                    Destination: &dbUser,
                },
                &cli.StringFlag{
                    Name:        "dbPass",
                    Aliases:     []string{"s", "S"},
                    Usage:       "Define the password of the database.",
                    Required:    true,
                    Destination: &dbPass,
                },
                &cli.StringFlag{
                    Name:        "dbSchema",
                    Aliases:     []string{"c", "C"},
                    Usage:       "Define the schema of the database.",
                    Required:    false,
                    Value:       "public",
                    Destination: &dbSchema,
                },
            },
            Action: func(c *cli.Context) error {
                dbConn, err = postgres.New(dbHost, dbPort, dbName, dbUser, dbPass)
                if err != nil {
                    err.LogError()
                    return err.Err
                }

                session := dbConn.NewSession(nil)
                repo := postgresRepository.New(dbName, dbSchema, session)
                decoratorService := decorator.New(repo, dbName)

                templateValues, err := decoratorService.GetTables().
                    GetColumnsOfAllTables().
                    GetPrimaryKeyOfAllTables().
                    GetForeignKeyOfAllTables().
                    GetGenericConstraintsOfAllTables().
                    PrepareTemplateValues()
                if err != nil {
                    err.LogError()
                    return err.Err
                }

                var generatedData string
                generatedData, err = tmplEngine.Generate(outputType, templateValues)
                if err != nil {
                    err.LogError()
                    return err.Err
                }

                switch output {
                case "std":
                    fmt.Print(generatedData)
                case "file":
                    fileExtension := filepath.Ext(outputFile)
                    if fmt.Sprintf(".%v", outputType) != fileExtension {
                        err = &pkg.Error{Err: fmt.Errorf("incompatible types provided for output type '%v' and file extention '%v'", outputType, fileExtension)}
                        err.LogError()
                        return err.Err
                    }

                    fileWriteErr := ioutil.WriteFile(outputFile, []byte(generatedData), 0755)
                    if fileWriteErr != nil {
                        err = &pkg.Error{Err: fmt.Errorf("failed to write data to file '%v' with error: %v", outputFile, fileWriteErr)}
                        err.LogError()
                        return err.Err
                    }
                default:
                    err = &pkg.Error{Err: fmt.Errorf("invalid output source was provided: %v", output)}
                    err.LogError()
                    return err.Err
                }

                return nil
            },
        },
    }

    execErr := app.Run(os.Args)
    if execErr != nil {
        os.Exit(1)
    }
}

// info sets up the information of the tool.
func info(app *cli.App, cfg *config.Config) {
    var appAuthors []*cli.Author
    for _, author := range cfg.Application.Authors {
        newAuthor := cli.Author{
            Name:  author.Name,
            Email: author.Email,
        }
        appAuthors = append(appAuthors, &newAuthor)
    }

    app.Authors = appAuthors
    app.Name = cfg.Application.Name
    app.Usage = cfg.Application.Usage
    app.Version = cfg.Application.Version
}
