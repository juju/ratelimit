// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package ratelimit

import (
	gc "launchpad.net/gocheck"

	"testing"
	"time"
)

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

type rateLimitSuite struct{}

var _ = gc.Suite(rateLimitSuite{})

type takeReq struct {
	time       time.Duration
	count      int64
	expectWait time.Duration
}

var takeTests = []struct {
	about        string
	fillInterval time.Duration
	capacity     int64
	reqs         []takeReq
}{{
	about:        "serial requests",
	fillInterval: 250 * time.Millisecond,
	capacity:     10,
	reqs: []takeReq{{
		time:       0,
		count:      0,
		expectWait: 0,
	}, {
		time:       0,
		count:      10,
		expectWait: 0,
	}, {
		time:       0,
		count:      1,
		expectWait: 250 * time.Millisecond,
	}, {
		time:       250 * time.Millisecond,
		count:      1,
		expectWait: 250 * time.Millisecond,
	}},
}, {
	about:        "concurrent requests",
	fillInterval: 250 * time.Millisecond,
	capacity:     10,
	reqs: []takeReq{{
		time:       0,
		count:      10,
		expectWait: 0,
	}, {
		time:       0,
		count:      2,
		expectWait: 500 * time.Millisecond,
	}, {
		time:       0,
		count:      2,
		expectWait: 1000 * time.Millisecond,
	}, {
		time:       0,
		count:      1,
		expectWait: 1250 * time.Millisecond,
	}},
}, {
	about:        "more than capacity",
	fillInterval: 1 * time.Millisecond,
	capacity:     10,
	reqs: []takeReq{{
		time:       0,
		count:      10,
		expectWait: 0,
	}, {
		time:       20 * time.Millisecond,
		count:      15,
		expectWait: 5 * time.Millisecond,
	}},
}, {
	about:        "sub-quantum time",
	fillInterval: 10 * time.Millisecond,
	capacity:     10,
	reqs: []takeReq{{
		time:       0,
		count:      10,
		expectWait: 0,
	}, {
		time:       7 * time.Millisecond,
		count:      1,
		expectWait: 3 * time.Millisecond,
	}, {
		time:       8 * time.Millisecond,
		count:      1,
		expectWait: 12 * time.Millisecond,
	}},
}, {
	about:        "within capacity",
	fillInterval: 10 * time.Millisecond,
	capacity:     5,
	reqs: []takeReq{{
		time:       0,
		count:      5,
		expectWait: 0,
	}, {
		time:       60 * time.Millisecond,
		count:      5,
		expectWait: 0,
	}, {
		time:       60 * time.Millisecond,
		count:      1,
		expectWait: 10 * time.Millisecond,
	}, {
		time:       80 * time.Millisecond,
		count:      2,
		expectWait: 10 * time.Millisecond,
	}},
}}

func (rateLimitSuite) TestTake(c *gc.C) {
	for i, test := range takeTests {
		tb := New(test.fillInterval, test.capacity)
		for j, req := range test.reqs {
			d := tb.take(tb.startTime.Add(req.time), req.count)
			if d != req.expectWait {
				c.Fatalf("test %d.%d, %s, got %v want %v", i, j, test.about, d, req.expectWait)
			}
		}
	}
}

type tryTakeReq struct {
	time   time.Duration
	count  int64
	expect int64
}

var tryTakeTests = []struct {
	about        string
	fillInterval time.Duration
	capacity     int64
	reqs         []tryTakeReq
}{{
	about:        "serial requests",
	fillInterval: 250 * time.Millisecond,
	capacity:     10,
	reqs: []tryTakeReq{{
		time:   0,
		count:  0,
		expect: 0,
	}, {
		time:   0,
		count:  10,
		expect: 10,
	}, {
		time:   0,
		count:  1,
		expect: 0,
	}, {
		time:   250 * time.Millisecond,
		count:  1,
		expect: 1,
	}},
}, {
	about:        "concurrent requests",
	fillInterval: 250 * time.Millisecond,
	capacity:     10,
	reqs: []tryTakeReq{{
		time:   0,
		count:  5,
		expect: 5,
	}, {
		time:   0,
		count:  2,
		expect: 2,
	}, {
		time:   0,
		count:  5,
		expect: 3,
	}, {
		time:   0,
		count:  1,
		expect: 0,
	}},
}, {
	about:        "more than capacity",
	fillInterval: 1 * time.Millisecond,
	capacity:     10,
	reqs: []tryTakeReq{{
		time:   0,
		count:  10,
		expect: 10,
	}, {
		time:   20 * time.Millisecond,
		count:  15,
		expect: 10,
	}},
}, {
	about:        "within capacity",
	fillInterval: 10 * time.Millisecond,
	capacity:     5,
	reqs: []tryTakeReq{{
		time:   0,
		count:  5,
		expect: 5,
	}, {
		time:   60 * time.Millisecond,
		count:  5,
		expect: 5,
	}, {
		time:   70 * time.Millisecond,
		count:  1,
		expect: 1,
	}},
}}

func (rateLimitSuite) TestTryTake(c *gc.C) {
	for i, test := range tryTakeTests {
		tb := New(test.fillInterval, test.capacity)
		for j, req := range test.reqs {
			d := tb.tryTake(tb.startTime.Add(req.time), req.count)
			if d != req.expect {
				c.Fatalf("test %d.%d, %s, got %v want %v", i, j, test.about, d, req.expect)
			}
		}
	}
}

func (rateLimitSuite) TestPanics(c *gc.C) {
	c.Assert(func() { New(0, 1) }, gc.PanicMatches, "token bucket fill interval is not > 0")
	c.Assert(func() { New(-2, 1) }, gc.PanicMatches, "token bucket fill interval is not > 0")
	c.Assert(func() { New(1, 0) }, gc.PanicMatches, "token bucket capacity is not > 0")
	c.Assert(func() { New(1, -2) }, gc.PanicMatches, "token bucket capacity is not > 0")
}
