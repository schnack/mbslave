package mbslave

type Server struct {
	DataModel DataModel
	Transport Transport
}

func NewRtuServer(config *Config) *Server {
	transport := NewRtuTransport(config)
	return NewServer(transport, NewDefaultDataModel(config))
}

func NewServer(transport Transport, dataModel DataModel) *Server {
	transport.SetHandler(dataModel.Handler)
	return &Server{
		DataModel: dataModel,
		Transport: transport,
	}
}

func (s *Server) Listen() error {
	return s.Transport.Listen()
}
