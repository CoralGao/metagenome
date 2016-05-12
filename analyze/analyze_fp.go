package main

import (
	"fmt"
	"github.com/vtphan/fmic"
	"math/rand"
	"time"
	"os"
	"bufio"
	"bytes"
	"strconv"
	"runtime"
)

//-----------------------------------------------------------------------------
func main() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS( runtime.NumCPU() )

	genome_index := os.Args[1]
	read := os.Args[2]
	genome_id, _ := strconv.Atoi(os.Args[3])
	fmt.Println("loading index...")
	idx := fmic.LoadCompressedIndex(genome_index)
	f, err := os.Open(read)
	check_for_error(err)
	r := bufio.NewReader(f)
	i := 0
	tp, fp, fn := 0, 0, 0
	var read1, read2 []byte

	fmt.Println("querying reads...")
	for {
		line, err := r.ReadBytes('\n')
		if err != nil { break }
		if len(line) > 1 {
			if i%4 == 1 {
				read1 = bytes.TrimSpace(line)
			} else if i%4 == 3 {
				read2 = bytes.TrimSpace(line)
				// randomized
				seq := idx.GuessPair(read1, reverse_complement(read2), 100, 1500)

				// deterministic
				// seq := idx.GuessPairD(read1, reverse_complement(read2))
				if seq == genome_id {
					tp++
				} else {
					if seq == -1 {
						fn++
					} else {
						fp++
					}
				}
				// fmt.Println(string(read1))
				// fmt.Println(string(read2))
				// fmt.Println(seq,"\n")
			}
		}
		i++
	}
	fmt.Println(tp, fp, fn, float64(tp)/float64(tp+fp), float64(tp)/float64(tp+fn))
}

func reverse_complement(s []byte) []byte {
	rs := make([]byte, len(s))
	for i:=0; i<len(s); i++ {
		if s[i] == 'A' {
			rs[len(s)-i-1] = 'T'
		} else if s[i] == 'T' {
			rs[len(s)-i-1] = 'A'
		} else if s[i] == 'C' {
			rs[len(s)-i-1] = 'G'
		} else if s[i] == 'G' {
			rs[len(s)-i-1] = 'C'
		} else {
			rs[len(s)-i-1] = s[i]
		}
	}
	return rs
}
func check_for_error(e error) {
	if e != nil {
		panic(e)
	}
}
