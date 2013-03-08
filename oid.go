package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"unsafe"
)

func OidFromRaw(raw string) *Oid {
	oid := new(Oid)
	craw := C.CString(raw)
	// Ugly hack to get around git_oid_fromraw using unsigned char*
	crawp := unsafe.Pointer(craw)
	defer C.free(crawp)
	C.git_oid_fromraw(oid.git_oid, (*C.uchar)(crawp))
	return oid
}

func OidFromString(str string) *Oid {
	oid := new(Oid)
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	length := C.size_t(len(str))
	C.git_oid_fromstrn(oid.git_oid, cstr, length)
	return oid
}

type Oid struct {
	git_oid *C.git_oid
}

func (oid *Oid) String() string {
	buffer := make([]C.char, git_OID_HEXSZ+1)
	length := C.size_t(cap(buffer))
	return C.GoString(C.git_oid_tostr(&buffer[0], length, oid.git_oid))
}

func (oid *Oid) Compare(other *Oid) int {
	return int(C.git_oid_cmp(oid.git_oid, other.git_oid))
}

func (oid *Oid) CompareN(other *Oid, n uint) int {
	return int(C.git_oid_ncmp(oid.git_oid, other.git_oid, C.uint(n)))
}

func (oid *Oid) Copy() *Oid {
	newOid := new(Oid)
	C.git_oid_cpy(newOid.git_oid, oid.git_oid)
	return newOid
}

func (oid *Oid) IsZero() bool {
	return C.git_oid_iszero(oid.git_oid) != c_FALSE
}

func (oid *Oid) Path() string {
	var path [git_OID_HEXSZ + 1]int8
	cpath := (*C.char)(&path[0])
	defer C.free(unsafe.Pointer(cpath))
	C.git_oid_pathfmt(cpath, oid.git_oid)
	return C.GoString(cpath)
}

func (oid *Oid) Equal(str string) (bool, error) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	equal := C.git_oid_streq(oid.git_oid, cstr)
	if equal == 0 {
		return true, nil
	}
	return false, gitError()
}

func OidShortenNew(minLength int) *OidShorten {
	os := new(OidShorten)
	os.git_oid_shorten = C.git_oid_shorten_new(C.size_t(minLength))
	return os
}

type OidShorten struct {
	git_oid_shorten *C.git_oid_shorten
}

func (os *OidShorten) Add(oid string) (int, error) {
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))
	num := C.git_oid_shorten_add(os.git_oid_shorten, coid)
	if num < 0 {
		return 0, gitError()
	}
	return int(num), nil
}

func (os *OidShorten) Free() {
	C.git_oid_shorten_free(os.git_oid_shorten)
}
