/* 
	Copyright 2016 Shanshan Gao
	Record the position of each genome in concatenate one
*/

package main

import (
	"os"
	"fmt"
	"bufio"
	"bytes"
)

func main() {
	file := os.Args[1]

	f, err := os.Open(file)
	check_for_error(err)
	defer f.Close()

	if file[len(file)-6:] != ".fasta" {
		panic("ReadFasta:" + file + "is not a fasta file.")
	}

	position := make([]int, 0)

	scanner := bufio.NewScanner(f)
	byte_array := make([]byte, 0)
	i := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > 0 {
			if line[0] == '>' {
				fmt.Println(string(line))
			}
			if line[0] != '>' {
				byte_array = append(byte_array, bytes.Trim(line, "\n\r ")...)
			} else if len(byte_array) > 0 {
				byte_array = append(byte_array, byte('|'))
				position = append(position, len(byte_array) - 1)
			}
			i++
		}
	}

	position = append(position, len(byte_array))
	begin := 0
	for _, p := range position {
		fmt.Println(p - begin)
		begin = p
	}
}

func check_for_error(e error) {
	if e != nil {
		panic(e)
	}
}