// Matrikelnummern:
// 9495107, 4706893, 9608900

package error

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateError(t *testing.T) {
	errorStruct := CreateError(DuplicateUserName, "/register")
	assert.Equal(t, "/register", errorStruct.Link)
	assert.Equal(t, string(DuplicateUserName), errorStruct.Text)
}
