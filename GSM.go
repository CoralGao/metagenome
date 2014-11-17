package main

import (
   "os"
   "fmt"
   "bufio"
   "bytes"
   "sync"
   "runtime"
   "time"
   "sort"
   "io/ioutil"
)

type ByMap struct{
    mapp map[int]int
    keys []int
}

func (s ByMap) Len() int {
    return len(s.keys)
}
func (s ByMap) Swap(i, j int) {
    s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
}
func (s ByMap) Less(i, j int) bool{
    return (s.mapp[s.keys[i]] < s.mapp[s.keys[j]])
}

func main() {
    if len(os.Args) != 2 {
        panic("must provide sequence folder file.")
    }

    files, _ := ioutil.ReadDir(os.Args[1])
    start_time := time.Now()
    gsm := make(map[int]int)

    core_num := 2
    kmer_len := 7
    distance := 0
    runtime.GOMAXPROCS(core_num+2)

    for index, fi := range files {
        fmt.Println(fi.Name())
        f,err := os.Open(os.Args[1] + "/" + fi.Name())
        if err != nil {
            fmt.Printf("%v\n",err)
            os.Exit(1)
        }

        br := bufio.NewReader(f)
        byte_array := bytes.Buffer{}

        _, isPrefix, err := br.ReadLine()

        if err != nil || isPrefix {
            fmt.Printf("%v\n",err)
            os.Exit(1)
        }

        for {
            line , isPrefix, err := br.ReadLine()
            if err != nil || isPrefix{
                break
            } else {
                byte_array.Write([]byte(line))
            }       
        }

        input := []byte(byte_array.String())
        var wg sync.WaitGroup
        result := make(chan int, core_num)

        for i := 0; i < core_num; i++ {
            wg.Add(1)
            go process(input, i, core_num, kmer_len, distance, result, &wg)
        }

        go func() {
            wg.Wait()
            close(result)
        }()

        gsm1 := make(map[int]int)

        for res := range result {
            // if gsm1[res] == 0 {
                gsm1[res] = index + 1
            // }
        }

        for k := range gsm1 {
            if gsm[k] == 0 {
                gsm[k] = gsm1[k]
            } else {
                gsm[k] = -1
            }
        }

        f.Close()

        /*for res := range result {
            if gsm[res] == 0 {
                gsm[res] = index+1
            } else if gsm[res] == index+1 {

            } else {
                gsm[res] = -1
            }
        }*/
    }
    fmt.Println(len(gsm))
    var keys []int
    for k := range gsm {
        keys = append(keys, k)
    }

    sort.Sort(ByMap{gsm, keys})
    for _, k := range keys {
        fmt.Println("Key:", k, "Value:", gsm[k], "end")
    }

    gsm_time := time.Since(start_time)
    fmt.Println("used time", gsm_time)
}

func process(genome []byte, i int, core_num int, kmer_len int, distance int, result chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    begin := len(genome)*i/core_num
    end := len(genome)*(i+1)/core_num
    if begin != 0 {
        begin = begin - 2*kmer_len-distance
    }
    for m := begin; m < end-2*kmer_len-distance; m++ {
        m1 := m
        m2 := m+kmer_len
        m3 := m+kmer_len+distance
        m4 := m+2*kmer_len+distance
        kmer := make([]byte, kmer_len)
        copy(kmer, genome[m1:m2])
        kmer = append(kmer, genome[m3:m4]...)
        repr := 0
        d:
        for j := 0; j<len(kmer); j++ {
            switch kmer[j] {
                case 'A': repr = 4*repr
                case 'C': repr = 4*repr + 1
                case 'G': repr = 4*repr + 2
                case 'T': repr = 4*repr + 3
                default:
                // we skip any qgram that contains a non-standard base, e.g. N
                  repr = -1
                  break d
            }
        }
        if repr!= -1 {
            result <- repr
        }
    }
}