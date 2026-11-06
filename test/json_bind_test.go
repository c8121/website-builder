package test

import (
	"os"
	"testing"

	"github.com/c8121/website-builder/internal/builder"
)

func TestJsonBind(t *testing.T) {

	html := `<div v-for="item in items">{{ item }}</div>`
	json := `{ "items": [ 1, 2, 3 ] }`

	err := builder.Bind(html, json, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}

}
