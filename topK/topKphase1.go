package main

import (
   "github.com/vtphan/kmers"
   "os"
   "sync"
   "bufio"
   "fmt"
   "runtime"
   "io/ioutil"
   "bytes"
)

func main() {
    if len(os.Args) != 2 {
        panic("Must provide sequence folder.")
    }

    kmer_len := 16
    files, _ := ioutil.ReadDir(os.Args[1])

    numCores := runtime.NumCPU()
    runtime.GOMAXPROCS(numCores)

    globalOccu := make(map[int]int16)
    for _, file := range files {
        localFreq := make(map[int]int32)
        kmerFreq(file.Name(), localFreq, kmer_len)
        for k := range localFreq {
            globalOccu[k]++
        }
    }

    for k := range globalOccu {
        fmt.Println(k, globalOccu[k])
    }
}

func kmerFreq(filename string, localFreq map[int]int32, kmer_len int) {
        genome := readGenome(os.Args[1] + "/" + filename)
        var wg sync.WaitGroup
        numCores := runtime.NumCPU()
        result := make(chan int, numCores)

        for i := 0; i < numCores; i++ {
            wg.Add(1)
            go func(i int) {
                defer wg.Done()
                start := len(genome)*i/numCores
                end := len(genome)*(i+1)/numCores
                if start != 0 {
                    start = start - kmer_len + 1
                }
                kmers.Slide(genome, kmer_len, start, end, result)
            }(i)
        }
        go func() {
            wg.Wait()
            close(result)  
        }()
        for k := range result{
            localFreq[k]++
        }
}
func readGenome(filename string) []byte{
    f,err := os.Open(filename)
    if err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }

    defer f.Close()
    // fi, _ := f.Stat()
    br := bufio.NewReader(f)
    byte_buffer := bytes.Buffer{}

    _, isPrefix, err := br.ReadLine()
    if err != nil || isPrefix {
        fmt.Printf("%v\n", nil)
        os.Exit(1)
    }
    // fmt.Fprintf(infofile, "%s\t%d\n", head, fi.Size())
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