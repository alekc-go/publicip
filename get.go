package publicip

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"gopkg.in/resty.v1"
)

var urls = [...]string{"api.ipify.org", "ifconfig.me", "icanhazip.com", "ipecho.net/plain", "ifconfig.co"}

// DefaultUserAgent is the user agent which will be used for connection.
// Note: changing this value may affect the result with some providers.
var DefaultUserAgent = "curl/7.58.0"

//Scheme to be used in request (http/https)
var Scheme = "https"

//Debug requests
var Debug = false

func Get() (string, error) {
	var pubIp string
	var err error
	client := resty.New()
	client.SetTimeout(time.Second)
	if Debug {
		client.SetDebug(Debug)
		client.SetHeader("User-Agent", DefaultUserAgent)
	}

	for _, url := range urls {
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
	resp, err := resty.R().
		SetHeader("User-Agent", DefaultUserAgent).
		Get(Scheme + "://" + url)
	if err != nil {
		return "", err
	}
	pubIp := strings.TrimSpace(string(resp.Body()))
	if net.ParseIP(pubIp) == nil {
		return "", fmt.Errorf("invalid ip [%s]", pubIp)
	}
	return pubIp, nil
}
