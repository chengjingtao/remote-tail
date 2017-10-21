package alauda

import (
	"fmt"
	"reflect"
)

type jsonValue interface{}

type jsonValueE struct {
	jsonValue
	err error
}

//NewJSONObjE new json object with errorable
func NewJSONObjE(jsonObj jsonValue) jsonValueE {
	return jsonValueE{
		jsonValue: jsonObj,
		err:       nil,
	}
}

func (jsonE jsonValueE) Value() interface{} {
	return jsonE.jsonValue
}

func (jsonE jsonValueE) Error() error {
	if jsonE.err == nil {
		return nil
	}
	return jsonE.err
}

func (jsonE jsonValueE) HasKey(key string) bool {
	if jsonE.err != nil {
		return false
	}
	v, ok := jsonE.jsonValue.(map[string]interface{})
	if !ok {
		return false
	}

	_, ok = v[key]
	return ok
}

func (jsonE jsonValueE) Get(key string) jsonValueE {
	if jsonE.err != nil {
		return jsonE
	}

	v, ok := jsonE.jsonValue.(map[string]interface{})
	if !ok {
		jsonE.err = fmt.Errorf("The key [%s] is not exist in %#v", key, jsonE.jsonValue)
		return jsonE
	}

	vv, ok := v[key]
	if ok {
		return NewJSONObjE(vv)
	}

	jsonE.err = fmt.Errorf("The key [%s] is not exist in %#v", key, jsonE.jsonValue)
	return jsonE
}

func (jsonE jsonValueE) Index(index int) jsonValueE {
	if jsonE.err != nil {
		return jsonE
	}
	arr, ok := jsonE.jsonValue.([]interface{})
	if !ok {
		jsonE.err = fmt.Errorf("%#v is not an array", jsonE.jsonValue)
		return jsonE
	}
	if index >= len(arr) {
		jsonE.err = fmt.Errorf("Out of range for data: %#v  when index=%d", jsonE.jsonValue, index)
		return jsonE
	}
	return NewJSONObjE(arr[index])
}

func (jsonE jsonValueE) AsInt() (int, error) {
	if jsonE.err != nil {
		return -1, jsonE.err
	}
	v, ok := jsonE.jsonValue.(float64)
	if ok {
		return int(v), nil
	}
	return -1, fmt.Errorf("%#v is not Int, But %s", jsonE.Value(), reflect.TypeOf(jsonE.jsonValue))
}

func (jsonE jsonValueE) AsString() (string, error) {
	if jsonE.err != nil {
		return "", jsonE.err
	}
	v, ok := jsonE.jsonValue.(string)
	if ok {
		return v, nil
	}
	return "", fmt.Errorf("%#v is not string", jsonE.Value())
}

func (jsonE jsonValueE) AsBool() (bool, error) {
	if jsonE.err != nil {
		return false, jsonE.err
	}
	v, ok := jsonE.jsonValue.(bool)
	if ok {
		return v, nil
	}
	return false, fmt.Errorf("%#v is not bool", jsonE.Value())
}

func (jsonE jsonValueE) AsArray() ([]interface{}, error) {
	if jsonE.err != nil {
		return []interface{}{}, jsonE.err
	}
	v, ok := jsonE.jsonValue.([]interface{})
	if ok {
		return v, nil
	}
	return []interface{}{}, fmt.Errorf("%#v is not []interface{}", jsonE.Value())
}
