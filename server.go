package mbslave

import (
	"github.com/goburrow/serial"
)

type Server struct {
	DataModel DataModel
	Transport Transport
}

func NewRtuServer(config serial.Config, dataModel DataModel) *Server {
	return &Server{
		DataModel: dataModel,
		Transport: NewRtuTransport(config, dataModel.Handler),
	}
}

func NewServer(transport Transport, dataModel DataModel) *Server {
	return &Server{
		DataModel: dataModel,
		Transport: transport,
	}
}

func (s *Server) Listen() error {
	s.DataModel.Init()
	return s.Transport.Listen()
}
