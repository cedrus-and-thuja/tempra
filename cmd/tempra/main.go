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

	"github.com/PaesslerAG/jsonpath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type KV struct {
	Key   string
	Value any
}

func MapToSlice(m map[string]any) []KV {
	var result []KV
	for k, v := range m {
		result = append(result, KV{Key: k, Value: v})
	}
	return result
}

func toJSON(v any) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	switch v := v.(type) {
	case string:
		return v == ""
	case []any:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case complex64:
		return v == 0
	case complex128:
		return v == 0
	case bool:
		return !v
	default:
		return false
	}
}

func hasValue(v any) bool {
	return !isEmpty(v)
}

func main() {
	var templatePath string
	var dataSource string
	flag.StringVar(&templatePath, "template", "", "file to load the template from")
	flag.StringVar(&dataSource, "source", "", "file to load the data from")
	flag.Parse()
	if templatePath != "" && dataSource != "" {
		globalValues := map[string]any{}
		contextData, err := loadData(string(dataSource))
		if err != nil {
			panic(err)
		}
		data, err := os.ReadFile(templatePath)
		if err != nil {
			panic(err)
		}
		funcMap := template.FuncMap{
			"add":  func(a, b int) int { return a + b },
			"sub":  func(a, b int) int { return a - b },
			"mul":  func(a, b int) int { return a * b },
			"div":  func(a, b int) int { return a / b },
			"odd":  func(n int) bool { return n%2 != 0 },
			"even": func(n int) bool { return n%2 == 0 },
			"env": func(key string, defaultValue ...string) string {
				value := os.Getenv(key)
				if value == "" && len(defaultValue) > 0 {
					return defaultValue[0]
				}
				return value
			},
			"sum": func(a []int) int {
				total := 0
				for _, v := range a {
					total += v
				}
				return total
			},
			"poke": func(name string, value any) error {
				globalValues[name] = value
				return nil
			},
			"peek": func(name string) (any, error) {
				value, ok := globalValues[name]
				if !ok {
					return nil, fmt.Errorf("no value found for key: %s", name)
				}
				return value, nil
			},
			"get": func(path string) (any, error) {
				if contextData == nil {
					return nil, errors.New("context data is nil")
				}
				if path == "" || path == "$" || path == "." {
					return contextData, nil
				}
				return jsonpath.Get(path, contextData)
			},
			"getFrom": func(path string, obj any) (any, error) {
				if obj == nil {
					obj = contextData
				}
				if path == "" || path == "$" || path == "." {
					return obj, nil
				}
				return jsonpath.Get(path, obj)
			},
			"isEmpty":    isEmpty,
			"hasValue":   hasValue,
			"toJSON":     toJSON,
			"MapToSlice": MapToSlice,
		}
		tmpl := template.Must(template.New(templatePath).Funcs(funcMap).Option("missingkey=zero").Parse(string(data)))
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
	if strings.HasSuffix(dataSource, "yaml") || strings.HasSuffix(dataSource, "yml") {
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
