# snippetbox

Golang project from Let's Go book by Alex Edwards

### DB

Create a snippetbox DB

```
-- Create a new UTF-8 `snippetbox` database.
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Switch to using the `snippetbox` database.
USE snippetbox;
```

Copy and paste content of file **snippetbox.sql** in your mysql client and run the script.

### Compile main program

This will create an executable specific to your OS in the same folder.

For Linux

```
go build -o snippetbox cmd/web/*
```

For Windows

```
go build -v -o snippetbox.exe cmd/web/*
```

### Run main

You can just run it using the below command on Linux

```
./snippetbox
```

Or, on Windows

```
snippetbox.exe
```

### Compile CLI program

This will create an executable specific to your OS in the same folder.

For Linux

```
go build cmd/cli
```

For Windows

```
go build cmd/cli
```

### Run CLI

You can just run it using the below command on Linux

```
./cli
```

Or, on Windows

```
cli.exe
```

### Execute CLI without compiling

```
go run cmd/cli/main.go -h
```
