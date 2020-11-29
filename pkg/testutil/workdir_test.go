package testutil

import (
	"fmt"
	"testing"
)

func TestFindRoot(t *testing.T) {
	wd := findWorkDir()
	fmt.Println(wd)
}
