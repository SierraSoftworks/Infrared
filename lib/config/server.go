package config

type Server struct {
	ListenOn string `json:"listen"`
	Database struct {
		Hosts    string `json:"url"`
		Database string `json:"db"`
	} `json:"database"`
}

func (s *Server) Save(filename string) error {
	return SaveJson(filename, s)
}

func (s *Server) Load(filename string) error {
	return LoadJson(filename, s)
}

func (s *Server) Log() {
	LogJson(s)
}
