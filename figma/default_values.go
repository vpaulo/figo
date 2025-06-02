package figma

import (
	"reflect"
	"strconv"
)

func SetDefaults(v interface{}) {
	setDefaultsRecursive(reflect.ValueOf(v))
}

func setDefaultsRecursive(v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("default")

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if field.String() == "" && tag != "" {
				field.SetString(tag)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() == 0 && tag != "" {
				if val, err := strconv.ParseInt(tag, 10, 64); err == nil {
					field.SetInt(val)
				}
			}
		case reflect.Bool:
			// NOTE: false may be intentional, use *bool instead
			if field.Bool() == false && tag != "" {
				if val, err := strconv.ParseBool(tag); err == nil {
					field.SetBool(val)
				}
			}
		case reflect.Ptr:
			elemKind := field.Type().Elem().Kind()
			if field.IsNil() {
				// Set default value if pointer is nil
				switch elemKind {
				case reflect.Bool:
					if val, err := strconv.ParseBool(tag); err == nil {
						field.Set(reflect.ValueOf(&val))
						continue
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if val, err := strconv.ParseInt(tag, 10, 64); err == nil {
						intVal := reflect.New(field.Type().Elem()).Elem()
						intVal.SetInt(val)
						field.Set(intVal.Addr())
						continue
					}
				case reflect.String:
					strVal := tag
					field.Set(reflect.ValueOf(&strVal))
					continue
				case reflect.Struct:
					field.Set(reflect.New(field.Type().Elem()))
					setDefaultsRecursive(field)
					continue
				}
			} else if elemKind == reflect.Struct {
				setDefaultsRecursive(field)
			}
		case reflect.Struct:
			setDefaultsRecursive(field)
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				// setDefaultsRecursive(elem)
				//
				//
				//
				// Only recurse if the element is a struct or a pointer to a struct
				// elemKind := elem.Kind()
				// if elemKind == reflect.Ptr && elem.Type().Elem().Kind() == reflect.Struct {
				// 	setDefaultsRecursive(elem)
				// } else if elemKind == reflect.Struct {
				// 	setDefaultsRecursive(elem.Addr()) // pass pointer to modify fields
				// }
				//
				//
				//
				// Recursively set defaults for slice of structs or pointers to structs
				switch elem.Kind() {
				case reflect.Ptr:
					if elem.Type().Elem().Kind() == reflect.Struct {
						setDefaultsRecursive(elem)
					}
				case reflect.Struct:
					setDefaultsRecursive(elem.Addr()) // take address to allow setting
				}
			}
		}
	}
}
