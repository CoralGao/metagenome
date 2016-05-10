package main

import (
	"fmt"
	"github.com/vtphan/fmic"
	"math/rand"
	"time"
	"os"
	"bytes"
	"bufio"
	"runtime"
)

//-----------------------------------------------------------------------------
func main() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS( runtime.NumCPU() )

	genome := os.Args[1]
	read := os.Args[2]
	saved_idx := fmic.LoadCompressedIndex(genome + ".fmi")
	f, err := os.Open(read)
	check_for_error(err)
	r := bufio.NewReader(f)
	i := 0
	fmt.Println("guessID")
	var read1, read2 []byte
	readsToGenome := make([]int, 50)
	for {
        line, err := r.ReadBytes('\n')
        if err != nil { break }
        if len(line) > 1 {
        	if i % 8 == 1 {
        		read1 = bytes.TrimSpace(line)
        	} else if i % 8 == 5{
        		read2 = bytes.TrimSpace(line)
        		seq := saved_idx.GuessPair(read1, reverse_complement(read2), 100, 1500)
        		fmt.Println(seq)

        		if seq != -1 {
        			readsToGenome[seq] += len(read2) + len(read1)
        		}
        	}
        }
        i++
    }
    
/*    lengths := []int{6337440,5060881,5455081,5484481,5387401,5192521,3331321,5577721,5422321,3904681,5668561,5429041,6056401,5463241,4017841,6205321,4680961,4730401,3719761,4666441,4089841,3784801,5171881,5583361,4107961,3724081,3799081,3919441,3573841,3572041,3718921,5008681,6874801,5314681,5098801,4502761,4539481,5177161,4012561,3818521,5514841,5342401,5293081,4749001,4216801,5730601,4786321,4045441,4348921,3663841}    
    for i:=0;i<50;i++ {
    	fmt.Println(float64(readsToGenome[i])/float64(lengths[i]))
    }*/
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
