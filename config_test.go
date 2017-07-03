package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{}
	assert.Equal(cfg.expiry, 0)
	assert.Equal(cfg.fqdn, "")
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{expiry: 1800, fqdn: "https://localhost"}
	assert.Equal(cfg.expiry, 1800)
	assert.Equal(cfg.fqdn, "https://localhost")
}
