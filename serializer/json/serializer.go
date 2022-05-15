package json

import (
	"encoding/json"

	"github.com/AntonyIS/go-loco/app"
	"github.com/pkg/errors"
)

type Locomotive struct{}

func (*Locomotive) Decode(input []byte) (*app.Locomotive, error) {
	Locomotive := &app.Locomotive{}

	if err := json.Unmarshal(input, Locomotive); err != nil {
		return nil, errors.Wrap(err, "serlializer.Locomotive.Decode")
	}
	return Locomotive, nil
}

func (*Locomotive) Encode(input *app.Locomotive) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serlializer.Locomotive.Decode")
	}
	return rawMsg, nil
}
