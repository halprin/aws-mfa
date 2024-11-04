package login

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_setSensibleDefaults_profileUnspecified(t *testing.T) {
	profile := ""

	outputProfile, _ := setSensibleDefaults(profile, 0)

	assert.Equal(t, "default", outputProfile)
}
