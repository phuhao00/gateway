package main

//Server gateway server
type Server struct {
	cliServer *ClientServer
	inServer  *InnerServer
}

func NewServer() *Server {
	return &Server{
		cliServer: NewClientServer(),
		inServer:  NewInnerServer(),
	}
}

func (s *Server) Start() {

}
