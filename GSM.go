package main

import (
   "os"
   "fmt"
   "bufio"
   "bytes"
   "sync"
   // "runtime"
   "time"
   "sort"
   "io/ioutil"
)

func main() {
    if len(os.Args) != 2 {
        panic("must provide sequence folder file.")
    }

    files, _ := ioutil.ReadDir(os.Args[1])
    start_time := time.Now()
    gsm := make(map[int]int)

    core_num := 3
    kmer_len := 5
    distance := 10
    // runtime.GOMAXPROCS(core_num+2)

    for index, fi := range files {
        fmt.Println(fi.Name())
    	f,err := os.Open(os.Args[1] + "/" + fi.Name())
        if err != nil {
            fmt.Printf("%v\n",err)
            os.Exit(1)
        }

        defer f.Close()
        br := bufio.NewReader(f)
        byte_array := bytes.Buffer{}

        _, isPrefix, err := br.ReadLine()

        if err != nil || isPrefix {
    		fmt.Printf("%v\n",err)
    		os.Exit(1)
        }

        // header := make([]byte, len(gname))
        // copy(header, gname)

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
        result := make(chan int, core_num)

        for i := 0; i < core_num; i++ {
            wg.Add(1)
            fmt.Println(i)
        	go process(input, i, core_num, kmer_len, distance, result, &wg)
        }

        go func() {
            wg.Wait()
            fmt.Println("close result")
            close(result)
        }()

        for res := range result {
            if gsm[res] == 0 {
                gsm[res] = index+1
            } else if gsm[res] == index+1 {

            } else {
                gsm[res] = -1
            }
        }
    }
    // fmt.Println(gsm)
    fmt.Println(len(gsm))
    var keys []int
    for k := range gsm {
        keys = append(keys, k)
    }

    sort.Ints(keys)
    for _, k := range keys {
        fmt.Println("Key:", k, "Value:", gsm[k])
    }

    gsm_time := time.Since(start_time)
    fmt.Println("used time", gsm_time)
}

func process(genome []byte, i int, core_num int, kmer_len int, distance int, result chan int, wg *sync.WaitGroup) {
    fmt.Println("iN", i)
    defer wg.Done()
    fmt.Println("loopin", i)
    begin := len(genome)*i/core_num
    end := len(genome)*(i+1)/core_num
    if begin != 0 {
        begin = begin - kmer_len
    }
    fmt.Println(begin, end)
    for m := begin; m < end-2*kmer_len-distance; m++ {
        // fmt.Println("m", m)
        kmer := append(genome[m:m+kmer_len], genome[m+kmer_len+distance:m+2*kmer_len+distance]...)
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