# ratelimit
--
    import "github.com/juju/ratelimit"

The ratelimit package provides an efficient token bucket implementation. See
http://en.wikipedia.org/wiki/Token_bucket.

## Usage

#### func  Reader

```go
func Reader(r io.Reader, bucket *TokenBucket) io.Reader
```
Reader returns a reader that is rate limited by the given token bucket. Each
token in the bucket represents one byte.

#### func  Writer

```go
func Writer(w io.Writer, bucket *TokenBucket) io.Writer
```
Writer returns a reader that is rate limited by the given token bucket. Each
token in the bucket represents one byte.

#### type TokenBucket

```go
type TokenBucket struct {
}
```

TokenBucket represents a token bucket that fills at a predetermined rate.
Methods on TokenBucket may be called concurrently.

#### func  New

```go
func New(fillInterval time.Duration, capacity int64) *TokenBucket
```
New returns a new token bucket that fills at the rate of one token every
fillInterval, up to the given maximum capacity. Both arguments must be positive.
The bucket is initially full.

#### func  NewWithRate

```go
func NewWithRate(rate float64, capacity int64) *TokenBucket
```
NewRate returns a token bucket that fills the bucket at the rate of rate tokens
per second up to the given maximum capacity. Because of limited clock
resolution, at high rates, the actual rate may be up to 1% different from the
specified rate.

#### func (*TokenBucket) Rate

```go
func (tb *TokenBucket) Rate() float64
```
Rate returns the fill rate of the bucket, in tokens per second.

#### func (*TokenBucket) Take

```go
func (tb *TokenBucket) Take(count int64) time.Duration
```
Take takes count tokens from the bucket without blocking. It returns the time
that the caller should wait until the tokens are actually available.

Note that if the request is irrevocable - there is no way to return tokens to
the bucket once this method commits us to taking them.

#### func (*TokenBucket) TakeAvailable

```go
func (tb *TokenBucket) TakeAvailable(count int64) int64
```
TakeAvailable takes up to count immediately available tokens from the bucket. It
returns the number of tokens removed, or zero if there are no available tokens.
It does not block.

#### func (*TokenBucket) Wait

```go
func (tb *TokenBucket) Wait(count int64)
```
Wait takes count tokens from the bucket, waiting until they are available.
