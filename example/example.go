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
	numberOfSampleEntries = 100000
)

func main() {
	// write first half
	parfile, err := os.Create("test.par")
	if err != nil {
		panic(err)
	}

	w := par.NewWriter(parfile)

	w.Header.EntryHeaderLength = 1
	err = w.WriteHeader()
	if err != nil {
		panic(err)
	}

	for i := 0; i < numberOfSampleEntries/2; i++ {
		_, err = w.Write([]byte("test data" + strconv.Itoa(i)))
		if err != nil {
			panic(err)
		}
	}
	w.Flush()
	parfile.Close()

	// append second half
	parfile1, _ := os.OpenFile("test.par", os.O_WRONLY|os.O_APPEND, 0666)

	w1 := par.NewWriter(parfile1)
	w1.Header.EntryHeaderLength = 1

	for i := numberOfSampleEntries / 2; i < numberOfSampleEntries; i++ {
		_, err = w1.Write([]byte("test data" + strconv.Itoa(i)))
		if err != nil {
			panic(err)
		}
	}
	w1.Flush()
	parfile1.Close()

	// read
	parfile2, _ := os.Open("test.par")
	r := par.NewReader(parfile2)

	fmt.Println("EntryHeaderLength: ", r.Header.EntryHeaderLength)

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

	fmt.Printf("found %d / %d entries\n", count, numberOfSampleEntries)
	parfile2.Close()

	// stats raw
	plain, _ := os.Open("test.par")
	defer plain.Close()

	stat, err := plain.Stat()
	if err != nil {
		return
	}
	fmt.Printf("File size is %5d kiB, avg. %2d-bytes per entry\n", stat.Size()/1024, stat.Size()/numberOfSampleEntries)

	// stats compressed
	zfile, err := os.Create("test.par.gz")
	if err != nil {
		panic(err)
	}
	defer plain.Close()

	zw := gzip.NewWriter(zfile)
	defer zw.Close()

	// chunks of 4MB
	chunk := make([]byte, 4<<20)

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
	fmt.Printf("File size is %5d kiB, avg. %2d-bytes per entry\n", zstat.Size()/1024, zstat.Size()/numberOfSampleEntries)
}
