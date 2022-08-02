//go:build component

package test

import (
	"testing"
)

func Test_404(t *testing.T) {
	test404(t, "/not/a/valid/path", "Route Not Found")
}
