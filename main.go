package main

import (
  "os"
  "os/signal"
  "io/ioutil"
  "fmt"
  "net"
  "syscall"
  "strconv"
  "log"
)

func main() {
  log.SetFlags(log.Lshortfile)
  tempDir, err := ioutil.TempDir("", "golang-sample-echo-server.")
  pid := strconv.Itoa(os.Getpid())
  socket := tempDir + "/server." + pid
  listener, err := net.Listen("unix", socket)
  if err != nil {
    log.Printf("error: %v\n", err)
    return
  }
  if err := os.Chmod(socket, 0700); err != nil {
    log.Printf("error: %v\n", err)
    return
  }
  close := make(chan int)
  shutdown(listener, tempDir, close)
  fmt.Printf("GOLANG_SAMPLE_SOCK=%v;export GOLANG_SAMPLE_SOCK;\n", socket)
  fmt.Printf("GOLANG_SAMPLE_PID=%v;export GOLANG_SAMPLE_PID;\n", pid)
  server(listener)
  _ = <-close
}

func shutdown(listener net.Listener, tempDir string, close chan int) {
  c := make(chan os.Signal, 2)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    interrupt := 0
    for {
      s := <-c
      switch s {
      case os.Interrupt:
        if (interrupt == 0) {
          fmt.Println("Interrupt...")
          interrupt++
          continue
        }
      }
      break
    }
    if err := listener.Close(); err != nil {
      log.Printf("error: %v\n", err)
    }
    if err := os.Remove(tempDir); err != nil {
      log.Printf("error: %v\n", err)
    }
    close <- 1
  }()
}

func server(listener net.Listener) {
  for {
    fd, err := listener.Accept()
    if err != nil {
      return
    }
    go process(fd)
  }
}

func process(fd net.Conn) {
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
      log.Printf("error: %v\n", err)
      break
    }
  }
}
