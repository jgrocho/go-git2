package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"unsafe"
)

type SortMode int

const SORT_NONE SortMode = iota
const (
	SORT_TOPOLOGICAL SortMode = 1 << iota
	SORT_TIME
	SORT_REVERSE
)

const (
	git_REVWALKOVER = -31
)

type Revwalk struct {
	git_revwalk *C.git_revwalk
}

func (revwalk *Revwalk) Free() {
	C.git_revwalk_free(revwalk.git_revwalk)
}

func (revwalk *Revwalk) Hide(oid *Oid) error {
	ecode := C.git_revwalk_hide(revwalk.git_revwalk, oid.git_oid)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (revwalk *Revwalk) HideGlob(glob string) error {
	cglob := C.CString(glob)
	defer C.free(unsafe.Pointer(cglob))
	ecode := C.git_revwalk_hide_glob(revwalk.git_revwalk, cglob)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (revwalk *Revwalk) HideHead() error {
	ecode := C.git_revwalk_hide_head(revwalk.git_revwalk)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (revwalk *Revwalk) HideRef(ref string) error {
	cref := C.CString(ref)
	defer C.free(unsafe.Pointer(cref))
	ecode := C.git_revwalk_hide_ref(revwalk.git_revwalk, cref)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (revwalk *Revwalk) Next() (*Oid, error) {
	oid := new(Oid)
	ecode := C.git_revwalk_next(oid.git_oid, revwalk.git_revwalk)
	if ecode == 0 {
		return oid, nil
	} else if ecode == git_REVWALKOVER {
		return nil, nil
	}
	return nil, gitError()
}

func (revwalk *Revwalk) Push(oid *Oid) error {
	ecode := C.git_revwalk_push(revwalk.git_revwalk, oid.git_oid)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (revwalk *Revwalk) PushGlob(glob string) error {
	cglob := C.CString(glob)
	defer C.free(unsafe.Pointer(cglob))
	ecode := C.git_revwalk_push_glob(revwalk.git_revwalk, cglob)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (revwalk *Revwalk) PushHead() error {
	ecode := C.git_revwalk_push_head(revwalk.git_revwalk)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (revwalk *Revwalk) PushRef(refname string) error {
	crefname := C.CString(refname)
	defer C.free(unsafe.Pointer(crefname))
	ecode := C.git_revwalk_push_ref(revwalk.git_revwalk, crefname)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (revwalk *Revwalk) Repository() *Repository {
	repo := new(Repository)
	repo.git_repository = C.git_revwalk_repository(revwalk.git_revwalk)
	if repo.git_repository == nil {
		return nil
	}
	return repo
}

func (revwalk *Revwalk) Reset() {
	C.git_revwalk_reset(revwalk.git_revwalk)
}

func (revwalk *Revwalk) Sorting(sort SortMode) {
	C.git_revwalk_sorting(revwalk.git_revwalk, C.uint(sort))
}

func (repo *Repository) NewRevwalk() (*Revwalk, error) {
	revwalk := new(Revwalk)
	ecode := C.git_revwalk_new(&revwalk.git_revwalk, repo.git_repository)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return revwalk, nil
}
