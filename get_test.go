package publicip

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"

	"gopkg.in/resty.v1"
)

func TestMain(m *testing.M) {
	//prevent execution of default client
	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
	returnCode := m.Run()
	httpmock.DeactivateAndReset()
	os.Exit(returnCode)
}
func TestSetMirrors(t *testing.T) {
	newMirrors := []string{
		"https://www.google.com",
		"https://www.msn.com",
	}
	require.NotEqual(t, newMirrors, mirrors)
	SetMirrors(newMirrors)
	require.Equal(t, newMirrors, mirrors, "mirrors should be overriden")
}

func TestGet(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	HttpClient = restyClient

	fixture := `1.2.3.4`
	responder := httpmock.NewStringResponder(503, fixture)

	//lets assume that all mirrors have failed so far.
	for _, url := range mirrors {
		httpmock.RegisterResponder("GET", url, responder)
	}
	res, err := Get()
	require.Equal(t, "", res)
	require.Error(t, err)
	require.IsType(t, MirrorsExausted{}, err)

	//ok, lets return a valid result
	httpmock.Reset()
	responder = httpmock.NewStringResponder(200, fixture)
	for _, url := range mirrors {
		httpmock.RegisterResponder("GET", url, responder)
	}
	res, err = Get()
	require.NoError(t, err)
	require.Equal(t, fixture, res)
}

func TestDownload(t *testing.T) {
	restyClient := resty.New()
	httpmock.ActivateNonDefault(restyClient.GetClient())
	res, err := download(restyClient, "")
	require.Error(t, err)
	require.Empty(t, res)

	fixture := `error`
	responder := httpmock.NewStringResponder(503, fixture)
	testUrl := "https://api.ipify.org"
	httpmock.RegisterResponder("GET", testUrl, responder)

	res, err = download(restyClient, testUrl)
	require.Error(t, err)
	require.IsType(t, DownloadError{}, err)
	require.Equal(t, err.(DownloadError).Url, testUrl)
	require.Equal(t, err.(DownloadError).StatusCode, 503)
	require.Equal(t, err.(DownloadError).Body, []byte(fixture))

	//invalid response
	httpmock.Reset()
	fixture = `not an ip`
	responder = httpmock.NewStringResponder(200, fixture)
	httpmock.RegisterResponder("GET", testUrl, responder)
	res, err = download(restyClient, testUrl)
	require.Error(t, err)
	require.IsType(t, InvalidResponseError{}, err)
	require.Equal(t, err.(InvalidResponseError).Response, fixture)

	httpmock.Reset()
	fixture = `1.2.3.4`
	responder = httpmock.NewStringResponder(200, fixture)
	httpmock.RegisterResponder("GET", testUrl, responder)
	res, err = download(restyClient, testUrl)
	require.NoError(t, err)
	require.Equal(t, fixture, res)
}
