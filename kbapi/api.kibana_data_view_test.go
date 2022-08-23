package kbapi

import (
	"github.com/stretchr/testify/assert"
)

func (s *KBAPITestSuite) TestKibanaDataView() {

	// Create new data view
	dv := &KibanaDataView{
		Title:         "test",
		TimeFieldName: "",
	}
	dv, err := s.KibanaDataView.Create(dv, false, false)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), dv.ID)

	// Get data view
	dvr, err := s.API.KibanaDataView.Get(dv.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), dvr)
	assert.Equal(s.T(), dvr.ID, dv.ID)
	assert.Equal(s.T(), "test", dvr.Title)

	// Update data view
	dv.Title = "test1"
	dv.TimeFieldName = "@timestamp"
	dv, err = s.KibanaDataView.Update(dv, false)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test1", dv.Title)

	// Delete data view
	err = s.KibanaDataView.Delete(dv.ID)
	assert.NoError(s.T(), err)

	// Check deleted
	dv, err = s.API.KibanaDataView.Get(dv.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), dv)
}
