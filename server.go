package mbslave

import (
	"go.bug.st/serial"
)

type Server struct {
	DataModel DataModel
	Transport Transport
}

func NewRtuServer(config serial.Mode, dataModel DataModel) *Server {
	transport := NewRtuTransport(config)
	return NewServer(transport, dataModel)
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
