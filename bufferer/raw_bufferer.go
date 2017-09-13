package bufferer

import (
	"bytes"

	"github.com/sirkon/logcarrier/bindec"
	"github.com/sirkon/logcarrier/binenc"
	"github.com/sirkon/logcarrier/fileio"
	"github.com/sirkon/logcarrier/logio"
)

// RawBufferer ...
type RawBufferer struct {
	l *logio.Writer
	d *fileio.File
}

// NewRawBufferer constructor
func NewRawBufferer(l *logio.Writer, d *fileio.File) *RawBufferer {
	return &RawBufferer{
		l: l,
		d: d,
	}
}

// Write implementation
func (b *RawBufferer) Write(p []byte) (n int, err error) {
	return b.l.Write(p)
}

// PostWrite implementation
func (b *RawBufferer) PostWrite() error {
	return b.l.Flush()
}

// Close implementation
func (b *RawBufferer) Close() error {
	if err := b.l.Flush(); err != nil {
		return err
	}
	if err := b.d.Close(); err != nil {
		return err
	}
	return nil
}

// Flush implementation
func (b *RawBufferer) Flush() error {
	if b.l.WorthFlushing() {
		if err := b.l.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// Logrotate implementation
func (b *RawBufferer) Logrotate(dir, name, group string) error {
	return b.d.Logrotate(dir, name, group)
}

// DumpState implementation
func (b *RawBufferer) DumpState(enc *binenc.Encoder, dest *bytes.Buffer) error {
	b.l.DumpState(enc, dest)
	if err := b.d.DumpState(enc, dest); err != nil {
		return err
	}
	return nil
}

// RestoreState implementation
func (b *RawBufferer) RestoreState(src *bindec.Decoder) error {
	b.l.RestoreState(src)
	if err := b.d.RestoreState(src); err != nil {
		return err
	}
	return nil
}
