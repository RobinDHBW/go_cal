// Matrikelnummern:
// 9495107, 4706893, 9608900

package configuration

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadFlags(t *testing.T) {
	ReadFlags()
	assert.EqualValues(t, 8080, flag.Lookup("port").Value.(flag.Getter).Get().(int))
	assert.EqualValues(t, 10, flag.Lookup("timeout").Value.(flag.Getter).Get().(int))
	assert.EqualValues(t, "./files", flag.Lookup("folder").Value.(flag.Getter).Get().(string))
	assert.EqualValues(t, "./server.crt", flag.Lookup("certPath").Value.(flag.Getter).Get().(string))
	assert.EqualValues(t, "./server.key", flag.Lookup("keyPath").Value.(flag.Getter).Get().(string))
}
