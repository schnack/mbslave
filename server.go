package mbslave

import "github.com/goburrow/serial"

type Server struct {
	DataModel DataModel
	transport Transport
}

func NewRtuServer(config *serial.Config, dataModel DataModel) *Server {
	return &Server{
		DataModel: dataModel,
		transport: NewRtuTransport(config, dataModel.Handler),
	}
}

func (s *Server) Listen() error {
	s.DataModel.Init()
	return s.transport.Listen()
}
