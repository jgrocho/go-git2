package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"unsafe"
)

func ObjectString2Type(str string) ObjectType {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	return ObjectType(C.git_object_string2type(cstr))
}

func ObjectType2String(form ObjectType) string {
	return C.GoString(C.git_object_type2string(C.git_otype(form)))
}

func ObjectTypeIsLoose(form ObjectType) bool {
	return C.git_object_typeisloose(C.git_otype(form)) != c_FALSE
}

type ObjectType int

const (
	OBJ_ANY ObjectType = iota - 2
	OBJ_BAD
	OBJ__EXT1
	OBJ_COMMIT
	OBJ_TREE
	OBJ_BLOB
	OBJ_TAG
	OBJ__EXT2
	OBJ_OFS_DELTA
	OBJ_REF_DELTA
)

func (objt *ObjectType) String() string {
	return ObjectType2String(*objt)
}

func (objt *ObjectType) IsLoose() bool {
	return ObjectTypeIsLoose(*objt)
}

type Object struct {
	git_object *C.git_object
}

func (obj *Object) Free() {
	C.git_object_free(obj.git_object)
	C.free(unsafe.Pointer(obj.git_object))
}

func (obj *Object) Id() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_object_id(obj.git_object)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (obj *Object) Owner() *Repository {
	repo := new(Repository)
	repo.git_repository = C.git_object_owner(obj.git_object)
	if repo.git_repository == nil {
		return nil
	}
	return repo
}

func (obj *Object) Type() ObjectType {
	return ObjectType(C.git_object_type(obj.git_object))
}

func (repo *Repository) LookupObject(id *Oid, form ObjectType) (*Object, error) {
	obj := new(Object)
	ecode := C.git_object_lookup(&obj.git_object, repo.git_repository, id.git_oid, C.git_otype(form))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return obj, nil
}

func (repo *Repository) LookupObjectPrefix(id *Oid, n uint, form ObjectType) (*Object, error) {
	obj := new(Object)
	ecode := C.git_object_lookup_prefix(&obj.git_object, repo.git_repository, id.git_oid, C.uint(n), C.git_otype(form))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return obj, nil
}
