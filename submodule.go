package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
// extern int go_submodule_callback(char *path, void *payload);
// extern int goSubmoduleForEach(git_repository *repo, void *payload);
import "C"
import (
	"unsafe"
)

type Submodule struct {
	git_submodule *C.git_submodule
}

func (repo *Repository) ForEachSubmodule(callback SubmoduleCallback, payload interface{}) error {
	data := unsafe.Pointer(&submoduleCallbackWrapper{callback, payload})
	ecode := C.goSubmoduleForEach(repo.git_repository, data)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

//export go_submodule_callback
func go_submodule_callback(submodule *C.char, payload unsafe.Pointer) C.int {
	wrap := (*submoduleCallbackWrapper)(payload)
	err := wrap.f(C.GoString(submodule), wrap.d)
	if err != nil {
		return C.int(git_SUCCESS - 1)
	}
	return C.int(git_SUCCESS)
}

type SubmoduleCallback func(submodule string, payload interface{}) error

type submoduleCallbackWrapper struct {
	f SubmoduleCallback
	d interface{}
}

func (repo *Repository) LookupSubmodule(name string) (*Submodule, error) {
	submodule := new(Submodule)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ecode := C.git_submodule_lookup(&submodule.git_submodule, repo.git_repository, cname)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return submodule, nil
}
