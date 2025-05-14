package main

import (
	"testing"
	"time"

	"github.com/jmsilvadev/go-pack-optimizer/pkg/config"
)

func TestRun(t *testing.T) {
	c := config.GetDefaultConfig()
	go run(c)
	time.Sleep(time.Second)
}
