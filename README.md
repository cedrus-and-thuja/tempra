# Tempra

A super bare bones template cli written in go.

## Installation

```
go install github.com/cedrus-and-thuja/tempra/cmd/tempra@latest
```

## usage

```
tempra -template test.template -source test-data.csv
```

Supports loading text go-templates (html not supported yet). Combined with a data source, supports CSV, json, and YAML data sources. YAML and JSON are passed in as is to the context of the template but CSV is passed in as a list of `map[string]string` as `.Data`. CSV also have all columns with the first character uppercased.
