package util

import "testing"

func TestSplit(t *testing.T) {
	l := Split(147888,60)
	t.Log(l)
	s := 0
	for _,v := range l {
		s+=v
	}
	t.Log(s)
}
