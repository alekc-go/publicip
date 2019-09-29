package publicip

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"gopkg.in/resty.v1"
)

var mirrors = []string{
	"https://api.ipify.org",
	"https://ifconfig.me",
	"https://icanhazip.com",
	"https://ipecho.net/plain",
	"https://ifconfig.co",
}

// DefaultUserAgent is the user agent which will be used for connection.
// Note: changing this value may affect the result with some providers.
var DefaultUserAgent = "curl/7.58.0"

//Debug requests
var Debug = false

// HttpClient is a placehold for alternative http client.
// If it's not set (equal nil), then new client instance is generated for
// every request
var HttpClient *resty.Client

// SetMirrors permit to override 3d party ip resolvers used in this library (if for some reason you want to
// use your own)
func SetMirrors(newUrls []string) {
	mirrors = newUrls
}

func Get() (string, error) {
	var pubIp string
	var err error
	client := resty.New()
	client.SetTimeout(time.Second)
	if Debug {
		client.SetDebug(Debug)
		client.SetHeader("User-Agent", DefaultUserAgent)
	}

	for _, url := range mirrors {
		pubIp, err = download(client, url)
		if err != nil {
			return pubIp, nil
		}
		if Debug {
			fmt.Printf("Potential ip [%v] obtained from %s is invalid\n", pubIp, url)
		}
	}
	return "", errors.New("couldn't obtain a valid ip from any source")
}

func download(cl *resty.Client, url string) (string, error) {
	resp, err := cl.R().
		SetHeader("User-Agent", DefaultUserAgent).
		Get(url)

	//check for errors and valid response
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", DownloadError{
			StatusCode: resp.StatusCode(),
			Body:       resp.Body(),
			Url:        url,
		}
	}

	pubIp := strings.TrimSpace(string(resp.Body()))
	if net.ParseIP(pubIp) == nil {
		return "", InvalidResponseError{Response: pubIp}
	}
	return pubIp, nil
}
