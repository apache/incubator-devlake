package main

import "testing"

const LOCAL_REPO_PATH = "/home/klesh/Projects/merico/tidb"

/*
goos: linux
goarch: amd64
pkg: github.com/merico-dev/lake/scripts
cpu: AMD Ryzen 7 5800H with Radeon Graphics
BenchmarkGit2GoOdb-16             1        88365798286 ns/op       53370704 B/op    2302267 allocs/op
PASS
ok      github.com/merico-dev/lake/scripts      88.383s
*/
func BenchmarkGit2GoOdb(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetCommitsByOdbIteration(LOCAL_REPO_PATH)
	}
}

func BenchmarkGit2GoWalk(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetCommitsByRevWalk(LOCAL_REPO_PATH)
	}
}
