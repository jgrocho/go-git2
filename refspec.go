package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"unsafe"
)

type Refspec struct {
	git_refspec *C.git_refspec
}

func (refspec *Refspec) Destination() string {
	return C.GoString(C.git_refspec_dst(refspec.git_refspec))
}

func (refspec *Refspec) Source() string {
	return C.GoString(C.git_refspec_src(refspec.git_refspec))
}

func (refspec *Refspec) SourceMatches(refname string) bool {
	crefname := C.CString(refname)
	defer C.free(unsafe.Pointer(crefname))
	return C.git_refspec_src_matches(refspec.git_refspec, crefname) != c_FALSE
}

func (refspec *Refspec) Transform(name string) (string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var path [git_PATH_MAX]int8
	cpath := (*C.char)(&path[0])
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_refspec_transform(cpath, C.size_t(git_PATH_MAX), refspec.git_refspec, cname)
	if ecode != git_SUCCESS {
		return "", gitError()
	}
	return C.GoString(cpath), nil
}
