package comb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func Post(theUrl string, data interface{}, result interface{}) error {
	content, _ := json.Marshal(data)
	resp, err := http.Post(theUrl, "application/json;charset=utf-8", strings.NewReader(string(content)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, result)
	return nil
}
