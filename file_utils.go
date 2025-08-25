package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	BUFF_SIZE       = 10000
	MAX_WORD_LENGTH = 500
)

func ReadFromFile(filename string, ch chan [][]byte) {
	defer close(ch)
	// build process
	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// reader:
	r := bufio.NewReader(file)

	// buffer alloc
	buffer := make([][]byte, BUFF_SIZE)
	for i := range len(buffer) {
		buffer[i] = make([]byte, MAX_WORD_LENGTH)
	}

	// bufferized file reading
	processedCount := 0
	err = BufferizedReading(r, buffer, BUFF_SIZE, ch)
	for err != io.EOF {
		processedCount += BUFF_SIZE

		if processedCount%50000 == 0 && processedCount > 0 {
			fmt.Printf("Read %d lines from file %s\n", processedCount, filename)
		}

		err = BufferizedReading(r, buffer, BUFF_SIZE, ch)
		if err != nil {
			break
		}
	}
}

func BufferizedReading(r *bufio.Reader, buffer [][]byte, max int, ch chan [][]byte) error {
	var err error

	for i := range max {
		var line []byte
		line, err = ReadLine(r)

		if len(line) > 0 {
			buffer[i] = line
		}

		if err != nil {
			return err
		}
	}
	ch <- buffer

	return nil
}

func ReadLine(r *bufio.Reader) ([]byte, error) {
	var b byte
	var err error
	var word []byte

	for b != '\n' {
		b, err = r.ReadByte()
		if err != nil {
			return word, err
		}
		word = append(word, b)
	}

	return word, nil
}
