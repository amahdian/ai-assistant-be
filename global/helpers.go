package global

import (
	"encoding/json"

	"github.com/google/uuid"
)

func MergeStruct(to interface{}, from interface{}) error {
	byte1, err := json.Marshal(to)
	if err != nil {
		return err
	}
	byte2, err := json.Marshal(from)
	if err != nil {
		return err
	}
	map1 := make(map[string]interface{})
	err = json.Unmarshal(byte1, &map1)
	if err != nil {
		return err
	}
	map2 := make(map[string]interface{})
	err = json.Unmarshal(byte2, &map2)
	if err != nil {
		return err
	}
	for k, v := range map2 {
		map1[k] = v
	}
	byteDest, err := json.Marshal(map1)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteDest, to)
	return err
}

func Uuid() string {
	return uuid.New().String()
}
