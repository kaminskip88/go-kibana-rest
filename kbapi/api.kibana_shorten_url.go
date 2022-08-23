package kbapi

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const (
	basePathKibanaShortenURL = "/api/short_url" // Base URL to access on Kibana shorten URL
)

// ShortenURL is the shorten URL object
type ShortenURL struct {
	URL string `json:"url"`
}

type shortURLReq struct {
	LocatorId string        `json:"locatorId"`
	Params    shortURLParam `json:"params"`
}

type shortURLParam struct {
	URL string `json:"url"`
}

// ShortenURLResponse is the shorten URL object response
type ShortenURLResponse struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

// KibanaShortenURLCreate permit to create new shorten URL
type KibanaShortenURLCreate func(shortenURL *ShortenURL) (*ShortenURLResponse, error)

// String permit to return ShortenURL object as JSON string
func (o *ShortenURL) String() string {
	json, _ := json.Marshal(o)
	return string(json)
}

// String permit to return ShortenURLResponse object as JSON string
func (o *ShortenURLResponse) String() string {
	json, _ := json.Marshal(o)
	return string(json)
}

// newKibanaShortenURLCreateFunc permit to create new shorten URL
func newKibanaShortenURLCreateFunc(c *resty.Client) KibanaShortenURLCreate {
	return func(shortenURL *ShortenURL) (*ShortenURLResponse, error) {

		if shortenURL == nil {
			return nil, NewAPIError(600, "You must provide shorten URL object")
		}
		log.Debug("Shorten URL: ", shortenURL)

		var req shortURLReq
		req.LocatorId = "LEGACY_SHORT_URL_LOCATOR"
		req.Params.URL = shortenURL.URL

		jsonData, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		log.Debugf("Shorten URL payload: %s", jsonData)

		resp, err := c.R().SetBody(jsonData).Post(basePathKibanaShortenURL)
		if err != nil {
			return nil, err
		}

		log.Debug("Response: ", resp)
		if resp.StatusCode() >= 300 {
			return nil, NewAPIError(resp.StatusCode(), resp.Status())
		}

		shortenURLResponse := &ShortenURLResponse{}
		err = json.Unmarshal(resp.Body(), shortenURLResponse)
		if err != nil {
			return nil, err
		}
		log.Debug("ShortenURLResponse: ", shortenURLResponse)

		return shortenURLResponse, nil
	}
}
