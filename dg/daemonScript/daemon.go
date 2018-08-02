package daemonScript

type Service struct {
	path        string
	syncChannel chan int64
}

func New(path string, syncChannel chan int64) *Service {
	return &Service{path: path, syncChannel: syncChannel}
}
