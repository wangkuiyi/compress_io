package compressed_io

import (
	"compress/bzip2"
	"compress/gzip"
	"io"
	"log"
)

// NewDecompressReader accepts an input io.ReadCloser and an error,
// checks if they are valid input reader, and returns an io.ReadCloser
// which decompress content read from the input.  The decompressing
// algorithm is specified by format, which could be "" for no
// decompressing, ".bz" for bzip2, and ".gz" for gzip.
//
// Example:
/*
   f, e := os.Open(filename)
   if r := NewDecompressReader(f, e, path.Ext(filename)); r != nil {
     defer r.Close()
     ... read from r ...
   }
*/
func NewDecompressReader(in io.ReadCloser, e error,
	format string) io.ReadCloser {
	if e != nil || in == nil {
		log.Printf("NewDecompressReader: %v", e)
		return nil
	}

	if len(format) > 0 {
		switch {
		case format == ".bz2":
			return &bzip2ReadCloser{bzip2.NewReader(in), in}
		case format == ".gz":
			if r, e := gzip.NewReader(in); e != nil {
				log.Printf("Cannot create gzip.Reader for %s: %v", format, e)
				return nil
			} else {
				return r
			}
		default:
			log.Printf("Unknown format: %s", format)
			return nil
		}
	}
	return in
}

// NewCompressWriter accepts an output io.WriteCloser and an error,
// checks if they are valid, and returns an io.WriteCloser which
// compress content written into it.  The compressing algorithm is
// specified by format, which could be "" for no compressing, ".gz"
// for gzip.
//
// It is noticable that Go standard library compress/bzip2 does not
// support compressing yet, so neither do we.
//
// Exmaple:
/*
   f, e := os.Create(filename)
   if w := NewCompressWriter(f, e, path.Ext(filename)); w != nil {
     defer w.Close()
     ... write to w ...
   }
*/
func NewCompressWriter(out io.WriteCloser, e error,
	format string) io.WriteCloser {
	if e != nil || out == nil {
		log.Printf("NewCompressWriter: %v", e)
		return nil
	}

	if len(format) > 0 {
		switch {
		case format == ".bz2":
			log.Print("Does not yet support bzip2 writer.")
			return nil
		case format == ".gz":
			if w := gzip.NewWriter(out); w == nil {
				log.Printf("Cannot create gzip writer for %s: %v", format, e)
			} else {
				return w
			}
		default:
			log.Printf("Unknown format: %s", format)
			return nil
		}
	}
	return out
}

type bzip2ReadCloser struct {
	io.Reader
	in io.ReadCloser
}

func (b *bzip2ReadCloser) Close() error {
	return b.in.Close()
}
