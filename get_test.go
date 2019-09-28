package publicip

import (
	"testing"

	"gopkg.in/resty.v1"

	"github.com/stretchr/testify/assert"
)

func TestSetMirrors(t *testing.T) {
	newMirrors := []string{"www.google.com", "www.msn.com"}
	assert.NotEqual(t, newMirrors, mirrors)
	SetMirrors(newMirrors)
	assert.Equal(t, newMirrors, mirrors, "mirrors should be overriden")
}
func TestDownload(t *testing.T) {
	client := resty.New()
	client.SetHeader("User-Agent", DefaultUserAgent)
	ip := ""
	for _, url := range mirrors {
		pubIP, err := download(client, url)
		if ip == "" {
			ip = pubIP
		}
		assert.NoError(t, err, "error during download of %s", err)
		assert.Equal(t, ip, pubIP, "ip should be the same")
	}
}
