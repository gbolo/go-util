package bench

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"
)

func benchmarkRun(h hash.Hash, i int, b *testing.B) {
	bs := make([]byte, i)
	_, err := rand.Read(bs)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(bs)
		h.Sum(nil)
	}

}

func BenchmarkMD5_1k(b *testing.B) {
	benchmarkRun(md5.New(), 1024, b)
}

func BenchmarkMD5_10k(b *testing.B) {
	benchmarkRun(md5.New(), 10*1024, b)
}

func BenchmarkMD5_100k(b *testing.B) {
	benchmarkRun(md5.New(), 100*1024, b)
}

func BenchmarkMD5_250k(b *testing.B) {
	benchmarkRun(md5.New(), 250*1024, b)
}

func BenchmarkMD5_500k(b *testing.B) {
	benchmarkRun(md5.New(), 500*1024, b)
}

func BenchmarkSHA1_1k(b *testing.B) {
	benchmarkRun(sha1.New(), 1024, b)
}

func BenchmarkSha1_10k(b *testing.B) {
	benchmarkRun(sha1.New(), 10*1024, b)
}

func BenchmarkSha1_100k(b *testing.B) {
	benchmarkRun(sha1.New(), 100*1024, b)
}

func BenchmarkSha1_250k(b *testing.B) {
	benchmarkRun(sha1.New(), 250*1024, b)
}

func BenchmarkSha1_500k(b *testing.B) {
	benchmarkRun(sha1.New(), 500*1024, b)
}

func BenchmarkSha256_1k(b *testing.B) {
	benchmarkRun(sha256.New(), 1024, b)
}

func BenchmarkSha256_10k(b *testing.B) {
	benchmarkRun(sha256.New(), 10*1024, b)
}

func BenchmarkSha256_100k(b *testing.B) {
	benchmarkRun(sha256.New(), 100*1024, b)
}

func BenchmarkSha256_250k(b *testing.B) {
	benchmarkRun(sha256.New(), 250*1024, b)
}

func BenchmarkSha256_500k(b *testing.B) {
	benchmarkRun(sha256.New(), 500*1024, b)
}

func BenchmarkSha512_1k(b *testing.B) {
	benchmarkRun(sha512.New(), 1024, b)
}

func BenchmarkSha512_10k(b *testing.B) {
	benchmarkRun(sha512.New(), 10*1024, b)
}

func BenchmarkSha512_100k(b *testing.B) {
	benchmarkRun(sha512.New(), 100*1024, b)
}

func BenchmarkSha512_250k(b *testing.B) {
	benchmarkRun(sha512.New(), 250*1024, b)
}

func BenchmarkSha512_500k(b *testing.B) {
	benchmarkRun(sha512.New(), 500*1024, b)
}


