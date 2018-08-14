package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"strings"
	"testing"
)

func Test_all(t *testing.T) {
	plainText := "show me the money!"
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Error(err)
		return
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Error(err)
		return
	}

	buff := bytes.NewBuffer(make([]byte, 0, 64*1024))
	writer := newWriter(buff, aead)
	n, err := writer.WriterFromRead(strings.NewReader(plainText))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("encrpt", n, buff.Len(), "overhead", writer.Overhead())

	reader := newReader(buff, aead)
	rn, err := reader.readSegment()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("decrpt", rn, string(reader.plainBuf[:n]), "overhead", reader.Overhead())

}
