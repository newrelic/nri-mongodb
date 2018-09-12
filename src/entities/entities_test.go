package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultCollector_GetSession_Error(t *testing.T) {
	coll := new(defaultCollector)
	session, err := coll.GetSession()
	assert.Error(t, err)
	assert.Nil(t, session)
}

func Test_logError_EmptyFormat(t *testing.T) {
	assert.NotPanics(t, func() {
		logError(assert.AnError, "")
	})
}

func Test_extractHostPort(t *testing.T) {
	type result struct {
		host string
		port string
	}
	type testCase struct {
		input    string
		expected result
	}
	tests := []testCase{
		{
			"host:123",
			result{
				host: "host",
				port: "123",
			},
		},
		{
			"host",
			result{
				host: "host",
				port: "27017",
			},
		},
		{
			"host:",
			result{
				host: "host",
				port: "27017",
			},
		},
	}
	for _, tc := range tests {
		hp := extractHostPort(tc.input)
		assert.Equal(t, tc.expected.host, hp.Host)
		assert.Equal(t, tc.expected.port, hp.Port)
	}
}
