/*
  Purpose : A Golang version of rolo
            (rolo is a Ruby version of solo)
            (solo is Perl program written by Tim Kay)
  Author  : Ky-Anh Huynh
  Github  : https://github.com/icy/golo
  Date    : 2018 June 1st
  License : MIT
*/

/*
  Examples:

  $ go run golo.go -timeout 10  -port 4040 --no-bind -- /usr/bin/ssh MyServer -o "LocalForward localhost:4040 localhost:8888" -fN
  :: Port is available. App is not running
  :: Now staring application '/usr/bin/ssh' from .

  $ go run golo.go -timeout 10  -port 4040 --no-bind -- /usr/bin/ssh MyServer -o "LocalForward localhost:4040 localhost:8888" -fN
  :: Port is not available. App is running?
*/

package main

import (
  "fmt"
  "flag"
  "net"
  "os"
  "os/exec"
  "syscall"
  "log"
  "regexp"
  "strconv"
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

/*
  Set close-on-exec state for all fds >= 3
  The idea comes from
    https://github.com/golang/gofrontend/commit/651e71a729e5dcbd9dc14c1b59b6eff05bfe3d26
*/
func closeOnExec(state bool) {
  out, err := exec.Command("ls", fmt.Sprintf("/proc/%d/fd/", syscall.Getpid())).Output();
  if err != nil {
    log.Fatal(err);
  }
  pids := regexp.MustCompile("[ \t\n]").Split(fmt.Sprintf("%s", out), -1);
  i := 0;
  for i < len(pids) {
    if len(pids[i]) < 1 {
      i ++;
      continue;
    }
    pid, err := strconv.Atoi(pids[i]);
    if err != nil {
      log.Fatal(err);
    }
    if (pid > 2) {
      // FIXME: Check if fd is close
      if state {
        syscall.Syscall(syscall.SYS_FCNTL, uintptr(pid), syscall.FD_CLOEXEC, 0);
      } else {
        syscall.Syscall(syscall.SYS_FCNTL, uintptr(pid), 0, 0);
      }
    }
    i ++;
  }
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
  if *noBind == false {
    warn(fmt.Sprintf("Making sure all fd >= 3 is not close-on-exec\n"));
    closeOnExec(false);
  } else {
    if (*conn) != nil {
      (*conn).Close();
    }
    // https://golang.org/src/syscall/exec_unix.go?s=7214:7279#L244
    // Ruby > 1.8 has option to not close other fds before Exec
    // but Golang syscall.Exec() doesn't have that option
  }

  warn(fmt.Sprintf("Now staring application '%s' from %s\n", execPath, *workDir));
  err = syscall.Exec(execPath, cmdArgs, syscall.Environ());
  if err != nil {
    warn(fmt.Sprintf("Executing got error '%s'\n", err));
    os.Exit(1);
  }
}
