package login

import (
	"github.com/stretchr/testify/assert"
	"os"
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

func Test_getMfaCodeUnlessSpecified_returnsSameCodeAsProvided(t *testing.T) {
	inputMfaCode := "266226"

	outputMfaCode, err := getMfaCodeUnlessSpecified(inputMfaCode)

	assert.NoError(t, err)
	assert.Equal(t, inputMfaCode, outputMfaCode)
}

func Test_getMfaCodeUnlessSpecified_returnsCodeFromStdin(t *testing.T) {
	stdinMfaCode := "266226"

	readPipe, writePipe, err := os.Pipe()
	assert.NoError(t, err)

	defer readPipe.Close()
	defer writePipe.Close()

	originalStdin := os.Stdin
	defer func() {
		// reset STDIN back to real STDIN
		os.Stdin = originalStdin
	}()
	os.Stdin = readPipe

	_, err = writePipe.WriteString(stdinMfaCode + "\n")
	assert.NoError(t, err)

	outputMfaCode, err := getMfaCodeUnlessSpecified("")

	assert.NoError(t, err)
	assert.Equal(t, stdinMfaCode, outputMfaCode)
}
