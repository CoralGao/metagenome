package main

import (
	"os"
	"fmt"
	"encoding/csv"
	"bufio"
	"math"
	"runtime"
	"sync"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		panic("must provide readfile and GSM from genomes and read GSM result file name")
	}

	resultfile, err := os.Create(os.Args[3]+".csv")
	if err != nil {
		fmt.Println("%v\n", err)
		os.Exit(1)
	}

	rw := csv.NewWriter(resultfile)
	head := make([]string, 2)
	head[0] = "kmer"
	head[1] = "b"
	returnError := rw.Write(head)
	if returnError != nil {
		fmt.Println(returnError)
	}

	rw.Flush()

	resultRead := make(chan int, 10)
	go CountFreq(os.Args[1], 14, resultRead)
	gsmread := make(map[int]int)
	for res := range resultRead {
		gsmread[res] = gsmread[res] + 1
	}

	csvfile, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvfile.Close()
	reader := csv.NewReader(csvfile)

	rawcsvdata, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, each := range rawcsvdata {
		if i != 0 {
			line := make([]string, 2)
			gsmid, _ := strconv.Atoi(each[0])
			line[0] = strconv.Itoa(gsmid)
			line[1] = strconv.Itoa(gsmread[gsmid])

			returnError := rw.Write(line)
			if returnError != nil {
				fmt.Println(returnError)
			}
		}
	}
	rw.Flush()
}

func CountFreq(readFile string, K int, result chan int) {

   // Get all reads into a channel
   reads := make(chan []byte)
   go func() {
      f, err := os.Open(readFile)
      if err != nil {
         panic("Error opening " + readFile)
      }
      defer f.Close()
      scanner := bufio.NewScanner(f)
      i := 0.0
      for scanner.Scan() {
        if math.Mod(i, 2.0) == 1 {
         reads <- []byte(scanner.Text())
        }
        i++
      }
      close(reads)
   }()

   // Spread the processing of reads into different cores
   numCores := runtime.NumCPU()
   runtime.GOMAXPROCS(numCores)
   var wg sync.WaitGroup

   for i:=0; i<numCores; i++ {
      wg.Add(1)
      go func() {
         defer wg.Done()
         ProcessRead(reads, K/2, 0, result)
      }()
   }
   go func() {
        wg.Wait()
        close(result)
   }()
}
func ProcessRead(reads chan []byte, kmer_len int, distance int, result chan int) {
    for read := range reads {
        for m := 0; m < len(read) - 2*kmer_len - distance; m++ {
            m1 := m
            m2 := m+kmer_len
            m3 := m+kmer_len+distance
            m4 := m+2*kmer_len+distance
            kmer := make([]byte, 2*kmer_len)
            copy(kmer, read[m1:m2])
            kmer = append(kmer, read[m3:m4]...)
            repr := 0
            for j := 0; j<len(kmer); j++ {
                switch kmer[j] {
                    case 'A': repr = 4*repr
                    case 'C': repr = 4*repr + 1
                    case 'G': repr = 4*repr + 2
                    case 'T': repr = 4*repr + 3
                    default:
                    // we skip any qgram that contains a non-standard base, e.g. N
                      repr = repr
                }
            }
            result <- repr
        }
    }
}