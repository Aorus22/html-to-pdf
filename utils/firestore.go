package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	projectID    = "cbt-01-a21ba"
	apiKey       = "AIzaSyAdZlIdf14WAPVIHoMZaIVQVrld5a6tpa0"
	firebaseHost = "https://firestore.googleapis.com/v1"
)

func GetDocument(path string) (map[string]any, error) {
	url := fmt.Sprintf("%s/projects/%s/databases/(default)/documents/%s?key=%s", firebaseHost, projectID, path, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("firestore error: %s", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func ParseFirestoreFields(fields map[string]any) map[string]any {
	result := make(map[string]any)

	for key, val := range fields {
		if val == nil {
			continue
		}
		v, ok := val.(map[string]any)
		if !ok {
			continue
		}

		switch {
		case v["stringValue"] != nil:
			result[key] = v["stringValue"].(string)

		case v["integerValue"] != nil:
			ival := 0
			fmt.Sscanf(v["integerValue"].(string), "%d", &ival)
			result[key] = ival

		case v["mapValue"] != nil:
			if mv, ok := v["mapValue"].(map[string]any); ok {
				if sf, ok := mv["fields"].(map[string]any); ok {
					result[key] = ParseFirestoreFields(sf)
				}
			}

		case v["arrayValue"] != nil:
			arr := []any{}
			rawArr, ok := v["arrayValue"].(map[string]any)["values"].([]any)
			if ok {
				for _, item := range rawArr {
					itemVal, ok := item.(map[string]any)
					if !ok {
						continue
					}

					if itemVal["mapValue"] != nil {
						if mv, ok := itemVal["mapValue"].(map[string]any); ok {
							if sf, ok := mv["fields"].(map[string]any); ok {
								arr = append(arr, ParseFirestoreFields(sf))
							}
						}
					} else {
						arr = append(arr, itemVal)
					}
				}
			}
			result[key] = arr

		default:
			result[key] = v
		}
	}

	return result
}
