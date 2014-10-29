package main

import (
   "os"
   "fmt"
   "bufio"
   "sync"
   "runtime"
   "time"
   "math"
   "sort"
)

func main() {
    if len(os.Args) != 2 {
        panic("must provide sequence folder file.")
    }

    start_time := time.Now()
    gsm := make(map[int]int)

    core_num := 2
    kmer_len := 5
    distance := 10
    runtime.GOMAXPROCS(core_num+2)

    reads := make(chan []byte, core_num)
    go ReadReads(reads, os.Args[1])

    var wg sync.WaitGroup
    result := make(chan int)

    for i := 0; i < core_num; i++ {
        wg.Add(1)
    	go process(reads, kmer_len, distance, result, &wg)
    }

    go func() {
        wg.Wait()
        close(result)
    }()

    for res := range result {
        gsm[res] = 1
    }

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

func ReadReads(reads chan []byte, file string) {
    fmt.Println(file)
    f,err := os.Open(file)
    if err != nil {
        fmt.Printf("%v\n",err)
        os.Exit(1)
    }
    defer f.Close()
    br := bufio.NewReader(f)

    _, isPrefix, err := br.ReadLine()

    if err != nil || isPrefix {
        fmt.Printf("%v\n",err)
        os.Exit(1)
    }

    i := 0
    for {
        i++
        line , isPrefix, err := br.ReadLine()
        if err != nil || isPrefix{
            break
        } else {
            if math.Mod(float64(i),2) == 0 {
                reads <- []byte(line)
            }
        }
    }
    close(reads)
}
func process(reads chan []byte, kmer_len int, distance int, result chan int, wg *sync.WaitGroup) {
    defer wg.Done()

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