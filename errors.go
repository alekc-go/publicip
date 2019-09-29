package publicip

type DownloadError struct {
	StatusCode int
	Body       []byte
	Url        string
}

func (DownloadError) Error() string {
	return "error on downloading resource"
}

//InvalidResponseError is returned when a mirrir is replying an invalid response (not an ip)
type InvalidResponseError struct {
	Response string
}

func (InvalidResponseError) Error() string {
	return "invalid response"
}

type MirrorsExausted struct {
}

func (MirrorsExausted) Error() string {
	return "ALl mirrors have been used, and no valid result has been found"
}
