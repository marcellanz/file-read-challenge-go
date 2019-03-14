package main

import "testing"

func BenchmarkReadfile0(b *testing.B) {
	readfile0("../indiv18/itcont.txt")
}
