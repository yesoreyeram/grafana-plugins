package anyframer_test

import (
	"errors"
	"testing"
)

func TestRootSelector(t *testing.T) {
	t.Run("invalid selector", testToFrame(testInputs{
		input:   `{ "users" : [{"name":"foo","salary": 123, "self_employed":false},{"name":"bar","salary": 456.789, "self_employed":true}] }`,
		framer:  Framer{RootSelector: "userslist"},
		wantErr: errors.New("error applying root selector"),
	}))
	t.Run("array selector with value", testToFrame(testInputs{
		input:  `{ "users" : [{"name":"foo","salary": 123, "self_employed":false},{"name":"bar","salary": 456.789, "self_employed":true}] }`,
		framer: Framer{RootSelector: "users"},
	}))
	// FIXME: Empty array should not throw error
	t.Run("array selector without value", testToFrame(testInputs{
		input:   `{ "users" : [] }`,
		framer:  Framer{RootSelector: "users"},
		wantErr: errors.New("error applying root selector"),
	}))
}
