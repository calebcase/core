// Package core provides facilities for creating core dumps.
package core

// #include <signal.h>
// #include <sys/types.h>
// #include <sys/wait.h>
// #include <unistd.h>
//
// pid_t
// dump()
// {
//   pid_t pid = fork();
//
//   if (pid == 0) {
//     pause();
//   } else {
//     int status = kill(pid, SIGABRT);
//
//     if (status != 0) {
//       return status;
//     }
//
//     return waitpid(pid, NULL, 0);
//   }
// }
import "C"
import (
	"os"
	"os/exec"
	"strconv"

	"github.com/inconshreveable/log15"
)

var Log = log15.New()

func init() {
	Log.SetHandler(log15.DiscardHandler())
}

// Dump the calling thread's state. This is done with a fork and SIGABRT. Only
// one thread will be present in the dump.
//
// The environment variable `GOTRACEBACK` needs to be set to `crash` otherwise
// the Go runtime will not produce a core dump.
//
// The ulimit for core dumps should also be set to allow them: `ulimit -u
// unlimited`.
func DumpSelf() (int, error) {
	pid, err := C.dump()
	Log.Debug("dump", log15.Ctx{
		"pid": pid,
		"err": err,
	})

	return int(pid), err
}

// Dump all threads. This is done using the external tool `gcore`. `gcore` is
// part of GDB's binutils and must be installed for this method to work. The
// core dump will contain all threads of the current process.
func DumpAll() (pid int, err error) {
	pid = os.Getpid()

	cmd := exec.Command("gcore", strconv.Itoa(pid))

	out, err := cmd.CombinedOutput()
	Log.Debug("gcore", log15.Ctx{
		"pid": pid,
		"out": string(out),
		"err": err,
	})
	if err != nil {
		return pid, err
	}

	return pid, nil
}

// Same as `DumpAll`, but use the provided prefix for the core file path.
func DumpAllTo(prefix string) (pid int, err error) {
	pid = os.Getpid()

	cmd := exec.Command("gcore", "-o", prefix, strconv.Itoa(pid))

	out, err := cmd.CombinedOutput()
	Log.Debug("gcore", log15.Ctx{
		"prefix": prefix,
		"pid":    pid,
		"out":    string(out),
		"err":    err,
	})
	if err != nil {
		return pid, err
	}

	return pid, nil
}
