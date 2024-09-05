package cloudflareeth

import "net/http"

const defaultCloudFlareBaseURL = "https://cloudflare-eth.com"

type Config struct {
	requester requester
}

type ConfigOptionResolver func(*Config)

var defaultConfigResolvers = []ConfigOptionResolver{
	WithHTTPClient(defaultRequester()),
}

func WithHTTPClient(requestClient requester) ConfigOptionResolver {
	return func(c *Config) {
		c.requester = requestClient
	}
}

func LoadDefaultConfig() Config {
	// Load default config
	var config Config

	for _, resolver := range defaultConfigResolvers {
		resolver(&config)
	}

	return config
}

// defaultRequester creates a new http client with a base url.
func defaultRequester() requester {
	c := *http.DefaultClient
	c.Timeout = defaultTimeout

	return httpClient{baseURL: defaultCloudFlareBaseURL, httpClient: c}
}
