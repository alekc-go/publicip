package publicip

type DownloadError struct {
	StatusCode int
	Body       []byte
	Url        string
}

func (DownloadError) Error() string {
	return "error on downloading resource"
}

type InvalidResponseError struct {
	Response string
}

func (InvalidResponseError) Error() string {
	return "invalid response"
}
