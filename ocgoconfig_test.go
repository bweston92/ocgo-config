package ocgoconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupWithInvalidURL(t *testing.T) {
	err := Setup("&://a?^")
	assert.Equal(t, ErrInvalidConfiguration, err)
}

func TestSetupWithInvalidExporterName(t *testing.T) {
	err := Setup("unknown://abc")
	assert.Equal(t, ErrUnknownExporter, err)
}

func TestAvailable(t *testing.T) {
	assert.Contains(t, Available(), "jaeger")
	assert.Contains(t, Available(), "stackdriver")
}

func TestSetupStackdriver(t *testing.T) {
	err := Setup("stackdriver://projectid")

	assert.Nil(t, err)
}
