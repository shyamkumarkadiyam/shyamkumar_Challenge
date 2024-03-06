package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func sanitizeValue(value string, dataType string) interface{} {
	value = strings.TrimSpace(value)

	switch dataType {
	case "S":
		if value == "" {
			return nil
		}
		if match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`, value); match {
			t, _ := time.Parse(time.RFC3339, value)
			return t.Unix()
		}
		return value

	case "N":
		if value == "" || !regexp.MustCompile(`^-?\d+(\.\d+)?$`).MatchString(value) {
			return nil
		}
		value = strings.TrimLeft(value, "0")
		if value == "-" || value == "" {
			value = "0"
		}
		n, _ := json.Number(value).Int64()
		return n

	case "BOOL":
		switch strings.ToUpper(value) {
		case "TRUE", "1", "T":
			return true
		case "FALSE", "0", "F":
			return false
		default:
			return nil
		}

	case "NULL":
		switch strings.ToUpper(value) {
		case "TRUE", "1", "T":
			return nil
		case "FALSE", "0", "F":
			return nil
		default:
			return nil
		}

	case "L":
		var list []interface{}
		l := make(map[string]interface{})
		_ = json.Unmarshal([]byte(value), &l)
		for _, v := range l["L"].([]interface{}) {
			item := v.(map[string]interface{})
			for dataType, value := range item {
				sanitizedValue := sanitizeValue(value.(string), dataType)
				if sanitizedValue != nil {
					list = append(list, sanitizedValue)
				}
			}
		}
		return list

	case "M":
		m := make(map[string]interface{})
		_ = json.Unmarshal([]byte(value), &m)
		sanitizedMap := make(map[string]interface{})
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := m[k]
			item := v.(map[string]interface{})
			for dataType, value := range item {
				sanitizedValue := sanitizeValue(value.(string), dataType)
				if sanitizedValue != nil {
					sanitizedMap[k] = sanitizedValue
				}
			}
		}
		return sanitizedMap
	}

	return nil
}

func transformJSON(inputJSON map[string]interface{}) []map[string]interface{} {
	var output []map[string]interface{}

	for key, value := range inputJSON {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		item := value.(map[string]interface{})
		for dataType, val := range item {
			sanitizedValue := sanitizeValue(val.(string), dataType)
			if sanitizedValue != nil {
				output = append(output, map[string]interface{}{key: sanitizedValue})
				break
			}
		}
	}

	return output
}

func main() {
	var inputJSON map[string]interface{}
	_ = json.NewDecoder(os.Stdin).Decode(&inputJSON)

	outputJSON := transformJSON(inputJSON)

	jsonBytes, _ := json.MarshalIndent(outputJSON, "", "  ")
	fmt.Println(string(jsonBytes))
}
