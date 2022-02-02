# data-dict

Create the schema documentation for any database. Currently, it is only PostgreSQL engine supported.

## Usage

To be able to create a diagram/schema dictionary for any database, the following command should be executed:

```shell script
➜ go run cmd/main.go generate -h
NAME:
   main generate - Generate the data dictionary / model representation from the database.

USAGE:
   main generate [command options] [arguments...]

OPTIONS:
   --outputType value, -t value, -T value  Define the output type. Allowed values: ['er', 'html', 'md'] (default: "md")
   --output value, -o value, -O value      Define the output of the generated data. Allowed values: ['std', 'file'] (default: "std")
   --outputFile value, -f value, -F value  Define the output file to publish the data to. This value will be used only in combination when [--output file] is provided. (default: "std")
   --dbHost value, -l value, -L value      Define the host of the database.
   --dbPort value, -p value, -P value      Define the port of the database. (default: 0)
   --dbName value, -n value, -N value      Define the name of the database.
   --dbUser value, -u value, -U value      Define the user of the database.
   --dbPass value, -s value, -S value      Define the password of the database.
   --dbSchema value, -c value, -C value    Define the schema of the database.
   --help, -h                              show help (default: false)
   
➜ 
```

An example for executing this is :

```shell script
➜ go run cmd/main.go generate -l localhost -p 5432 -n my_database -u my_user -s my_password -c public -t html -o file -f file.html
```
