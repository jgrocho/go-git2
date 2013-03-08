package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
// #include <git2/branch.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type BranchType int

const (
	BRANCH_LOCAL BranchType = 1 << iota
	BRANCH_REMOTE
	BRANCH_ALL = BRANCH_LOCAL | BRANCH_REMOTE
)

func (repo *Repository) CreateBranch(name string, target *Object, force bool) (*Oid, error) {
	oid := new(Oid)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cforce := C.int(c_FALSE)
	if force {
		cforce = C.int(c_TRUE)
	}
	ecode := C.git_branch_create(oid.git_oid, repo.git_repository, cname, target.git_object, cforce)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

func (repo *Repository) DeleteBranch(name string, flag BranchType) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cflag := C.git_branch_t(flag)
	ecode := C.git_branch_delete(repo.git_repository, cname, cflag)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (repo *Repository) ListBranches(flags ...BranchType) ([]string, error) {
	var cnames C.git_strarray
	defer C.git_strarray_free(&cnames)
	var cflags C.uint
	if len(flags) == 0 {
		cflags = C.uint(BRANCH_ALL)
	} else {
		for _, flag := range flags {
			cflags |= C.uint(flag)
		}
	}
	ecode := C.git_branch_list(&cnames, repo.git_repository, cflags)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}

	// TODO: Find a safer way if one exists.
	var namesSlice reflect.SliceHeader
	length := int(cnames.count)
	namesSlice.Data = uintptr(unsafe.Pointer(cnames.strings))
	namesSlice.Len = length
	namesSlice.Cap = length
	cnameStrings := *(*[]*C.char)(unsafe.Pointer(&namesSlice))

	names := make([]string, length)
	for i := 0; i < len(cnameStrings); i++ {
		names[i] = C.GoString(cnameStrings[i])
	}

	return names, nil
}

func (repo *Repository) MoveBranch(oldName, newName string, force bool) error {
	coldName := C.CString(oldName)
	defer C.free(unsafe.Pointer(coldName))
	cnewName := C.CString(newName)
	defer C.free(unsafe.Pointer(cnewName))
	cforce := C.int(c_FALSE)
	if force {
		cforce = C.int(c_TRUE)
	}
	ecode := C.git_branch_move(repo.git_repository, coldName, cnewName, cforce)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}
