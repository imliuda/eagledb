package point_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/eagledb/eagledb/point"
)

func BenchmarkParse(b *testing.B) {
	sbuf := "cpu,host=debian,cpu=0 user=123,system=456 153456433"
	for i := 0; i < b.N; i++ {
		_, err := point.Parses(sbuf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func ExampleParse() {
	sbuf := "cpu,host=debian,cpu=0 user=123,system=456 153456433"
	points, err := point.Parses(sbuf)
	if err != nil {
		log.Println(err)
	}
	p := points[0]

	fmt.Println(string(p.Name()))

	for _, tag := range p.Tags() {
		fmt.Println(string(tag.Key), string(tag.Value))
	}

	iter := point.NewFieldIterator(p)
	for iter.Iterate() {
		fmt.Println(string(iter.Key()), iter.Value().(int64))
	}

	fmt.Println(p.Time())
	// Output:
	// cpu
	// cpu 0
	// host debian
	// user 123
	// system 456
	// 153456433
}
