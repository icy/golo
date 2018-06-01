/*
  Purpose : A Golang version of rolo
            (rolo is a Ruby version of solo)
            (solo is Perl program written by Tim Kay)
  Author  : Ky-Anh Huynh
  Github  : https://github.com/icy/golo
  Date    : 2018 June 1st
  License : MIT
*/

package main

import (
  "fmt"
  "flag"
  "net"
  "os"
)

func isPortAvailable(ip string, port int, timeout int) bool {
  conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port));
  if err != nil {
    return false
  }
  conn.Close();
  return true;
}

func warn(msg string) {
  fmt.Fprintf(os.Stderr, fmt.Sprintf(":: %s", msg));
}

func main() {
  ipAddress := flag.String("address", "127.0.0.1", "Address to listen on or to check");
  workDir   := flag.String("dir", ".", "Working diretory");
  port      := flag.Int("port", 0, "Port to listen on or to check");
  timeout   := flag.Int("timeout", 1, "Timeout when checking. Default: 1 second.");
  noBind    := flag.Bool("no-bind", false, "Do not bind on address:port specified");
  flag.Parse();

  var conn = new(net.Listener);

  if *noBind {
    if isPortAvailable(*ipAddress, *port, *timeout) {
      warn("Port is available. App is not running\n");
    } else {
      warn("Port is not available. App is running?\n");
      os.Exit(0);
    }
  } else {
    conn_, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ipAddress, *port));
    if err != nil {
      warn(fmt.Sprintf("Unable to bind on %s:%d. App is running?\n", *ipAddress, *port));
      os.Exit(1);
    }
    *conn = conn_;
  }

  if *conn != nil {
    (*conn).Close();
  }

  warn(fmt.Sprintf("Now staring application from %s\n", *workDir));
}
