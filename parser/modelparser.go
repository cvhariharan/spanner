package parser

import (
	"encoding/json"
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

	modelMap := ParseModel(m)
	return modelMap, nil
}

func ParseModel(m map[string]interface{}) map[string]map[string]string {
	return convertToStruct(m, make(map[string]map[string]string), "User")
}

func convertToStruct(m map[string]interface{}, modelMap map[string]map[string]string, structName string) map[string]map[string]string {
	fm := make(map[string]string)
	mapName := make([]string, 0)
	mapValue := make([]interface{}, 0)
	for k, v := range m {
		typeName := reflect.TypeOf(v).String()
		if reflect.TypeOf(v).String() == "map[string]interface {}" {
			typeName = strings.ToLower(k)
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
