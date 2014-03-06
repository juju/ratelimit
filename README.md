# ratelimit
--
    import "github.com/juju/ratelimit"

The ratelimit package provides an efficient token bucket implementation. See
http://en.wikipedia.org/wiki/Token_bucket.

## Usage

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

#### func (*TokenBucket) Get

```go
func (tb *TokenBucket) Get(count int64)
```
Get gets count tokens from the bucket, waiting until the tokens are available.

#### func (*TokenBucket) GetNB

```go
func (tb *TokenBucket) GetNB(count int64) time.Duration
```
GetNB gets count tokens from the bucket without blocking. It returns the time to
wait until the tokens are actually available.

Note that if the request is irrevocable - there is no way to return tokens to
the bucket once this method commits us to taking them.
