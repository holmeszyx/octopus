package main

import (
	"crypto/cipher"
	"io"
)

const chunkSize = 1024 // 1K

type writer struct {
	io.Writer
	cipher.AEAD
	nonce []byte
	buf   []byte
}

func newWriter(w io.Writer, aead cipher.AEAD) *writer {
	return &writer{
		Writer: w,
		AEAD:   aead,
		nonce:  make([]byte, aead.NonceSize()),
		buf:    make([]byte, chunkSize),
	}
}

func (w *writer) WriterFromRead(r io.Reader) (n int64, err error) {

	maxPayload := chunkSize - w.Overhead()

	for {
		payloadBuf := w.buf[:maxPayload]

		ar, re := r.Read(payloadBuf)

		if ar > 0 {
			n += int64(ar)
			outBuf := w.buf[:ar+w.Overhead()] // mostly , len is chunkSize. if the len < chunkSize, it's the last chunk
			payloadBuf = payloadBuf[:ar]
			w.Seal(outBuf[:0], w.nonce, payloadBuf, nil)
			increment(w.nonce)

			_, we := w.Write(outBuf)
			if we != nil {
				err = we
				break
			}
		}

		if re != nil && re != io.EOF {
			err = re
		}
		if re != nil {
			break
		}

	}

	return n, err
}

type reader struct {
	io.Reader
	cipher.AEAD
	nonce    []byte
	buf      []byte
	plainBuf []byte
}

func newReader(r io.Reader, aead cipher.AEAD) *reader {
	return &reader{
		Reader:   r,
		AEAD:     aead,
		nonce:    make([]byte, aead.NonceSize()),
		buf:      make([]byte, chunkSize),
		plainBuf: make([]byte, chunkSize),
	}
}

// read a segment which is a full buf
func (r *reader) readSegment() (plainSize int, err error) {
	buf := r.buf[:]
	ar, re := r.Read(buf)

	maxPlain := chunkSize
	usedChunk := 0
	if ar > 0 {
		for {
			remain := ar - usedChunk
			if remain <= 0 {
				break
			}

			// fmt.Println(remain, usedChunk)

			plain := r.plainBuf[plainSize:maxPlain]
			limit := chunkSize

			if remain < chunkSize {
				limit = remain
			}
			chunk := buf[usedChunk : usedChunk+limit]
			r.Open(plain[:0], r.nonce, chunk, nil)

			plainSize += (len(chunk) - r.Overhead())
			usedChunk += len(chunk)

			increment(r.nonce)

		}
	}

	if re != nil {
		err = re
	}

	return plainSize, err
}

type FuncSegmentOuter func(data []byte, err error)

func (r *reader) Output(outer FuncSegmentOuter) {
	for {
		n, err := r.readSegment()
		if err != nil && err != io.EOF {
			outer(r.plainBuf, err)
			break
		}
		outer(r.plainBuf[:n], err)
		if err == io.EOF {
			break
		}
	}
}

// something inspire from shadowsocks-go2
func increment(b []byte) {
	for i := range b {
		b[i]++
		if b[i] != 0 {
			return
		}
	}
}
