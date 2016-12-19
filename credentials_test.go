package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCredentialsReadFromEnvironment(t *testing.T) {
	os.Setenv("DOCKERHUB_USERNAME", "user")
	os.Setenv("DOCKERHUB_PASSWORD", "pass")

	username, password := credentials()

	assert.Equal(t, "user", username)
	assert.Equal(t, "pass", password)
}
