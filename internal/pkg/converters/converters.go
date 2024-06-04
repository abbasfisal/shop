package converters

import "encoding/json"

func MapToString(data map[string]string) string {
	str, _ := json.Marshal(data)

	return string(str)
}

func StringToMap(data string) map[string]string {
	var out map[string]string
	_ = json.Unmarshal([]byte(data), &out)
	return out
}

func UrlValuesToString(data map[string][]string) string {
	str, _ := json.Marshal(data)
	return string(str)
}

func StringToUrlValues(data string) map[string][]string {
	var out map[string][]string
	_ = json.Unmarshal([]byte(data), &out)
	return out
}
