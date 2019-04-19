package util

import (
	"time"

	"github.com/imroc/req"
	jsoniter "github.com/json-iterator/go"
)

// JSONStructToMap convert struct to map
func JSONStructToMap(obj interface{}) map[string]interface{} {
	jsonBytes, _ := jsoniter.Marshal(obj)
	var data map[string]interface{}
	jsoniter.Unmarshal(jsonBytes, &data)
	return data
}

// BindGetJSONData bind the json data of method GET
// body must be a point
func BindGetJSONData(url string, param req.Param, body interface{}) error {
	r, err := req.Get(url, param)
	if err != nil {
		return err
	}
	err = r.ToJSON(body)
	if err != nil {
		return err
	}
	return nil
}

func GetNowTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}
