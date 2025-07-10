package encoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// since a valid data struct of Hive is an array of different types, and Go
// does not like that, this function deserialize the value of each element in
// the array as struct exported field, in the order they are defined.
func JsonArrayDeserialize(buf any, rawJson []byte) error {
	bufPtr := reflect.ValueOf(buf)
	if bufPtr.Kind() != reflect.Ptr ||
		bufPtr.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("buf must be a pointer to a struct")
	}

	// deserialize json
	rawValues := []any{}
	if err := json.Unmarshal(rawJson, &rawValues); err != nil {
		return fmt.Errorf("failed to json deserialize AccountAuth: %v", err)
	}

	bufV := bufPtr.Elem()

	// check for pub fields that can be set
	canSetFieldCounter := 0
	canSetFieldIndexes := make([]int, 0, bufV.NumField())
	for i := 0; i < bufV.NumField(); i++ {
		if bufV.Field(i).CanSet() {
			canSetFieldCounter++
			canSetFieldIndexes = append(canSetFieldIndexes, i)
		}
	}

	if len(rawValues) != canSetFieldCounter {
		return errors.New("array item and struct pub field count mismatch.")
	}

	// iterate + set fields
	for i, v := range rawValues {
		field := bufV.Field(
			canSetFieldIndexes[i],
		) // skipping over private fields
		if v == nil {
			field.Set(reflect.Zero(field.Type()))
			continue
		}

		err := setFieldValue(field, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, value any) error {
	valueV := reflect.ValueOf(value)

	if valueV.Type().AssignableTo(field.Type()) {
		field.Set(valueV)
		return nil
	}

	// thank you claude
	switch field.Kind() {

	case reflect.String:
		if str, ok := value.(string); ok {
			field.SetString(str)
		} else {
			return fmt.Errorf("cannot convert %T to string", value)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case int:
			field.SetInt(int64(v))
		case int32:
			field.SetInt(int64(v))
		case int64:
			field.SetInt(v)
		case float64: // JSON numbers come as float64
			field.SetInt(int64(v))
		default:
			return fmt.Errorf("cannot convert %T to int", value)
		}

	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case float32:
			field.SetFloat(float64(v))
		case float64:
			field.SetFloat(v)
		case int:
			field.SetFloat(float64(v))
		default:
			return fmt.Errorf("cannot convert %T to float", value)
		}

	case reflect.Bool:
		if b, ok := value.(bool); ok {
			field.SetBool(b)
		} else {
			return fmt.Errorf("cannot convert %T to bool", value)
		}

	case reflect.Struct:
		// For struct fields, marshal the value back to JSON and unmarshal to the struct
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal struct value: %v", err)
		}

		// Create a new instance of the struct type
		newStruct := reflect.New(field.Type())
		if err := json.Unmarshal(jsonBytes, newStruct.Interface()); err != nil {
			return fmt.Errorf("failed to unmarshal to struct: %v", err)
		}

		field.Set(newStruct.Elem())

	default:
		panic(fmt.Errorf("unhandled field type %s", field.Kind()))
	}

	return nil
}
