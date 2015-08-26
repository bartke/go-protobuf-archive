package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"

	par "github.com/bartke/go-protobuf-archive"
)

const (
	numberOfSampleEntries = 128 * 1024
)

func main() {
	parfile, err := os.Create("test.par")
	if err != nil {
		panic(err)
	}
	defer parfile.Close()

	// write
	w := par.NewWriter(parfile)

	w.Header.EntryHeaderLength = 1
	w.Header.NumberOfEntries = numberOfSampleEntries
	err = w.WriteHeader()
	if err != nil {
		panic(err)
	}

	for i := 0; i < numberOfSampleEntries; i++ {
		_, err = w.Write([]byte("test data" + strconv.Itoa(i)))
		if err != nil {
			panic(err)
		}
	}

	w.Flush()

	// read
	parfile2, err := os.Open("test.par")
	r := par.NewReader(parfile2)

	fmt.Println("EntryHeaderLength: ", r.Header.EntryHeaderLength)
	fmt.Println("Number of Entries: ", r.Header.NumberOfEntries)

	var count int
	for {
		h, err := r.Next()
		if err != nil {
			if err == io.EOF {
				// we go here after the last element is read
				break
			}
			panic(err)
		}
		buf := make([]byte, h.Length)
		_, err = r.Read(buf)
		if err != nil {
			panic(err)
		}

		//fmt.Println(string(buf))
		count++
	}

	fmt.Println(count, "entries found")

	// raw
	stat, err := parfile.Stat()
	if err != nil {
		return
	}

	fmt.Printf("File size is %5d kiB, avg. is %2d-bytes per entry\n", stat.Size()/1024, stat.Size()/numberOfSampleEntries)

	// compressed
	plain, err := os.Open("test.par")
	if err != nil {
		panic(err)
	}
	defer plain.Close()

	zfile, err := os.Create("test.par.gz")
	if err != nil {
		panic(err)
	}
	defer plain.Close()

	zw := gzip.NewWriter(zfile)
	defer zw.Close()

	// chunk
	chunk := make([]byte, 4<<20) // Read 4MB at a time

	for {
		n, err := plain.Read(chunk)
		if n > 0 {
			zw.Write(chunk)
		}
		if err != nil {
			if err == io.EOF {
				// we go here after the last element is read
				break
			}
			panic(err)
		}
	}
	zw.Flush()

	zstat, err := zfile.Stat()
	if err != nil {
		return
	}

	fmt.Printf("File size is %5d kiB, avg. is %2d-bytes per entry\n", zstat.Size()/1024, zstat.Size()/numberOfSampleEntries)
}
