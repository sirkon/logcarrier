/*
package frameio
Bufferized writer that keeps integrity of blocks produced by data compressors.
This is important, since decomporessors usually work upon the stream of data splitted
into frames, where the frame is a solid chunk that is required to be a single piece to be
decompressed.
*/

package frameio

import (
	"bindec"
	"binenc"
	"bytes"
	"io"
)

const (
	defaultBufferSize = 128 * 1024 * 1024
)

// Writer is a bufferized writer which takes care of comporession frames integrity, i.e.
// it only splits data at frame bounds, not at frame internals.
type Writer struct {
	bufsize int
	writer  io.Writer
	buffer  *bytes.Buffer

	flushCounter     int
	frameInsert      int
	prevFlushCounter int

	worthFlushing bool
}

// NewWriter constructs writer whose buffer has the default size
func NewWriter(writer io.Writer) *Writer {
	return NewWriterSize(writer, defaultBufferSize)
}

// NewWriterSize return a new writer whose buffer has at least specified size
func NewWriterSize(writer io.Writer, size int) *Writer {
	res := &Writer{
		bufsize:       size,
		writer:        writer,
		buffer:        &bytes.Buffer{},
		worthFlushing: true,
	}
	res.buffer.Grow(size)
	return res
}

// Flush flushes all buffered data
func (w *Writer) Flush() error {
	if w.buffer.Len() > 0 {
		w.flushCounter = w.frameInsert
		if _, err := w.buffer.WriteTo(w.writer); err != nil {
			return err
		}
	}
	w.buffer.Reset()
	return nil
}

// Write writes the content of data into the buffer
func (w *Writer) Write(data []byte) (nn int, err error) {
	if w.buffer.Len() > 0 && w.buffer.Len()+len(data) > w.bufsize {
		w.worthFlushing = false
		err = w.Flush()
		if err != nil {
			return
		}
	}
	if len(data) > w.bufsize {
		nn, err = w.writer.Write(data)
		return
	}
	w.frameInsert++
	nn, err = w.buffer.Write(data)
	return
}

// WorthFlushing checks if any write was done after the last check
func (w *Writer) WorthFlushing() bool {
	res := w.worthFlushing && w.frameInsert != w.prevFlushCounter && w.prevFlushCounter == w.flushCounter
	w.prevFlushCounter = w.flushCounter
	w.worthFlushing = true
	return res
}

// DumpState ...
func (w *Writer) DumpState(enc *binenc.Encoder, dest *bytes.Buffer) {
	dest.Write(enc.Uint32(uint32(w.bufsize)))
	dest.Write(enc.Uint32(uint32(w.buffer.Len())))
	dest.Write(enc.Uint32(uint32(w.flushCounter)))
	dest.Write(enc.Uint32(uint32(w.frameInsert)))
	dest.Write(enc.Uint32(uint32(w.prevFlushCounter)))
	dest.Write(enc.Bool(w.worthFlushing))
}

// RestoreState ...
func (w *Writer) RestoreState(src *bindec.Decoder) {
	bufsize, ok := src.Uint32()
	if !ok {
		panic("Cannot restore bufsize")
	}
	buflen, ok := src.Uint32()
	if !ok {
		panic("Cannot restore buffer length")
	}
	buffer, ok := src.Bytes(int(buflen))
	if !ok {
		panic("Cannot restore a buffer")
	}
	flushcounter, ok := src.Uint32()
	if !ok {
		panic("Cannot restore flush counter")
	}
	frameInsert, ok := src.Uint32()
	if !ok {
		panic("Cannot restore frame insert")
	}
	prevFlushCounter, ok := src.Uint32()
	if !ok {
		panic("Cannot restore previous flush counter value")
	}
	worthFlushing, ok := src.Bool()
	if !ok {
		panic("Cannot restore worth flushing indicator")
	}
	w.bufsize = int(bufsize)
	w.buffer.Reset()
	w.buffer.Write(buffer)
	w.flushCounter = int(flushcounter)
	w.frameInsert = int(frameInsert)
	w.prevFlushCounter = int(prevFlushCounter)
	w.worthFlushing = worthFlushing
}
