package error

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateError(t *testing.T) {
	errorStruct := CreateError(DuplicateUserName, "/register")
	// TODO: http und localhost
	assert.Equal(t, "/register", errorStruct.Link)
	assert.Equal(t, string(DuplicateUserName), errorStruct.Text)
}
