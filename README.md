# Compress IO

Go standard libraries support reading/writing Gzip and Bzip2 data
streams.  However they do not have common interface.  For example,
`compress/bzip2.NewReader` returns `io.Reader`, where as
`compress/gzip.NewReader` returns `io.ReadCloser`.  This package
encapsulates the `compress` package and provide unified interface.


## Read

Read from a file while decompressing:

     f, e := os.Open(filename)
     if r := compress_io.NewReader(f, e, path.Ext(filename)); r != nil {
       defer r.Close()
       ... read from r ...
     }

If filename has extension ".bz2", `r` would decompress bzip2 data
stream, or if it is ".gz", then `r` would decompress gzip data stream,
otherwise, `r` does not invoke any decompression algorithm.

## Write

Write to a file while compressing:

      f, e := os.Create(filename)
      if w := compress_io.NewWriter(f, e, path.Ext(filename)); w != nil {
        defer w.Close()
        ... write to w ...
      }

If filename has extension ".bz2", `w` would compress bzip2 data
stream, or if it is ".gz", then `r` would compress gzip data stream,
otherwise, `r` does not invoke any compression algorithm.

## Caveat

Unfortunately, Go does not support taking multiple return values as
parameters, as described in this
[discussion](https://code.google.com/p/go/issues/detail?id=973).  This
prevents us from writing:

    if w := compress_io.NewCompressWriter(os.Create(filename), path.Ext(filename))

Instead, we have to write two lines:

    f, e := os.Create(filename)
    w := compress_io.NewWriter(f, e, path.Ext(filename))

So, there is a possibility to write:

    defer f.Close()

instead of

    defer w.Close()

However, `defer w.Close()` is what we want.
