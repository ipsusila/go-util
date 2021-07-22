package util

import (
	"errors"
	"reflect"
	"strings"
)

// list of known struct tag key
const (
	JsonTag = "json"
	DbTag   = "db"
	YamlTag = "yaml"
	TomlTag = "toml"
)

// IsZero return true if given parameter is set to zero value
func IsZero(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

// StructTag return associated tag for given field and tag key
func StructTag(v interface{}, fieldName, key string) string {
	vt := reflect.TypeOf(v)
	done := false
	for !done {
		// get information about element
		switch vt.Kind() {
		case reflect.Array, reflect.Slice, reflect.Chan, reflect.Ptr, reflect.Map:
			vt = vt.Elem()
		default:
			done = true
		}
	}

	// if not struct, return false
	if vt.Kind() != reflect.Struct {
		return ""
	}

	if field, ok := vt.FieldByName(fieldName); ok {
		return field.Tag.Get(key)
	}

	return ""
}

// StructTags list tags for given data
func StructTags(v interface{}, key string) ([]string, error) {
	vt := reflect.TypeOf(v)
	done := false
	for !done {
		// get information about element
		switch vt.Kind() {
		case reflect.Array, reflect.Slice, reflect.Chan, reflect.Ptr, reflect.Map:
			vt = vt.Elem()
		default:
			done = true
		}
	}

	// collect tag recursively
	tagsMap := make(map[string]int)

	var collectTags func(t reflect.Type)
	collectTags = func(t reflect.Type) {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Anonymous {
				collectTags(field.Type)
			} else if val, ok := field.Tag.Lookup(key); ok && val != "-" {
				//tags = append(tags, val)
				tagsMap[val] = 1
			}
		}
	}
	collectTags(vt)

	// retur collected tags
	i := 0
	n := len(tagsMap)
	tags := make([]string, n)
	for tag := range tagsMap {
		tags[i] = tag
		i++
	}

	return tags, nil
}

// StructToMap expands struct/map to map[string]interface{}
// If v already a map, it simply return v
func StructToMap(v interface{}) (map[string]interface{}, error) {
	if res, ok := v.(map[string]interface{}); ok {
		return res, nil
	}
	// get value
	vv := reflect.ValueOf(v)
	vv = reflect.Indirect(vv)

	// only handle struct
	if vv.Kind() != reflect.Struct {
		return nil, errors.New("arg v is not struct")
	}

	vt := vv.Type()
	result := make(map[string]interface{})
	for i := 0; i < vv.NumField(); i++ {
		field := vv.Field(i)
		if field.CanInterface() {
			result[vt.Field(i).Name] = field.Interface()
		}
	}

	return result, nil
}

// StructToMapTag expands struct/map to map[string]interface{}
// If v already a map, it simply return v
func StructToMapTag(v interface{}, key string) (map[string]interface{}, error) {
	if res, ok := v.(map[string]interface{}); ok {
		return res, nil
	}

	result := make(map[string]interface{})
	var fnExpander func(reflect.Value)
	fnExpander = func(rv reflect.Value) {
		rv = reflect.Indirect(rv)

		// only handle struct
		if rv.Kind() != reflect.Struct {
			return
		}

		vt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			sf := vt.Field(i)
			sv := rv.Field(i)
			if sf.Anonymous {
				fnExpander(sv)
				continue
			}

			// get tag
			var varName string
			if tag, ok := sf.Tag.Lookup(key); ok {
				if tag == "-" {
					continue // ignore fields
				}
				tags := strings.Split(tag, ",")
				if len(tags) >= 1 {
					varName = tags[0]
				}
			} else {
				varName = sf.Name
			}

			if sv.CanInterface() {
				result[varName] = sv.Interface()
			}
		}
	}

	// expand all values
	fnExpander(reflect.ValueOf(v))

	return result, nil
}

// StringsContains check wether str is in the slice
func StringsContains(strSlice []string, str string) bool {
	for _, v := range strSlice {
		if v == str {
			return true
		}
	}
	return false
}

// StringsContainsFold check wether str is in the slice
func StringsContainsFold(strSlice []string, str string) bool {
	for _, v := range strSlice {
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}

// MapKeys return map keys
func MapKeys(m map[string]interface{}) []string {
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
