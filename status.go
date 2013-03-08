package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
// extern int go_status_callback(char *path, unsigned int flags, void *payload);
// extern int goStatusForEach(git_repository *repo, void *payload);
// extern int goStatusForEachExt(git_repository *repo, git_status_options *opts, void *payload);
import "C"
import (
	"unsafe"
)

type StatusFlag uint

const STATUS_CURRENT StatusFlag = iota
const (
	STATUS_INDEX_NEW StatusFlag = 1 << iota
	STATUS_INDEX_MODIFIED
	STATUS_INDEX_DELETED
	STATUS_WT_NEW
	STATUS_WT_MODIFIED
	STATUS_WT_DELETED
	STATUS_IGNORED
)

type StatusOptions struct {
	git_status_options *C.git_status_options
}

func (repo *Repository) ForEachStatus(callback StatusCallback, payload interface{}) error {
	data := unsafe.Pointer(&statusCallbackWrapper{callback, payload})
	ecode := C.goStatusForEach(repo.git_repository, data)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (repo *Repository) ForEachExtStatus(opts StatusOptions, callback StatusCallback, payload interface{}) error {
	data := unsafe.Pointer(&statusCallbackWrapper{callback, payload})
	ecode := C.goStatusForEachExt(repo.git_repository, opts.git_status_options, data)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

//export go_status_callback
func go_status_callback(path *C.char, flags C.uint, payload unsafe.Pointer) C.int {
	wrap := (*statusCallbackWrapper)(payload)
	err := wrap.f(C.GoString(path), StatusFlag(flags), wrap.d)
	if err != nil {
		return C.int(git_SUCCESS - 1)
	}
	return C.int(git_SUCCESS)
}

type StatusCallback func(path string, flags StatusFlag, payload interface{}) error

type statusCallbackWrapper struct {
	f StatusCallback
	d interface{}
}

func (repo *Repository) StatusFile(path string) (StatusFlag, error) {
	var cflags C.uint
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_status_file(&cflags, repo.git_repository, cpath)
	if ecode != git_SUCCESS {
		return StatusFlag(cflags), gitError()
	}
	return StatusFlag(cflags), nil
}

func (repo *Repository) ShouldIgnore(path string) (bool, error) {
	var cignored C.int
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_status_should_ignore(&cignored, repo.git_repository, cpath)
	if ecode != git_SUCCESS {
		return false, gitError()
	}
	return (cignored != c_FALSE), nil
}
