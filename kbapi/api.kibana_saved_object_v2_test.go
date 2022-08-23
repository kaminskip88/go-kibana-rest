package kbapi

import (
	"github.com/stretchr/testify/assert"
)

func (s *KBAPITestSuite) TestKibanaSaveObjectV2() {

	// Create new saved object
	so := &KibanaSavedObject{
		ID:   "test",
		Type: "index-pattern",
		Attributes: map[string]interface{}{
			"title": "test-pattern-*",
		},
		References: []Reference{{
			ID:   "1",
			Name: "one",
			Type: "index-pattern",
		}},
	}
	so, err := s.KibanaSavedObjectV2.Create(so, false)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), so.UpdatedAt)

	// Get saved object by id and type
	sor, err := s.KibanaSavedObjectV2.Get(so.ID, so.Type)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), sor)
	assert.Equal(s.T(), sor.UpdatedAt, so.UpdatedAt)
	assert.Equal(s.T(), "test", so.ID)
	assert.NotEmpty(s.T(), so.References)

	// Update saved object
	so.Attributes = map[string]interface{}{
		"title": "test-pattern-xxx-*",
	}
	so, err = s.KibanaSavedObjectV2.Update(so)
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), so.Attributes, "title")
	assert.Equal(s.T(), so.Attributes["title"], "test-pattern-xxx-*")

	// Delete saved object
	err = s.KibanaSavedObjectV2.Delete(so.ID, so.Type)
	assert.NoError(s.T(), err)

	// Check deleted
	so, err = s.API.KibanaSavedObjectV2.Get(so.ID, so.Type)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), so)
}
