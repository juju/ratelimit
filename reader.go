// +build ignore

type Reader struct {
	r io.Reader
	flow byteFlow
}


// NewReader returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte. Note that means this reader
// cannot be used for rate limiting of more than
// 1GB/s, as that would require a fill interval of
// less than 1 ns. Also, very high flow rates will
// be inaccurate
func NewReader(r io.Reader, bucket *TokenBucket, quantum int) io.Reader {
	
}

func (r *Reader) Read(buf []byte) (int, error) {
	n, err := r.Read(buf)
	if n <= 0 {
		return n, err
	}
	r.flow.wait(n)
	return n, err
}

type byteFlow struct {
	bucket *tokenBucket
	borrowed int
	quantum int
}

func (f *byteFlow) wait(n int) {
	if n < r.borrowed {
		r.borrowed -= n
		return
	}
	// Pay back any bytes we borrowed earlier.
	n -= r.borrowed
	// Calculate number of tokens we need, rounding
	// up if necessary.
	tokens := (n + r.quantum - 1) / r.quantum
	r.bucket.Wait(tokens)

	// If we took more than we strictly need,
	// pay it back later.
	r.borrowed = tokens * r.quantum - len(buf)
}

type Writer struct {
	w io.Writer
	flow byteFlow
}


func NewWriter(w io.Writer, bucket *TokenBucket, quantum int) io.Reader {
	
}

func (w *Writer) Write(buf []byte) (int, error) {
	w.flow.wait(len(buf))
	return w.w.Write(buf)
}
