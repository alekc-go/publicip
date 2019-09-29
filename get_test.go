package publicip

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"

	"gopkg.in/resty.v1"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	//prevent execution of default client
	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
	returnCode := m.Run()
	httpmock.DeactivateAndReset()
	os.Exit(returnCode)
}
func TestSetMirrors(t *testing.T) {
	newMirrors := []string{"www.google.com", "www.msn.com"}
	assert.NotEqual(t, newMirrors, mirrors)
	SetMirrors(newMirrors)
	assert.Equal(t, newMirrors, mirrors, "mirrors should be overriden")
}

func TestDownload(t *testing.T) {
	httpmock.Reset()
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	res, err := download(restyClient, "")
	assert.Error(t, err)
	assert.Empty(t, res)

	fixture := `error`
	responder := httpmock.NewStringResponder(503, fixture)
	testUrl := "https://api.ipify.org"
	httpmock.RegisterResponder("GET", testUrl, responder)

	res, err = download(restyClient, testUrl)
	assert.Error(t, err)
	assert.IsType(t, DownloadError{}, err)
	assert.Equal(t, err.(DownloadError).Url, testUrl)
	assert.Equal(t, err.(DownloadError).StatusCode, 503)
	assert.Equal(t, err.(DownloadError).Body, []byte(fixture))

	//invalid response
	httpmock.Reset()
	fixture = `not an ip`
	responder = httpmock.NewStringResponder(200, fixture)
	httpmock.RegisterResponder("GET", testUrl, responder)
	res, err = download(restyClient, testUrl)
	assert.Error(t, err)
	assert.IsType(t, InvalidResponseError{}, err)
	assert.Equal(t, err.(InvalidResponseError).Response, fixture)

	httpmock.Reset()
	fixture = `1.2.3.4`
	responder = httpmock.NewStringResponder(200, fixture)
	httpmock.RegisterResponder("GET", testUrl, responder)
	res, err = download(restyClient, testUrl)
	assert.NoError(t, err)
	assert.Equal(t, fixture, res)
}
