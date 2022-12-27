package error

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateError(t *testing.T) {
	errorStruct := CreateError(Authentification, "/")
	// TODO: http und localhost
	assert.Equal(t, "http://localhost:8080/", errorStruct.Link)
	assert.Equal(t, string(Authentification), errorStruct.Text)
}
