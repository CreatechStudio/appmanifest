package main

import (
	"bytes"
	"os"
	"regexp"
	"testing"
)

func TestCalucalteMD5s(t *testing.T) {
	ok := []string{
		"f1c9645dbc14efddc7d8a322685f26eb",
		"f1c9645dbc14efddc7d8a322685f26eb",
		"f1c9645dbc14efddc7d8a322685f26eb",
		"93b885adfe0da089cdf634904fd59f71",
	}
	buf := make([]byte, DefaultMD5Size*3+1)
	r := bytes.NewReader(buf)
	md5s, err := calculateMD5s(r, DefaultMD5Size)
	if err != nil {
		t.Fatal(err)
	}
	if len(md5s) != len(ok) {
		t.Fatal("expected", len(ok), "got", len(md5s))
	}
	for i, h := range md5s {
		if ok[i] != h {
			t.Fatal("expected", ok[i], "got", h)
		}
	}
}

/*
Benchmark10MB-8  	  500000	      3297 ns/op	   32896 B/op	       3 allocs/op
Benchmark100MB-8 	  300000	      3528 ns/op	   32897 B/op	       3 allocs/op
Benchmark1000MB-8	       1	1967685209 ns/op	 3334064 B/op	     829 allocs/op
*/

func benchmarkSize(b *testing.B, size int) {
	var buf = make([]byte, size)
	r := bytes.NewReader(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateMD5s(r, DefaultMD5Size)
	}
}

func Benchmark10MB(b *testing.B) {
	benchmarkSize(b, DefaultMD5Size)
}

func Benchmark100MB(b *testing.B) {
	benchmarkSize(b, DefaultMD5Size*10)
}

func Benchmark1000MB(b *testing.B) {
	benchmarkSize(b, DefaultMD5Size*100)
}

func TestCreateAppManifestMD5SizeSmallFileUsesFileSize(t *testing.T) {
	tmp, err := os.CreateTemp("", "appmanifest-small-*.pkg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	const fileSize = 1024
	if err := tmp.Truncate(fileSize); err != nil {
		t.Fatal(err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := createAppManifest(tmp.Name(), "", &out, DefaultMD5Size); err != nil {
		t.Fatal(err)
	}

	if !regexp.MustCompile(`(?s)<key>md5-size</key>\s*<integer>1024</integer>`).MatchString(out.String()) {
		t.Fatalf("expected md5-size to be file size %d, plist was:\n%s", fileSize, out.String())
	}
}

func TestCreateAppManifestMD5SizeLargeFileUsesDefaultChunkSize(t *testing.T) {
	tmp, err := os.CreateTemp("", "appmanifest-large-*.pkg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	if err := tmp.Truncate(DefaultMD5Size + 1); err != nil {
		t.Fatal(err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := createAppManifest(tmp.Name(), "", &out, DefaultMD5Size); err != nil {
		t.Fatal(err)
	}

	if !regexp.MustCompile(`(?s)<key>md5-size</key>\s*<integer>10485760</integer>`).MatchString(out.String()) {
		t.Fatalf("expected md5-size to be %d for files larger than 10MB, plist was:\n%s", DefaultMD5Size, out.String())
	}
}
