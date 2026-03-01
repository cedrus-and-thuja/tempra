# Tempra

A lightweight command-line tool for applying Go templates to various data sources.

## Overview

Tempra allows you to combine Go templates with data from CSV, JSON, or YAML files to generate text output. It's a simple but powerful tool for generating reports, configuration files, or any text-based output from structured data.

## Installation

```bash
go install github.com/cedrus-and-thuja/tempra/cmd/tempra@latest
```

## Usage

```bash
tempra -template <template_file> -source <data_file>
```

### Required Flags

- `-template`: Path to the Go template file
- `-source`: Path to the data source file (CSV, JSON, or YAML)

## Data Sources

Tempra supports the following data source formats:

### CSV

CSV files are loaded as a list of maps, accessible in the template via `.Data`. Column headers are automatically title-cased (first character uppercase).

Example CSV:
```csv
name,age,city
john,30,new york
jane,25,san francisco
```

Access in template:
```
{{range .Data}}
Name: {{.Name}}, Age: {{.Age}}, City: {{.City}}
{{end}}
```

### JSON

JSON files are loaded directly into the template context.

Example JSON:
```json
{
  "users": [
    { "name": "John", "age": 30 },
    { "name": "Jane", "age": 25 }
  ]
}
```

Access in template:
```
{{range .users}}
Name: {{.name}}, Age: {{.age}}
{{end}}
```

### YAML

YAML files are loaded directly into the template context.

Example YAML:
```yaml
users:
  - name: John
    age: 30
  - name: Jane
    age: 25
```

Access in template:
```
{{range .users}}
Name: {{.name}}, Age: {{.age}}
{{end}}
```

## Examples

### Generate a simple report from CSV

**Data file (people.csv):**
```csv
name,age,occupation
Alice,28,Engineer
Bob,35,Designer
Charlie,42,Manager
```

**Template file (report.template):**
```
# Team Members Report

{{range .Data}}
## {{.Name}}
- Age: {{.Age}}
- Role: {{.Occupation}}
{{end}}
```

**Command:**
```bash
tempra -template report.template -source people.csv
```

**Output:**
```
# Team Members Report

## Alice
- Age: 28
- Role: Engineer

## Bob
- Age: 35
- Role: Designer

## Charlie
- Age: 42
- Role: Manager
```

### Generate a configuration file from JSON

**Data file (config.json):**
```json
{
  "app": {
    "name": "MyApp",
    "version": "1.0.0"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "admin"
  }
}
```

**Template file (app-config.template):**
```
# Application Configuration

APP_NAME={{.app.name}}
APP_VERSION={{.app.version}}

# Database Settings
DB_HOST={{.database.host}}
DB_PORT={{.database.port}}
DB_USER={{.database.username}}
```

**Command:**
```bash
tempra -template app-config.template -source config.json
```

**Output:**
```
# Application Configuration

APP_NAME=MyApp
APP_VERSION=1.0.0

# Database Settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
```

## Development

Tempra uses Taskfile for common development tasks.

- `task setup`: Install development tools
- `task lint`: Run linters
- `task test`: Run tests with coverage
- `task build`: Build the binary
- `task run`: Run the application

## License

[MIT License]