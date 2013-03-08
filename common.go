package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"errors"
)

const (
	git_SUCCESS = iota
)

const (
	c_FALSE = iota
	c_TRUE
)

const (
	git_PATH_MAX  = 4096
	git_OID_RAWSZ = 20
	git_OID_HEXSZ = git_OID_RAWSZ * 2
)

func Version() (int, int, int) {
	var cmajor C.int
	var cminor C.int
	var crev C.int
	C.git_libgit2_version(&cmajor, &cminor, &crev)
	return int(cmajor), int(cminor), int(crev)
}

func Init() {
	C.git_threads_init()
}

func Shutdown() {
	C.git_threads_shutdown()
}

func gitError() error {
	ge := C.giterr_last()
	msg := C.GoString(ge.message)
	C.giterr_clear()
	return errors.New(msg)
}
