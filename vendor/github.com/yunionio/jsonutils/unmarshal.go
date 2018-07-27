package jsonutils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/yunionio/log"
	"github.com/yunionio/pkg/gotypes"
	"github.com/yunionio/pkg/util/reflectutils"
	"github.com/yunionio/pkg/util/timeutils"
	"github.com/yunionio/pkg/utils"
)

func (this *JSONValue) Unmarshal(obj interface{}, keys ...string) error {
	return fmt.Errorf("unsupported operation Unmarshall")
}

func (this *JSONArray) Unmarshal(obj interface{}, keys ...string) error {
	return jsonUnmarshal(this, obj, keys)
}

func (this *JSONDict) Unmarshal(obj interface{}, keys ...string) error {
	return jsonUnmarshal(this, obj, keys)
}

func jsonUnmarshal(jo JSONObject, o interface{}, keys []string) error {
	if len(keys) > 0 {
		var err error = nil
		jo, err = jo.Get(keys...)
		if err != nil {
			return err
		}
	}
	value := reflect.Indirect(reflect.ValueOf(o))
	return jo.unmarshalValue(value)
}

func (this *JSONValue) unmarshalValue(val reflect.Value) error {
	// return fmt.Errorf("JSONValue: type mismatch")
	// null value
	if val.CanSet() {
		zeroVal := reflect.New(val.Type()).Elem()
		val.Set(zeroVal)
	}
	return nil
}

func (this *JSONInt) unmarshalValue(val reflect.Value) error {
	switch val.Type() {
	case JSONIntType:
		json := val.Interface().(JSONInt)
		json.data = this.data
		return nil
	case JSONIntPtrType, JSONObjectType:
		val.Set(reflect.ValueOf(this))
		return nil
	case JSONStringType:
		json := val.Interface().(JSONString)
		json.data = fmt.Sprintf("%d", this.data)
		return nil
	case JSONStringPtrType:
		json := val.Interface().(*JSONString)
		data := fmt.Sprintf("%d", this.data)
		if json == nil {
			json = NewString(data)
			val.Set(reflect.ValueOf(json))
		} else {
			json.data = data
		}
		return nil
	case JSONBoolType, JSONFloatType, JSONArrayType, JSONDictType, JSONBoolPtrType, JSONFloatPtrType, JSONArrayPtrType, JSONDictPtrType:
		return fmt.Errorf("JSONInt type mismatch %s", val.Type())
	}
	switch val.Kind() {
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8,
		reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		val.SetInt(this.data)
	case reflect.Float32, reflect.Float64:
		val.SetFloat(float64(this.data))
	case reflect.Bool:
		if this.data == 0 {
			val.SetBool(false)
		} else {
			val.SetBool(true)
		}
	case reflect.String:
		val.SetString(fmt.Sprintf("%d", this.data))
	default:
		return fmt.Errorf("JSONInt type mismatch: %s", val.Type())
	}
	return nil
}

func (this *JSONBool) unmarshalValue(val reflect.Value) error {
	switch val.Type() {
	case JSONBoolType:
		json := val.Interface().(JSONBool)
		json.data = this.data
		return nil
	case JSONBoolPtrType, JSONObjectType:
		val.Set(reflect.ValueOf(this))
		return nil
	case JSONStringType:
		json := val.Interface().(JSONString)
		json.data = strconv.FormatBool(this.data)
		return nil
	case JSONStringPtrType:
		json := val.Interface().(*JSONString)
		data := strconv.FormatBool(this.data)
		if json == nil {
			json = NewString(data)
			val.Set(reflect.ValueOf(json))
		} else {
			json.data = data
		}
		return nil
	case JSONIntType, JSONFloatType, JSONArrayType, JSONDictType, JSONIntPtrType, JSONFloatPtrType, JSONArrayPtrType, JSONDictPtrType:
		return fmt.Errorf("JSONBool type mismatch %s", val.Type())
	}
	switch val.Kind() {
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8,
		reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		if this.data {
			val.SetInt(1)
		} else {
			val.SetInt(0)
		}
	case reflect.Float32, reflect.Float64:
		if this.data {
			val.SetFloat(1.0)
		} else {
			val.SetFloat(0.0)
		}
	case reflect.Bool:
		val.SetBool(this.data)
	case reflect.String:
		if this.data {
			val.SetString("true")
		} else {
			val.SetString("false")
		}
	default:
		return fmt.Errorf("JSONBool type mismatch: %s", val.Type())
	}
	return nil
}

func (this *JSONFloat) unmarshalValue(val reflect.Value) error {
	switch val.Type() {
	case JSONFloatType:
		json := val.Interface().(JSONFloat)
		json.data = this.data
		return nil
	case JSONFloatPtrType, JSONObjectType:
		val.Set(reflect.ValueOf(this))
		return nil
	case JSONStringType:
		json := val.Interface().(JSONString)
		json.data = fmt.Sprintf("%f", this.data)
		return nil
	case JSONStringPtrType:
		json := val.Interface().(*JSONString)
		data := fmt.Sprintf("%f", this.data)
		if json == nil {
			json = NewString(data)
			val.Set(reflect.ValueOf(json))
		} else {
			json.data = data
		}
		return nil
	case JSONIntType:
		json := val.Interface().(JSONInt)
		json.data = int64(this.data)
		return nil
	case JSONIntPtrType:
		json := val.Interface().(*JSONInt)
		if json == nil {
			json = NewInt(int64(this.data))
			val.Set(reflect.ValueOf(json))
		} else {
			json.data = int64(this.data)
		}
		return nil
	case JSONBoolType, JSONArrayType, JSONDictType, JSONBoolPtrType, JSONArrayPtrType, JSONDictPtrType:
		return fmt.Errorf("JSONFloat type mismatch %s", val.Type())
	}
	switch val.Kind() {
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8,
		reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		val.SetInt(int64(this.data))
	case reflect.Float32, reflect.Float64:
		val.SetFloat(this.data)
	case reflect.Bool:
		if this.data == 0 {
			val.SetBool(false)
		} else {
			val.SetBool(true)
		}
	case reflect.String:
		val.SetString(fmt.Sprintf("%f", this.data))
	default:
		return fmt.Errorf("JSONFloat type mismatch: %s", val.Type())
	}
	return nil
}

func (this *JSONString) unmarshalValue(val reflect.Value) error {
	switch val.Type() {
	case JSONStringType:
		json := val.Interface().(JSONString)
		json.data = this.data
		return nil
	case JSONStringPtrType, JSONObjectType:
		val.Set(reflect.ValueOf(this))
		return nil
	case gotypes.TimeType:
		var tm time.Time
		var err error
		if len(this.data) > 0 {
			tm, err = timeutils.ParseTimeStr(this.data)
			if err != nil {
				return err
			}
		} else {
			tm = time.Time{}
		}
		val.Set(reflect.ValueOf(tm))
		return nil
	case JSONBoolType, JSONIntType, JSONFloatType, JSONArrayType, JSONDictType,
		JSONBoolPtrType, JSONIntPtrType, JSONFloatPtrType, JSONArrayPtrType, JSONDictPtrType:
		return fmt.Errorf("JSONString type mismatch %s", val.Type())
	}
	switch val.Kind() {
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8,
		reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		intVal, err := strconv.ParseInt(this.data, 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(intVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(this.data, 64)
		if err != nil {
			return err
		}
		val.SetFloat(floatVal)
	case reflect.Bool:
		val.SetBool(utils.ToBool(this.data))
	case reflect.String:
		val.SetString(this.data)
	default:
		return fmt.Errorf("JSONString type mismatch: %s", val.Type())
	}
	return nil
}

func (this *JSONArray) unmarshalValue(val reflect.Value) error {
	switch val.Type() {
	case JSONArrayType:
		array := val.Interface().(JSONArray)
		if this.data != nil {
			array.Add(this.data...)
		}
		return nil
	case JSONArrayPtrType, JSONObjectType:
		val.Set(reflect.ValueOf(this))
		return nil
	case JSONDictType, JSONIntType, JSONStringType, JSONBoolType, JSONFloatType,
		JSONDictPtrType, JSONIntPtrType, JSONStringPtrType, JSONBoolPtrType, JSONFloatPtrType:
		return fmt.Errorf("JSONArray type mismatch %s", val.Type())
	}
	switch val.Kind() {
	case reflect.String:
		val.SetString(this.String())
		return nil
	case reflect.Slice, reflect.Array:
		for _, json := range this.data {
			newEle := reflect.New(val.Type().Elem()).Elem()
			err := json.unmarshalValue(newEle)
			if err != nil {
				return err
			}
			newVal := reflect.Append(val, newEle)
			val.Set(newVal)
		}
	default:
		return fmt.Errorf("JSONArray type mismatch: %s", val.Type())
	}
	return nil
}

func (this *JSONDict) unmarshalValue(val reflect.Value) error {
	switch val.Type() {
	case JSONDictType:
		dict := val.Interface().(JSONDict)
		dict.Update(this)
		return nil
	case JSONDictPtrType, JSONObjectType:
		val.Set(reflect.ValueOf(this))
		return nil
	case JSONArrayType, JSONIntType, JSONBoolType, JSONFloatType, JSONStringType,
		JSONArrayPtrType, JSONIntPtrType, JSONBoolPtrType, JSONFloatPtrType, JSONStringPtrType:
		return fmt.Errorf("JSONDict type mismatch %s", val.Type())
	}
	switch val.Kind() {
	case reflect.String:
		val.SetString(this.String())
		return nil
	case reflect.Map:
		return this.unmarshalMap(val)
	case reflect.Struct:
		return this.unmarshalStruct(val)
	default:
		return fmt.Errorf("JSONDict type mismatch: %s", val.Type())
	}
	return nil
}

func (this *JSONDict) unmarshalMap(val reflect.Value) error {
	if val.IsNil() {
		mapVal := reflect.MakeMap(val.Type())
		val.Set(mapVal)
	}
	valType := val.Type()
	keyType := valType.Key()
	if keyType.Kind() != reflect.String {
		return fmt.Errorf("map key must be string")
	}
	for k, v := range this.data {
		keyVal := reflect.ValueOf(k)
		valVal := reflect.New(valType.Elem()).Elem()

		err := v.unmarshalValue(valVal)
		if err != nil {
			log.Debugf("unmarshalMap field %s error %s", k, err)
			return err
		}
		val.SetMapIndex(keyVal, valVal)
	}
	return nil
}

func (this *JSONDict) unmarshalStruct(val reflect.Value) error {
	fieldValues := reflectutils.FetchStructFieldNameValues(val)
	for k, v := range this.data {
		fieldValue, ok := fieldValues[k] // first try original key
		if !ok {                         // try kebab
			k = utils.CamelSplit(k, "_")
			fieldValue, ok = fieldValues[k]
		}
		if ok {
			err := v.unmarshalValue(fieldValue)
			if err != nil {
				log.Debugf("unmarshalStruct field %s error %s", k, err)
				return err
			}
		}
	}
	return nil
}
