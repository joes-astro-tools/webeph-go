package zabinski_test

import (
	"testing"
	"webeph/zabinski"
)

func TestSeekLatitude(t *testing.T) {
	// There are 438336 hours between July 15, 2022@00:00 and July 15, 2072@00:00
	max, min, ave, maxT := zabinski.SeekLatitudes(438336)
	t.Log(max)
	t.Log(min)
	t.Log(ave)
	t.Log(maxT)
}
