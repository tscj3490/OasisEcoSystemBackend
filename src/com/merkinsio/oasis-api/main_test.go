package main

import (
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()

	// TODO: Initialize the mongo testing db

	code := m.Run()

	// TODO: Restore the mongo testing db

	os.Exit(code)
}