package parser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func ParseModelFromJSON(filename string) (map[string]map[string]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	keys := reflect.ValueOf(m).MapKeys()
	if len(keys) != 1 {
		return nil, errors.New("Expected only one model definition.")
	}

	modelName := keys[0].String()

	modelMap := ParseModel(m, modelName)
	return modelMap, nil
}

func ParseModel(m map[string]interface{}, modelName string) map[string]map[string]string {
	return convertToStruct(m, make(map[string]map[string]string), modelName)
}

func convertToStruct(m map[string]interface{}, modelMap map[string]map[string]string, structName string) map[string]map[string]string {
	fm := make(map[string]string)
	mapName := make([]string, 0)
	mapValue := make([]interface{}, 0)
	for k, v := range m {
		typeName := reflect.TypeOf(v).String()
		if reflect.TypeOf(v).String() == "map[string]interface {}" {
			typeName = strings.Title(k) + "0"
			mapName = append(mapName, typeName)
			mapValue = append(mapValue, v)
		}
		fm[strings.Title(k)] = typeName
	}
	modelMap[structName] = fm
	if len(mapName) > 0 {
		for i, n := range mapName {
			modelMap = convertToStruct(mapValue[i].(map[string]interface{}), modelMap, n)
		}
	}
	return modelMap
}
