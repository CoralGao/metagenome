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
   "container/heap"
   "encoding/csv"
   "strconv"
)
// An Item is something we manage in a priority queue.
type Item struct {
    value    int // The value of the item; arbitrary.
    priority []int    // The priority of the item in the queue.
    // The index is needed by update and is maintained by the heap.Interface methods.
    index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    // We want Pop to give us the highest, not lowest, priority so we use greater than here.
    // return pq[i].priority > pq[j].priority
    if pq[i].priority[0] < pq[j].priority[0] {
        return false
    } else if pq[i].priority[0] > pq[j].priority[0] {
        return true
    } else{
        return pq[i].priority[1] < pq[j].priority[1]
    }
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*Item)
    item.index = n
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    item.index = -1 // for safety
    *pq = old[0 : n-1]
    return item
}

func main() {
    if len(os.Args) != 3 {
        panic("Must provide sequence folder and result file name.")
    }

    kmer_len := 16
    files, _ := ioutil.ReadDir(os.Args[1])

    resultfile, err := os.Create(os.Args[2]+".csv")
    if err != nil {
        fmt.Printf("%v\n",err)
        os.Exit(1)
    }
    rw := csv.NewWriter(resultfile)
    head := make([]string, len(files)+1)
    head[0] = "kmer"
    for index, fi := range files {
        head[index+1] = fi.Name()
    }

    returnError := rw.Write(head)
    if returnError != nil {
        fmt.Println(returnError)
    }
    rw.Flush()

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

    topKnum := 100
    var item *Item
    kmernum := topKnum*len(files)
    topK := make([]int, kmernum)
    for i := range topK{
        topK[i] = -1
    }
    m := 0
    for _, file := range files {
        fmt.Println(file.Name())
        localFreq := make(map[int]int32)
        kmerFreq(file.Name(), localFreq, kmer_len)
        pq := make(PriorityQueue, topKnum)
        k := 0
        for j := range localFreq {
            if !contains(topK, j) {
                if k < topKnum {
                    pq[k] = &Item{
                        value:  j,
                        priority:   []int{int(globalOccu[j]), int(localFreq[j])},
                        index:  k,
                    }
                    k++
                    if k == topKnum {
                        heap.Init(&pq)
                    }
                } else {
                    item = &Item{
                        value:  j,
                        priority: []int{int(globalOccu[j]), int(localFreq[j])},
                    }
                    if compare(item.priority, pq[0].priority) {
                        _ = heap.Pop(&pq).(*Item)
                        heap.Push(&pq, item)
                    }
                }
            }
        }
        for pq.Len() > 0 {
            item := heap.Pop(&pq).(*Item)
            fmt.Println(kmers.NumToKmer(item.value, kmer_len))
            topK[m] = item.value
            m++
        }
    }
    matrix := make([][]int, kmernum)
    for i := 0; i < kmernum; i++ {
        matrix[i] = make([]int, len(files))
    }
    for index, file := range files{
        localFreq := make(map[int]int32)
        kmerFreq(file.Name(), localFreq, kmer_len)
        for i := 0; i < kmernum; i++ {
            matrix[i][index] = int(localFreq[topK[i]])
        }
    }
    for i := 0; i < kmernum; i++ {
        head[0] = strconv.Itoa(topK[i])
        for j := 1; j < len(files)+1; j++{
            head[j] = strconv.Itoa(matrix[i][j-1])
        }
        returnError := rw.Write(head)
        if returnError != nil {
            fmt.Println(returnError)
        }
        rw.Flush()
    }
    resultfile.Close()
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