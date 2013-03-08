package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type RefType int

const REF_INVALID RefType = iota
const (
	REF_OID RefType = 1 << iota
	REF_SYMBOLIC
	REF_PACKED
	REF_HAS_PEEL
	REF_LISTALL = REF_OID | REF_SYMBOLIC | REF_PACKED
)

type Reference struct {
	git_reference *C.git_reference
}

func (ref *Reference) Compare(other *Reference) int {
	return int(C.git_reference_cmp(ref.git_reference, other.git_reference))
}

func (ref *Reference) Delete() error {
	ecode := C.git_reference_delete(ref.git_reference)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (ref *Reference) Free() {
	C.git_reference_free(ref.git_reference)
}

func (ref *Reference) IsPacked() bool {
	return C.git_reference_is_packed(ref.git_reference) != c_FALSE
}

func (ref *Reference) Name() string {
	return C.GoString(C.git_reference_name(ref.git_reference))
}

func (ref *Reference) Oid() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_reference_oid(ref.git_reference)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (ref *Reference) SetOid(oid *Oid) error {
	ecode := C.git_reference_set_oid(ref.git_reference, oid.git_oid)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (ref *Reference) Owner() *Repository {
	repo := new(Repository)
	repo.git_repository = C.git_reference_owner(ref.git_reference)
	if repo.git_repository == nil {
		return nil
	}
	return repo
}

func (ref *Reference) Reload() error {
	ecode := C.git_reference_reload(ref.git_reference)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (ref *Reference) Rename(newName string, force bool) error {
	cnewName := C.CString(newName)
	defer C.free(unsafe.Pointer(cnewName))
	cforce := C.int(c_FALSE)
	if force {
		cforce = C.int(c_TRUE)
	}
	ecode := C.git_reference_rename(ref.git_reference, cnewName, cforce)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (ref *Reference) Resolve() (*Reference, error) {
	resolved := new(Reference)
	ecode := C.git_reference_resolve(&resolved.git_reference, ref.git_reference)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return ref, nil
}

func (ref *Reference) Target() string {
	return C.GoString(C.git_reference_target(ref.git_reference))
}

func (ref *Reference) SetTarget(target string) error {
	ctarget := C.CString(target)
	defer C.free(unsafe.Pointer(ctarget))
	ecode := C.git_reference_set_target(ref.git_reference, ctarget)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (ref *Reference) Type() RefType {
	return RefType(C.git_reference_type(ref.git_reference))
}

func (repo *Repository) CreateOidRef(name string, oid *Oid, force bool) (*Reference, error) {
	ref := new(Reference)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cforce := C.int(c_FALSE)
	if force {
		cforce = C.int(c_TRUE)
	}
	ecode := C.git_reference_create_oid(&ref.git_reference, repo.git_repository, cname, oid.git_oid, cforce)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return ref, nil
}

func (repo *Repository) CreateSymbolicRef(name, target string, force bool) (*Reference, error) {
	ref := new(Reference)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ctarget := C.CString(target)
	defer C.free(unsafe.Pointer(ctarget))
	cforce := C.int(c_FALSE)
	if force {
		cforce = C.int(c_TRUE)
	}
	ecode := C.git_reference_create_symbolic(&ref.git_reference, repo.git_repository, cname, ctarget, cforce)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return ref, nil
}

func (repo *Repository) ListReferences(flags RefType) ([]string, error) {
	var crefs C.git_strarray
	defer C.git_strarray_free(&crefs)
	ecode := C.git_reference_list(&crefs, repo.git_repository, C.uint(flags))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}

	// TODO: Find a safer way if one exists.
	var refsSlice reflect.SliceHeader
	length := int(crefs.count)
	refsSlice.Data = uintptr(unsafe.Pointer(crefs.strings))
	refsSlice.Len = length
	refsSlice.Cap = length
	crefStrings := *(*[]*C.char)(unsafe.Pointer(&refsSlice))

	refs := make([]string, length)
	for i := 0; i < len(crefStrings); i++ {
		refs[i] = C.GoString(crefStrings[i])
	}
	return refs, nil
}

func (repo *Repository) LookupReference(name string) (*Reference, error) {
	ref := new(Reference)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ecode := C.git_reference_lookup(&ref.git_reference, repo.git_repository, cname)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return ref, nil
}

func (repo *Repository) ReferenceNameToOid(name string) (*Oid, error) {
	oid := new(Oid)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ecode := C.git_reference_name_to_oid(oid.git_oid, repo.git_repository, cname)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

func (repo *Repository) PackAllRefs() error {
	ecode := C.git_reference_packall(repo.git_repository)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}
