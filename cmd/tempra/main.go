package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func main() {
	var templatePath string
	var dataSource string
	flag.StringVar(&templatePath, "template", "", "file to load the template from")
	flag.StringVar(&dataSource, "source", "", "file to load the data from")
	flag.Parse()
	args := flag.Args()
	print(args)
	if templatePath != "" && dataSource != "" {
		data, err := os.ReadFile(templatePath)
		if err != nil {
			panic(err)
		}
		tmpl := template.Must(template.New(templatePath).Parse(string(data)))
		contextData, err := loadData(string(dataSource))
		if err != nil {
			panic(err)
		}
		tmpl.Execute(os.Stdout, contextData)
		return
	}
	panic("template and source are required")
}

func loadData(dataSource string) (*interface{}, error) {
	data, err := os.ReadFile(dataSource)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(dataSource, "json") {
		var obj interface{}
		err = json.Unmarshal(data, &obj)
		return &obj, err
	}
	if strings.HasSuffix(dataSource, "yaml") {
		var obj interface{}
		err = yaml.Unmarshal(data, &obj)
		return &obj, err
	}
	if strings.HasSuffix(dataSource, "csv") {
		reader := csv.NewReader(strings.NewReader(string(data)))
		titler := cases.Title(language.English)
		records := []map[string]string{}
		headers, err := reader.Read()
		for i := range headers {
			headers[i] = titler.String(headers[i])
		}
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Error reading record:", err)
				continue
			}
			recordMap := map[string]string{}
			for i, value := range record {
				recordMap[headers[i]] = value
			}
			records = append(records, recordMap)
		}
		var result interface{} = map[string]interface{}{
			"Data": records,
		}
		return &result, err
	}
	return nil, errors.New("unsupported data source")
}
