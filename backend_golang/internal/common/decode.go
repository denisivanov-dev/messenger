package common

import "encoding/json"

func DecodePayload(payload any, out any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}
