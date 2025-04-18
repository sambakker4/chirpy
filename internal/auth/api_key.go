package auth

import(
	"net/http"
	"strings"
	"errors"
)

func GetAPIKey(headers http.Header) (string, error) {
	key := headers.Get("Authorization")
	apiKey, ok := strings.CutPrefix(key, "ApiKey ")
	if !ok {
		return "", errors.New("Authorization header should be in the format: ApiKey <apikey>")
	}

	return apiKey, nil
}
