package eagledb

type Config struct {
	ListenAddr string
	ListenPort string
}

type EagleServer struct {
	config *Config
	databases []*Database
	clients []Client
}

func (s *EagleSever) Start() error {

}

func (s *EagleServer) CreateDatabase(name string) error {

}

func main() {
	s := EagleServer{}

	err := e.Start()
	if err != nil {
		fmt.Println("failed to start eagledb:", err)
	}
}
