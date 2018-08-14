package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

// encode or decode
var decode = false
var outFile string
var inFile string
var keys = "a2f126e3ea0f141e52b343024f33805e" // showmethemoney10010

func init() {
	flag.BoolVar(&decode, "d", false, "is decoding")
	flag.StringVar(&outFile, "o", "", "the output file")
	flag.Parse()
	if flag.NArg() > 0 {
		inFile = flag.Arg(0)
	}
}

func main() {

	if inFile == "" {
		fmt.Println("no input file")
		os.Exit(1)
	}

	if !isFileExist(inFile) {
		fmt.Println("in file not exits")
		os.Exit(2)
	}

	if outFile == "" {
		fmt.Println("require a out file.")
		os.Exit(3)
	}

	if decode {
		decoding()
	} else {
		encoding()
	}
}

func decoding() {

	outWriter, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
		return
	}
	defer outWriter.Close()

	inReader, err := os.Open(inFile)
	defer inReader.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
		return
	}

	key, _ := hex.DecodeString(keys)
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
		return
	}

	aead, _ := cipher.NewGCM(block)

	bufreader := bufio.NewReader(inReader)
	r := newReader(bufreader, aead)

	cb := func(data []byte, err error) {
		if err != nil && err != io.EOF {
			fmt.Println(err)
			os.Exit(5)
			return
		}
		outWriter.Write(data)
	}

	r.Output(cb)
}

func encoding() {

	outWriter, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
		return
	}
	defer outWriter.Close()

	inReader, err := os.Open(inFile)
	defer inReader.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
		return
	}

	key, _ := hex.DecodeString(keys)
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
		return
	}

	aead, _ := cipher.NewGCM(block)

	w := newWriter(outWriter, aead)
	_, err = w.WriterFromRead(inReader)
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
		return
	}

}

func isFileExist(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil && err != os.ErrExist {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}
