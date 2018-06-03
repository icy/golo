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
  "syscall"
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

    warn(fmt.Sprintf("Bind successfully on %s:%d\n", *ipAddress, *port));
    *conn = conn_;
  }

  err := syscall.Chdir(*workDir);
  if err != nil {
    warn(fmt.Sprintf("Switching to '%s' got error '%s'\n", *workDir, err));
    os.Exit(1);
  }

  cmdArgs   := flag.Args();
  if len(cmdArgs) < 1 {
    warn(fmt.Sprintf("You must specify a command\n"));
    os.Exit(1);
  }
  execPath := cmdArgs[0];
  warn(fmt.Sprintf("Now staring application '%s' from %s\n", execPath, *workDir));
  err = syscall.Exec(execPath, cmdArgs, syscall.Environ());
  if err != nil {
    warn(fmt.Sprintf("Executing got error '%s'\n", err));
    os.Exit(1);
  }
}
