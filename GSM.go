package main

import (
   "os"
   "fmt"
   "bufio"
   "bytes"
   "sync"
   "runtime"
   "sort"
   "io/ioutil"
   "github.com/vtphan/kmers"
)

func main() {
    if len(os.Args) != 2 {
        panic("must provide sequence folder file.")
    }

    files, _ := ioutil.ReadDir(os.Args[1])
    gsm := make(map[int]int)

    numCores := runtime.NumCPU()
    runtime.GOMAXPROCS(numCores)

    kmer_len := 14
    for index, file := range files {
        genome := readGenome(os.Args[1] + "/" + file.Name())
        var wg sync.WaitGroup
        result := make(chan int, numCores)

        for i := 0; i < numCores; i++ {
            wg.Add(1)
            go func(i int) {
                defer wg.Done()
                start := len(genome)*i/numCores
                end := len(genome)*(i+1)/numCores
                if start != 0 {
                    start = start - kmer_len
                }
                fmt.Println(start, end)
                kmers.Slide(genome, kmer_len, start, end, result)
            }(i)
        }

        go func() {
            wg.Wait()
            close(result)
        }()
        for k := range result {
            if gsm[k] == 0 {
                gsm[k] = index + 1
            } else if gsm[k] == index + 1 {
                
            } else {
                gsm[k] = -1
            }
        }
    }

    var keys []int
    for k := range gsm {
        keys = append(keys, k)
    }    

    sort.Ints(keys)
    for _, k := range keys {
        fmt.Println("Key:", kmers.NumToKmer(k, kmer_len),k, "Value:", gsm[k], "end")
    }
}

func readGenome(filename string) []byte{
    f,err := os.Open(filename)
    if err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }

    defer f.Close()
    br := bufio.NewReader(f)
    byte_buffer := bytes.Buffer{}

    _, isPrefix, err := br.ReadLine()
    if err != nil || isPrefix {
        fmt.Printf("%v\n", nil)
        os.Exit(1)
    }
    for {
        line, isPrefix, err := br.ReadLine()
        if err != nil || isPrefix {
            break
        } else {
            if bytes.Contains(line, []byte(">")) {
                byte_buffer.Write([]byte("NNNNNNNNNNNNNN"))
            } else {
                byte_buffer.Write(line)
            }
        }
    }

    return byte_buffer.Bytes()
}