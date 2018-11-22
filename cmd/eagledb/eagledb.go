package eagledb

import (
	"github.com/eagledb/eagledb"
)

func main() {
	server := eagledb.Server{}

	err := server.Start()
	if err != nil {
		fmt.Println("failed to start eagledb:", err)
	}
}
