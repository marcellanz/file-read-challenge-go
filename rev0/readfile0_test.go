package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func BenchmarkReadfile0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readfile0("../itcont.txt", log.New(ioutil.Discard, "", log.LstdFlags))
	}
}
