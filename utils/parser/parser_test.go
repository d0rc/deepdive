package parser

import (
	"encoding/json"
	"testing"
)

func TestTryParseData(t *testing.T) {
	context := "[]"

	result := make([]string, 0)
	_, err := TryParseData(func(subContext string) ([]string, error) {
		err := json.Unmarshal([]byte(subContext), &result)
		return result, err
	}, context)

	if err != nil {
		panic(err)
	}
}
