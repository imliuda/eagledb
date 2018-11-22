package point_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/eagledb/eagledb/point"
)

var sbuf = "cpu\\ us\\,age,h\\ o\\=st=d\\,e\\ bian,cpu=0 u\\=s\\ er=123,system=456 153456433\n" +
	"disk_io,host=debian,disk=sda write=123.456,read=45.63 153456433 \n" +
	"ip_addr,host=debian,if=eth0 ip=\"192.168.1.1\\n\\tqwe\" 153456433\n " +
	"service,host=debian,service=apache2 up=true,down=false 153456433 \n " +
	"null,host=debian value=null 153456433"

// Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz
// BenchmarkParse-4   	  100000	     11789 ns/op
// It's about 40000 point per second
// Need Optimise
func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := point.Parses(sbuf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func ExampleParse() {
	points, err := point.Parses(sbuf)
	if err != nil {
		log.Println(err)
	}

	for _, p := range points {
		log.Println(p.String())
		fmt.Printf("%s ", p.Name())

		for _, tag := range p.Tags() {
			fmt.Printf("%s %s ", tag.Key, tag.Value)
		}

		iter := point.NewFieldIterator(p)
		for iter.Next() {
			if iter.Type() == point.String {
				fmt.Printf("%s %s ", iter.Key(), iter.Value())
			} else {
				fmt.Printf("%s %v ", iter.Key(), iter.Value())
			}
		}

		fmt.Printf("%d\n", p.Time())
	}

	// Output:
	// cpu us,age cpu 0 h o=st d,e bian u=s er 123 system 456 153456433
	// disk_io disk sda host debian write 123.456 read 45.63 153456433
	// ip_addr host debian if eth0 ip 192.168.1.1
	// 	qwe 153456433
	// service host debian service apache2 up true down false 153456433
	// null host debian value <nil> 153456433
}
