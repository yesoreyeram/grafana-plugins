package anyframer

import (
	"errors"

	jsonata "github.com/blues/jsonata-go"
)

func applySelector(input any, selector string) (output any, err error) {
	if input == nil {
		return nil, errors.New("invalid/empty data")
	}
	if selector == "" {
		return input, nil
	}
	e := jsonata.MustCompile(selector)
	res, err := e.Eval(input)
	if err != nil {
		return nil, errors.New("error applying root selector")
	}
	return res, nil
}
