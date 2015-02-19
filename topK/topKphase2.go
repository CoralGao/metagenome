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
   "github.com/vtphan/pqueue"
   "strings"
   "strconv"
)
type MyQ struct {
    occurance map[int]int16
    frequency map[int]int32
}
func (M* MyQ) is_better(x, y int) bool {
    if M.occurance[x] > M.occurance[y] {
        return false
    } else if M.occurance[x] < M.occurance[y] {
        return true
    } else {
        return M.frequency[x] > M.frequency[y]
    }
}
func main() {
    if len(os.Args) != 3 {
        panic("Must provide sequence folder and global occurance file.")
    }

    kmer_len := 16
    files, _ := ioutil.ReadDir(os.Args[1])

    numCores := runtime.NumCPU()
    runtime.GOMAXPROCS(numCores)

    globalOccu := make(map[int]int16)
    loadOccu(os.Args[2], globalOccu)

    topKnum := 32
    topKmer := make([][]int, len(files))
    var wg1 sync.WaitGroup
    for index, file := range files {
        wg1.Add(1)
        localFreq := make(map[int]int32)
        kmerFreq(file.Name(), localFreq, kmer_len)
        go func(i int) {
            defer wg1.Done()
            topKmer[i] = localFreqCalc(topKnum,localFreq,globalOccu) 
        }(index)
    }
    wg1.Wait()
    fmt.Println(topKmer)
}
func localFreqCalc(topKnum int, localFreq map[int]int32, globalOccu map[int]int16) []int {
    topK := make([]int, topKnum)
    Q := MyQ{globalOccu, localFreq}
    pq := pqueue.New(topKnum, Q.is_better)

    for i := range localFreq {
        pq.Push(i)
    }
    for i:=0; i<topKnum; i++ {
        key:= pq.Pop()
        topK[i] = key
   }
   return topK
}
func loadOccu(filename string, globalOccu map[int]int16) {
    f,err := os.Open(filename)
    if err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }

    defer f.Close()
    br := bufio.NewReader(f)

    for {
        line, isPrefix, err := br.ReadLine()
        if err != nil || isPrefix {
            break
        } else {
            elements := strings.Split(string(line), " ")
            k,_ := strconv.Atoi(elements[0])
            mapk,_ := strconv.Atoi(elements[1])
            globalOccu[k] = int16(mapk)
        }
    }
}
func contains(s []int, e int) bool {
    for _, a := range s { if a == e { return true } }
    return false
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
func compare(i, j []int) bool {
    // We want Pop to give us the highest, not lowest, priority so we use greater than here.
    // return pq[i].priority > pq[j].priority
    if i[0] < j[0] {
        return true
    } else if i[0] > j[0] {
        return false
    } else{
        return i[1] > j[1]
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