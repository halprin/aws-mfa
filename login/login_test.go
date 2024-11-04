package login

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_setSensibleDefaults_profileAndDurationUnspecified(t *testing.T) {
	outputProfile, outputDuration := setSensibleDefaults("", 0)

	assert.Equal(t, "default", outputProfile)
	assert.Equal(t, 43200, outputDuration)
}

func Test_setSensibleDefaults_profileAndDurationSpecified(t *testing.T) {
	inputProfile := "DogCow"
	inputDuration := 26

	outputProfile, outputDuration := setSensibleDefaults(inputProfile, inputDuration)

	assert.Equal(t, inputProfile, outputProfile)
	assert.Equal(t, inputDuration, outputDuration)
}
