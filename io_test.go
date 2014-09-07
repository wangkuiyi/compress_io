package compressed_io

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

const (
	d = "Hello, world or 你好，世界 or καλημ ́ρα κóσμ or こんにちはせかい\n"
)

func TestNewDecmpressReader_CannotOpenFile(t *testing.T) {
	filename := "a_file_not_there"
	f, e := os.Open(filename)
	r := NewDecompressReader(f, e, path.Ext(filename))
	if r != nil {
		t.Errorf("Expecting nil io.ReadCloser, got non-nil")
	}
}

func TestNewDecmpressReader_ReadPlainData(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte(d))
	if r := NewDecompressReader(ioutil.NopCloser(&buf), nil, ""); r == nil {
		t.Skipf("NewDecompressReader failed")
	} else {
		defer r.Close()
		if b, e := ioutil.ReadAll(r); e != nil {
			t.Errorf("ioutil.ReadAll: %v", e)
		} else {
			if string(b) != d {
				t.Errorf("Expecting %s, got %s", d, string(b))
			}
		}
	}
}

func TestNewDecmpressReader_ReadGzipData(t *testing.T) {
	var buf bytes.Buffer
	if w := gzip.NewWriter(&buf); w == nil {
		t.Skipf("gzip.NewWriter failed")
	} else {
		w.Write([]byte(d))
		w.Close()
	}

	if r := NewDecompressReader(ioutil.NopCloser(&buf), nil, ".gz"); r == nil {
		t.Skipf("NewDecompressedReader failed")
	} else {
		defer r.Close()
		if b, e := ioutil.ReadAll(r); e != nil {
			t.Errorf("ioutil.ReadAll: %v", e)
		} else {
			if string(b) != d {
				t.Errorf("Expecting %s, got %s", d, string(b))
			}
		}
	}
}

func TestNewDecmpressReader_ReadBzip2Data(t *testing.T) {
	log.Println("TODO(wyi): Implemente TestNewDecmpressReader_ReadBzip2Data")
}

func TestNewCompressWriter_CannotCreateFile(t *testing.T) {
	filename := "/tmp/not_exist_dir/not_there_file"
	f, e := os.Create(filename)
	w := NewCompressWriter(f, e, path.Ext(filename))
	if w != nil {
		t.Errorf("Expecting nil io.WriteCloser, got non-nil")
	}
}

func writeAndRead(format string) error {
	var buf bytes.Buffer
	if w := NewCompressWriter(NopWriteCloser(&buf), nil, format); w == nil {
		return errors.New("NewCompressWriter failed")
	} else {
		w.Write([]byte(d))
		w.Close()
	}

	r := NewDecompressReader(ioutil.NopCloser(&buf), nil, format)
	if r == nil {
		return errors.New("NewDecompressReader failed")
	} else {
		defer r.Close()
		if b, e := ioutil.ReadAll(r); e != nil {
			return fmt.Errorf("ioutil.ReadAll: %v", e)
		} else {
			if string(b) != d {
				return fmt.Errorf("Expecting %s, got %s", d, string(b))
			}
		}
	}
	return nil
}

func TestNewCompressWriter_WritePlainData(t *testing.T) {
	if e := writeAndRead(""); e != nil {
		t.Error(e)
	}
}

func TestNewCompressWriter_WriteGzipData(t *testing.T) {
	if e := writeAndRead(".gz"); e != nil {
		t.Error(e)
	}
}

type nopWriteCloser struct {
	io.Writer
}

func (w *nopWriteCloser) Close() error { return nil }

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}
