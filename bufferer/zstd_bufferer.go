package bufferer

import (
	"bytes"

	"sync"

	"github.com/sirkon/logcarrier/bindec"
	"github.com/sirkon/logcarrier/binenc"
	"github.com/sirkon/logcarrier/fileio"
	"github.com/sirkon/logcarrier/frameio"
	"github.com/sirkon/logcarrier/logio"
)

// ZSTDBufferer ...
type ZSTDBufferer struct {
	l *logio.Writer
	c *ZSTDWriter
	f *frameio.Writer
	d *fileio.File
	p *sync.Pool
}

// NewZSTDBufferer constructor
func NewZSTDBufferer(l *logio.Writer, c *ZSTDWriter, f *frameio.Writer, d *fileio.File) *ZSTDBufferer {
	res := &ZSTDBufferer{
		l: l,
		c: c,
		f: f,
		d: d,
	}
	return res
}

// Write implementation
func (b *ZSTDBufferer) Write(p []byte) (n int, err error) {
	return b.l.Write(p)
}

// PostWrite implementation
func (b *ZSTDBufferer) PostWrite() error {
	if b.l.OvergrownExtra(nil) {
		return b.l.Flush()
	}
	return nil
}

// Close implementation
func (b *ZSTDBufferer) Close() error {
	if err := b.l.Flush(); err != nil {
		return err
	}
	if err := b.c.Close(); err != nil {
		return err
	}
	b.c.Reset()
	if err := b.f.Flush(); err != nil {
		return err
	}
	if err := b.d.Close(); err != nil {
		return err
	}
	return nil
}

// Flush implementation
func (b *ZSTDBufferer) Flush() error {
	if b.l.WorthFlushing() {
		if err := b.l.Flush(); err != nil {
			return err
		}
	}
	if b.f.WorthFlushing() {
		if err := b.c.Close(); err != nil {
			return err
		}
		b.c.Reset()
		if err := b.f.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// Logrotate implementation
func (b *ZSTDBufferer) Logrotate(dir, name, group string) error {
	return b.d.Logrotate(dir, name, group)

}

// DumpState implementation
func (b *ZSTDBufferer) DumpState(enc *binenc.Encoder, dest *bytes.Buffer) error {
	b.l.DumpState(enc, dest)
	b.c.w.Backup()
	b.f.DumpState(enc, dest)
	if err := b.d.DumpState(enc, dest); err != nil {
		return err
	}
	return nil
}

// RestoreState implementation
func (b *ZSTDBufferer) RestoreState(src *bindec.Decoder) error {
	b.l.RestoreState(src)
	b.c.w.Restore()
	b.f.RestoreState(src)
	if err := b.d.RestoreState(src); err != nil {
		return err
	}
	return nil
}
