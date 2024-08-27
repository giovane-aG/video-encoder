package utils_test

import (
	"encoding/json"
	"testing"

	"github.com/giovane-aG/video-encoder/encoder/infrastructure/utils"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnErrorWhenJsonIsNotValid(t *testing.T) {
	var invalidJson string = "{id: 123"

	err := utils.IsJson(invalidJson)

	require.NotNil(t, err)
}

func TestShouldReturnNilIfJsonIsValid(t *testing.T) {
	st := struct {
		Id int
	}{
		Id: 1234,
	}

	validJson, err := json.Marshal(&st)

	require.Nil(t, err)

	err = utils.IsJson(string(validJson))
	require.Nil(t, err)

}
