package main

import (
  "net"
  "log"
  "os"
  "fmt"
)

type Server struct {
  listener net.Listener
}

func NewServer() *Server{
  s := new(Server)
  return s;
}

func (s *Server) Open(socket string) {
  listener, err := net.Listen("unix", socket)
  if err != nil {
    return err
  }
  s.listener = listener;
  if err := os.Chmod(socket, 0600); err != nil {
    s.Close()
    return err
  }
}

func (s *Server) Close() {
  if err := s.listener.Close(); err != nil {
    return err;
  }
}

func (s *Server) Start() {
  for {
    fd, err := s.listener.Accept()
    if err != nil {
      return
    }
    go s.Process(fd)
  }
}

func (s *Server) Process(fd net.Conn) {
  defer fd.Close()
  for {
    buf := make([]byte, 512)
    nr, err := fd.Read(buf)
    if err != nil {
      break
    }
    data := buf[0:nr]
    fmt.Printf("Recieved: %v", string(data));
    _, err = fd.Write(data)
    if err != nil {
      return err
    }
  }
}
