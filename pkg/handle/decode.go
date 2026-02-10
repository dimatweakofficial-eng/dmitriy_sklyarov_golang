package handle

import (
	"encoding/json"
	"io"
)

func Decode[T any](body io.ReadCloser) (*T, error) {
	var payload T
	//из запроса в payload по полям json 
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
