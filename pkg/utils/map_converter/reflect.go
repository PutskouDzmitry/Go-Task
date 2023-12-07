package map_converter

import "reflect"

func ConvertStructToMap(dto interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.Indirect(reflect.ValueOf(dto))
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i).Tag.Get("json")

		if valueField.Kind() == reflect.Ptr {
			if !valueField.IsNil() {
				result[typeField] = valueField.Elem().Interface()
			}
		} else {
			result[typeField] = valueField.Interface()
		}
	}

	return result
}
