package main

import (
    "fmt"
    "os"
    
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

    var outputType string
    var dbHost, dbName, dbUser, dbPass string
    var dbPort int

    var dbConn *dbr.Connection
    tmplEngine := template.New()

    app.Commands = []*cli.Command{
        {
            Name:    "create",
            Aliases: []string{"c"},
            Usage:   "Create the data dictionary from the database.",
            Flags: []cli.Flag{
                &cli.StringFlag{
                    Name:        "output",
                    Aliases:     []string{"o", "O"},
                    Usage:       "Define the output type. Allowed values: ['md', 'html']",
                    Required:    false,
                    Value:       "markdown",
                    Destination: &outputType,
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
            },
            Action: func(c *cli.Context) error {
                dbConn, err = postgres.New(dbHost, dbPort, dbName, dbUser, dbPass)
                if err != nil {
                    err.LogError()
                    return err.Err
                }

                session := dbConn.NewSession(nil)

                repo := postgresRepository.New(dbName, session)
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

                switch outputType {
                case "md":
                    tmplEngine.GenerateMarkdown(templateValues)
                case "html":
                    tmplEngine.GenerateHTML(templateValues)
                default:
                    err = &pkg.Error{Err: fmt.Errorf("invalid value provided for output type: %v", outputType)}
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
