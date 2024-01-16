package common

import (
	"net/url"
	"regexp"
)

// Neutralization for unwanted command injections in domain string
func NeutralizeDomain(input string) string {
	pattern := "^(?:https?://)?(?:www.)?([^:/\n&=?¿\"!| %]+)"
	regex := regexp.MustCompile(pattern)
	return regex.FindString(input)
}

// Neutralization HTTP parameter pollution. CWE:235
func NeutralizeParamPollution(queryParams map[string][]string) string {
	params := url.Values{}
	for key, values := range queryParams {
		for _, value := range values {
			if !params.Has(key) {
				params.Add(key, value)
			}
		}
	}
	return params.Encode()
}
