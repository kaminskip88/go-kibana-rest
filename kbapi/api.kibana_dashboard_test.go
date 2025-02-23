package kbapi

import (
	"encoding/json"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

func (s *KBAPITestSuite) TestKibanaDashboard() {

	// Import dashboard from fixtures
	b, err := ioutil.ReadFile("../fixtures/kibana-dashboard.json")
	if err != nil {
		panic(err)
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(b, &data)
	assert.NoError(s.T(), err)
	err = s.API.KibanaDashboard.Import(data, nil, true, "default")
	assert.NoError(s.T(), err)

	// Export dashboard
	data, err = s.API.KibanaDashboard.Export([]string{"edf84fe0-e1a0-11e7-b6d5-4dc382ef7f5b"}, "default")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), data)

	// Import dashboard from fixtures in specific space
	b, err = ioutil.ReadFile("../fixtures/kibana-dashboard.json")
	if err != nil {
		panic(err)
	}
	data = make(map[string]interface{})
	err = json.Unmarshal(b, &data)
	assert.NoError(s.T(), err)
	err = s.API.KibanaDashboard.Import(data, nil, true, "testacc")
	assert.NoError(s.T(), err)

	// Export dashboard from specific space
	data, err = s.API.KibanaDashboard.Export([]string{"edf84fe0-e1a0-11e7-b6d5-4dc382ef7f5b"}, "testacc")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), data)

}
