package utils_test

import (
	"testing"

	"github.com/giovane-aG/video-encoder/encoder/infrastructure/utils"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnErrorWhenJsonIsNotValid(t *testing.T) {
	var invalidJson string = "{id: 123"

	err := utils.IsJson(invalidJson)

	require.NotNil(t, err)
}
