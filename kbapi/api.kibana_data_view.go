package kbapi

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const (
	basePathKibanaDataView = "/api/data_views" // Base URL to access on Kibana data view API
)

// KibanaDataView is the Data View API object
type KibanaDataView struct {
	ID            string `json:"id,omitempty"`
	Title         string `json:"title"`
	TimeFieldName string `json:"timeFieldName"`
}

type KibanaDataViewRequest struct {
	Override      bool            `json:"override,omitempty"`
	RefreshFields bool            `json:"refresh_fields,omitempty"`
	DataView      *KibanaDataView `json:"data_view"`
}

type KibanaDataViewResponse struct {
	DataView *KibanaDataView `json:"data_view"`
}

// KibanaDataViewGet permit to get data view
type KibanaDataViewGet func(id string) (*KibanaDataView, error)

// KibanaDataViewCreate permit to create data view
type KibanaDataViewCreate func(dv *KibanaDataView, override, refreshFields bool) (*KibanaDataView, error)

// KibanaDataViewDelete permit to delete data view
type KibanaDataViewDelete func(id string) error

// KibanaDataViewUpdate permit to update data view
type KibanaDataViewUpdate func(dv *KibanaDataView, refreshFields bool) (*KibanaDataView, error)

// String permit to return KibanaDataView object as JSON string
func (k *KibanaDataView) String() string {
	json, _ := json.Marshal(k)
	return string(json)
}

// newKibanaDataViewGetFunc permit to get the kibana data view with it id
func newKibanaDataViewGetFunc(c *resty.Client) KibanaDataViewGet {
	return func(id string) (*KibanaDataView, error) {

		if id == "" {
			return nil, NewAPIError(600, "You must provide kibana data view ID")
		}
		log.Debug("KibanaDataViewGet ID: ", id)

		path := fmt.Sprintf("%s/data_view/%s", basePathKibanaDataView, id)
		log.Debug("Path: ", path)
		result := &KibanaDataViewResponse{}
		resp, err := c.R().SetResult(result).Get(path)
		if err != nil {
			return nil, err
		}
		log.Debug("KibanaDataViewGet Response: ", resp)
		if !resp.IsSuccess() {
			if resp.StatusCode() == 404 {
				return nil, nil
			}
			return nil, NewAPIError(resp.StatusCode(), resp.Status())
		}

		log.Debug("KibanaDataViewGet Result: ", result.DataView)

		return result.DataView, nil
	}

}

// newKibanaDataViewCreateFunc permit to create new Kibana data view
func newKibanaDataViewCreateFunc(c *resty.Client) KibanaDataViewCreate {
	return func(dv *KibanaDataView, override, refreshFields bool) (*KibanaDataView, error) {

		if dv == nil {
			return nil, NewAPIError(600, "You must provide KibanaDataView object")
		}
		log.Debug("KibanaDataViewCreate: ", dv)

		var body KibanaDataViewRequest
		body.DataView = dv
		body.Override = override
		body.RefreshFields = refreshFields

		result := &KibanaDataViewResponse{}

		path := fmt.Sprintf("%s/data_view", basePathKibanaDataView)
		resp, err := c.R().SetHeader("Content-Type", "application/json").
			SetBody(body).SetResult(result).Post(path)
		if err != nil {
			return nil, err
		}

		log.Debug("KibanaDataViewCreate Response: ", resp)
		if !resp.IsSuccess() {
			return nil, NewAPIError(resp.StatusCode(), resp.Status())
		}

		log.Debug("KibanaDataViewCreate Result: ", result.DataView)

		return result.DataView, nil
	}

}

// newKibanaDataViewDeleteFunc permit to delete the kibana data view by id
func newKibanaDataViewDeleteFunc(c *resty.Client) KibanaDataViewDelete {
	return func(id string) error {

		if id == "" {
			return NewAPIError(600, "You must provide kibana data view ID")
		}

		log.Debug("KibanaDataViewDelete ID: ", id)

		path := fmt.Sprintf("%s/data_view/%s", basePathKibanaDataView, id)
		resp, err := c.R().Delete(path)
		if err != nil {
			return err
		}
		log.Debug("KibanaDataViewDelete Response: ", resp)
		if !resp.IsSuccess() {
			return NewAPIError(resp.StatusCode(), resp.Status())
		}

		return nil
	}

}

// newKibanaDataViewUpdateFunc permit to update the Kibana data view by id
func newKibanaDataViewUpdateFunc(c *resty.Client) KibanaDataViewUpdate {
	return func(dv *KibanaDataView, refreshFields bool) (*KibanaDataView, error) {

		if dv == nil {
			return nil, NewAPIError(600, "You must provide KibanaDataViewRequest object")
		}
		log.Debug("KibanaDataViewUpdate: ", dv)

		path := fmt.Sprintf("%s/data_view/%s", basePathKibanaDataView, dv.ID)

		var body KibanaDataViewRequest
		body.DataView = dv
		body.RefreshFields = refreshFields
		body.DataView.ID = ""

		result := &KibanaDataViewResponse{}

		resp, err := c.R().SetBody(body).SetResult(result).Post(path)
		if err != nil {
			return nil, err
		}

		log.Debug("KibanaDataViewUpdate Response: ", resp)
		if !resp.IsSuccess() {
			return nil, NewAPIError(resp.StatusCode(), resp.Status())
		}

		log.Debug("KibanaDataViewUpdate Result: ", result.DataView)

		return result.DataView, nil
	}

}
