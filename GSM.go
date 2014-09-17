package main

import (
   "os"
   "fmt"
   "bufio"
   "bytes"
   "sync"
   "runtime"
   "time"
   "io/ioutil"
)

func main() {
    if len(os.Args) != 2 {
        panic("must provide sequence folder file.")
    }

    files, _ := ioutil.ReadDir(os.Args[1])
    start_time := time.Now()
    gsm := make(map[int]int)

    for j, fi := range files {
    	f,err := os.Open(os.Args[1] + "/" + fi.Name())
        if err != nil {
            fmt.Printf("%v\n",err)
            os.Exit(1)
        }

        defer f.Close()
        br := bufio.NewReader(f)
        byte_array := bytes.Buffer{}

        gname, isPrefix, err := br.ReadLine()

        if err != nil || isPrefix {
    		fmt.Printf("%v\n",err)
    		os.Exit(1)    	
        }

        header := make([]byte, len(gname))
        copy(header, gname)

        for {
            line , isPrefix, err := br.ReadLine()
            if err != nil || isPrefix{
                break
            } else {
                byte_array.Write([]byte(line))
            }    	
        }

        input := []byte(byte_array.String())
        fmt.Println(len(input))
        var wg sync.WaitGroup
        core_num := 3
        kmer_len := 15
        result := make(chan int, 1000000)
        runtime.GOMAXPROCS(core_num)

        go func() {
            wg.Wait()
            close(result)
        }()

        for i := 0; i < core_num; i++ {
            wg.Add(1)
        	go func(genome []byte, index int, core_num int, result chan int) {
                defer wg.Done()
        		begin := len(genome)*index/core_num
        		end := len(genome)*(index+1)/core_num
                if begin != 0 {
                    begin = begin - kmer_len
                }
        		fmt.Println(begin, end)
        		for i := begin; i < end-kmer_len; i++ {
    				kmer := genome[i:i+kmer_len]
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
        	}(input, i, core_num, result)
        }

        for res := range result {
            if gsm[res] == 0 {
                gsm[res] = j+1
            } else if gsm[res] == j+1 {
                break
            } else {
                gsm[res] = -1
            }
        }
    }

    fmt.Println(gsm)

    gsm_time := time.Since(start_time)
    fmt.Println("used time", gsm_time)
}