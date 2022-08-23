package kbapi

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const (
// basePathKibanaSavedObject = "/api/saved_objects" // Base URL to access on Kibana data view API
)

// KibanaDataView is the Data View API object
type KibanaSavedObject struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	UpdatedAt  string                 `json:"updated_at"`
	Version    string                 `json:"version"`
	Attributes map[string]interface{} `json:"attributes"`
	References []Reference            `json:"references"`
	// Attributes KibanaSavedObjectAttrs `json:"attributes"`
}

type KibanaSavedObjectRequest struct {
	Attributes map[string]interface{} `json:"attributes"`
	References []Reference            `json:"references"`
	// initialNamespaces
}

type Reference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// type KibanaSavedObjectAttrs struct {
// 	ID   string `json:"id"`
// 	Type string `json:"type"`
// }

// KibanaDataViewGet permit to get data view
type KibanaSavedObjectV2Get func(id, objType string) (*KibanaSavedObject, error)

// KibanaDataViewCreate permit to create data view
type KibanaSavedObjectV2Create func(so *KibanaSavedObject, overwrite bool) (*KibanaSavedObject, error)

// KibanaDataViewDelete permit to delete data view
type KibanaSavedObjectV2Delete func(id, objType string) error

// KibanaDataViewUpdate permit to update data view
type KibanaSavedObjectV2Update func(dv *KibanaSavedObject) (*KibanaSavedObject, error)

// String permit to return KibanaDataView object as JSON string
func (k *KibanaSavedObject) String() string {
	json, _ := json.Marshal(k)
	return string(json)
}

// newKibanaDataViewGetFunc permit to get the kibana data view with it id
func newKibanaSavedObjectV2GetFunc(c *resty.Client) KibanaSavedObjectV2Get {
	return func(id, objType string) (*KibanaSavedObject, error) {

		if id == "" || objType == "" {
			return nil, NewAPIError(600, "You must provide both ID and object type")
		}
		log.Debug("ID: ", id)
		log.Debug("Type: ", objType)

		path := fmt.Sprintf("%s/%s/%s", basePathKibanaSavedObject, objType, id)
		log.Debug("Path: ", path)
		result := &KibanaSavedObject{}
		resp, err := c.R().SetResult(result).Get(path)
		if err != nil {
			return nil, err
		}
		log.Debug("Response: ", resp)
		if !resp.IsSuccess() {
			if resp.StatusCode() == 404 {
				return nil, nil
			}
			return nil, NewAPIError(resp.StatusCode(), resp.Status())
		}

		log.Debug("Result: ", result)

		return result, nil
	}

}

// newKibanaDataViewCreateFunc permit to create new Kibana data view
func newKibanaSavedObjectV2CreateFunc(c *resty.Client) KibanaSavedObjectV2Create {
	return func(so *KibanaSavedObject, overwrite bool) (*KibanaSavedObject, error) {

		if so == nil {
			return nil, NewAPIError(600, "You must provide KibanaSavedObject")
		}
		log.Debug("KibanaSavedObject: ", so)

		path := fmt.Sprintf("%s/%s/%s", basePathKibanaSavedObject, so.Type, so.ID)
		log.Debug("Path: ", path)

		req := newRequest(so)
		log.Debug("Request: ", req)
		resp, err := c.R().SetHeader("Content-Type", "application/json").
			SetQueryParam("overwrite", strconv.FormatBool(overwrite)).
			SetBody(req).SetResult(so).Post(path)
		if err != nil {
			return nil, err
		}

		log.Debug("Response: ", resp)
		if !resp.IsSuccess() {
			return nil, NewAPIError(resp.StatusCode(), resp.Status())
		}

		log.Debug("Result: ", so)

		return so, nil
	}

}

// newKibanaDataViewDeleteFunc permit to delete the kibana data view by id
func newKibanaSavedObjectV2Func(c *resty.Client) KibanaSavedObjectV2Delete {
	return func(id, objType string) error {

		if id == "" || objType == "" {
			return NewAPIError(600, "You must provide both ID and object type")
		}

		log.Debug("ID: ", id)
		log.Debug("Type: ", objType)

		path := fmt.Sprintf("%s/%s/%s", basePathKibanaSavedObject, objType, id)
		log.Debug("Path: ", path)
		resp, err := c.R().Delete(path)
		if err != nil {
			return err
		}
		log.Debug("Response: ", resp)
		if !resp.IsSuccess() {
			return NewAPIError(resp.StatusCode(), resp.Status())
		}

		return nil
	}

}

// newKibanaDataViewUpdateFunc permit to update the Kibana data view by id
func newKibanaSavedObjectV2UpdateFunc(c *resty.Client) KibanaSavedObjectV2Update {
	return func(so *KibanaSavedObject) (*KibanaSavedObject, error) {

		if so == nil {
			return nil, NewAPIError(600, "You must provide KibanaSavedObject")
		}
		log.Debug("KibanaSavedObject: ", so)

		path := fmt.Sprintf("%s/%s/%s", basePathKibanaSavedObject, so.Type, so.ID)
		log.Debug("Path: ", path)

		req := newRequest(so)
		resp, err := c.R().SetHeader("Content-Type", "application/json").
			SetBody(req).SetResult(so).Put(path)
		if err != nil {
			return nil, err
		}

		log.Debug("Response: ", resp)
		if !resp.IsSuccess() {
			return nil, NewAPIError(resp.StatusCode(), resp.Status())
		}

		log.Debug("Result: ", so)

		return so, nil
	}

}

func newRequest(so *KibanaSavedObject) KibanaSavedObjectRequest {
	var req KibanaSavedObjectRequest
	req.Attributes = so.Attributes
	if so.References == nil {
		req.References = make([]Reference, 0)
	} else {
		req.References = so.References
	}
	return req
}
