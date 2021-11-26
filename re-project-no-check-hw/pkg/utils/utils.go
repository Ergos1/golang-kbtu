package utils

import (
	"fmt"
	"reflect"
)

func IsDefaultValue(value interface{}) bool {
	return value == reflect.Zero(reflect.TypeOf(value)).Interface()
}

func DBValueByField(arg interface{}, _field string) string {
	value := ""

	v := reflect.ValueOf(arg)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		fallthrough
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i).Tag.Get("db")
			if field == _field {
				value = fmt.Sprintf("%v", v.Field(i).Interface())
				break
			}
		}
	default:
		panic(fmt.Errorf("[error] DBValues requires a struct, found: %s", v.Kind().String()))
	}

	return value
}

func DBValues(arg interface{}) []string {
	values := make([]string, 0)

	v := reflect.ValueOf(arg)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		fallthrough
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			value := v.Field(i).Interface()
			if !IsDefaultValue(value) {
				values = append(values, fmt.Sprintf("%v", value))
			}
		}
	default:
		panic(fmt.Errorf("[error] DBValues requires a struct, found: %s", v.Kind().String()))
	}

	return values
}

func DBFields(arg interface{}) []string {
	fields := make([]string, 0)

	v := reflect.ValueOf(arg)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		fallthrough
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i).Tag.Get("db")
			if field != "" {
				fields = append(fields, field)
			}
		}
	default:
		panic(fmt.Errorf("[error] DBFields requires a struct, found: %s", v.Kind().String()))
	}

	return fields
}

func SetCsv(fields []string, values []string) string {
	if len(fields) != len(values) {
		panic(fmt.Errorf("[error] SetCsv get args with different len"))
	}

	result := ""
	for i := 0; i < len(fields); i++ {
		result += fields[i] + "=" + values[i] + ", "
	}
	result = result[:len(result)-2]
	return result
}

func SetCsvIgnoreId(fields []string, values []string) string {
	if len(fields) != len(values) {
		panic(fmt.Errorf("[error] SetCsv get args with different len"))
	}

	result := ""
	for i := 0; i < len(fields); i++ {
		if fields[i] == "id" {
			continue
		}
		result += fields[i] + "=" + values[i] + ", "
	}
	result = result[:len(result)-2]
	return result
}

func FieldsCSVIgnoreId(fields []string) string {
	result := ""
	for _, field := range fields {
		if field == "id"{
			continue
		}
		result += field + ", "
	}
	result = result[:len(result)-2]
	return result
}

func FieldsCSV(fields []string) string {
	result := ""
	for _, field := range fields {
		result += field + ", "
	}
	result = result[:len(result)-2]
	return result
}

func FieldsCSVColonsIgnoreId(fields []string) string {
	result := ""
	for _, field := range fields {
		if field == "id" {
			continue
		}
		result += ":" + field + ", "
	}
	result = result[:len(result)-2]
	return result
}

func FieldsCSVColons(fields []string) string {
	result := ""
	for _, field := range fields {
		result += ":" + field + ", "
	}
	result = result[:len(result)-2]
	return result
}
