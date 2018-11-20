package point

import (
	"testing"
	"github.com/eagledb/eagledb/point"
)

func TestParse(t *testing.T) {
	sbuf := "cpu,host=debian,cpu=0 user=123,system=456 153456433"
	points, err := point.Parses(sbuf)
	if err != nil {
		t.Fatal(err)
	}
	p := points[0]
	t.Log(*p)
}
