package h

import (
	"io"
	"net/http"
)

func GetString(url string) (string, error) {
	body, err := Get(url)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Get(url string) ([]byte, error) {
	var body []byte
	/* #nosec GG107 */
	res, err := http.Get(url)
	if err != nil {
		return body, err
	}

	/* #nosec G307 */
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
