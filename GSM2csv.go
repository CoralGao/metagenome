package main

import (
   "os"
   "fmt"
   "bufio"
   "bytes"
   "sync"
   "runtime"
   "strconv"
   "time"
   "io/ioutil"
   "math"
   "encoding/csv"
)

func main() {
    if len(os.Args) != 5 {
        panic("must provide sequence folder file, readfile and result file names.")
    }

    // Count GSM for read file
    resultRead := make(chan int, 10)
    go CountFreq(os.Args[2], 14, resultRead)
    gsmread := make(map[int]int)
    for res := range resultRead {
        gsmread[res] = gsmread[res] + 1
    }

    files, _ := ioutil.ReadDir(os.Args[1])
    start_time := time.Now()
    gsm := make(map[int]int)
    gsmFreq := make(map[int]int)

    core_num := 1
    kmer_len := 7
    distance := 0
    runtime.GOMAXPROCS(core_num+2)

    // Build a csv file to store the GSM from genomes&reads
    resultfile, err1 := os.Create(os.Args[3]+".csv")
    if err1 != nil {
        fmt.Printf("%v\n",err1)
        os.Exit(1)
    }

    resultbfile, err2 := os.Create(os.Args[4]+".csv")
    if err2 != nil {
        fmt.Printf("%v\n",err2)
        os.Exit(1)
    }

    rw := csv.NewWriter(resultfile)
    rwb := csv.NewWriter(resultbfile)
    head := make([]string, len(files)+1)
    head[0] = "kmer"

    for index, fi := range files {
        head[index+1] = fi.Name()
    }
    headb := make([]string, 2)
    headb[0] = "kmer"
    headb[1] = "b"

    returnError := rw.Write(head)
    if returnError != nil {
        fmt.Println(returnError)
    }

    returnErrorb := rwb.Write(headb)
    if returnErrorb != nil {
        fmt.Println(returnErrorb)
    }
    rw.Flush()
    rwb.Flush()

    // Get the unique GSM frequency from genome files
    for index, fi := range files {
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

        // GSM info from the current genome
        gsm1 := make(map[int]int)
        gsmFreq1 := make(map[int]int)

        for res := range result {
            gsm1[res] = index + 1
            gsmFreq1[res] = gsmFreq1[res] + 1
        }

        for k := range gsm1 {
            if gsm[k] == 0 {
                gsm[k] = gsm1[k]
                gsmFreq[k] = gsmFreq1[k]
            } else {
                gsm[k] = -1
                gsmFreq[k] = 0
            }
        }

        f.Close()

    }

    // Merge the unique GSM from genome files & reads to csv file
    for k := range gsm {
        if gsm[k] != -1 {
            line := make([]string, len(files)+1)
            lineb := make([]string, 2)
            for i := range line {
                if i == 0 {
                    line[0] = strconv.Itoa(k)
                    lineb[0] = strconv.Itoa(k)
                } else if i == gsm[k] {
                    line[gsm[k]] = strconv.Itoa(gsmFreq[k])
                } else {
                    line[i] = strconv.Itoa(0)
                } 
            }
            lineb[1] = strconv.Itoa(gsmread[k])
            returnError := rw.Write(line)
            returnErrorb := rwb.Write(lineb)
            if returnError != nil {
                fmt.Println(returnError)
            }
            if returnErrorb != nil {
                fmt.Println(returnErrorb)
            }
        }
    }
    rw.Flush()
    rwb.Flush()

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