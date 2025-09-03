package sanitize

import (
	"fmt"
	"net/url"
)

// SanitizeURLParams receives a URL as input and returns a sanitized version of it,
// where the parameter values have been hardcoded to prevent malicious injections.
// Avoid HTTP parameter pollution [CWE:235]
func SanitizeURLParams(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL to sanitize it: %w", err)
	}

	// Get URL query parameters
	queryParams := parsedURL.Query()

	// Encode parameter values
	for param, values := range queryParams {
		for _, value := range values {
			// The parameter value is encoded here to ensure that special characters are handled securely.
			// special characters are handled safely.
			// `url.QueryEscape` encodes the value so that it is safe to use in a URL.
			// Then, the original value is replaced with the encoded value using the `Set` method.
			queryParams.Set(param, url.QueryEscape(value))
		}
	}

	// Establecer la cadena de consulta codificada en la URL
	parsedURL.RawQuery = queryParams.Encode()

	return parsedURL.String(), nil
}
