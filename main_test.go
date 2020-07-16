package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
