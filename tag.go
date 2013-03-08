package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type Tag struct {
	git_tag *C.git_tag
}

func (tag *Tag) Free() {
	C.git_tag_free(tag.git_tag)
}

func (tag *Tag) Id() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_tag_id(tag.git_tag)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (tag *Tag) Message() string {
	return C.GoString(C.git_tag_message(tag.git_tag))
}

func (tag *Tag) Name() string {
	return C.GoString(C.git_tag_name(tag.git_tag))
}

func (tag *Tag) Peel() (*Object, error) {
	obj := new(Object)
	ecode := C.git_tag_peel(&obj.git_object, tag.git_tag)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return obj, nil
}

func (tag *Tag) Tagger() *Signature {
	sig := new(Signature)
	sig.git_signature = C.git_tag_tagger(tag.git_tag)
	if sig.git_signature == nil {
		return nil
	}
	return sig
}

func (tag *Tag) Target() (*Object, error) {
	obj := new(Object)
	ecode := C.git_tag_target(&obj.git_object, tag.git_tag)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return obj, nil
}

func (tag *Tag) TargetOid() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_tag_target_oid(tag.git_tag)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (tag *Tag) Type() ObjectType {
	return ObjectType(C.git_tag_type(tag.git_tag))
}

func (repo *Repository) CreateTag(name string, target *Object, tagger *Signature, message string, force bool) (*Oid, error) {
	oid := new(Oid)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cmessage := C.CString(message)
	defer C.free(unsafe.Pointer(cmessage))
	cforce := C.int(c_FALSE)
	if force {
		cforce = C.int(c_TRUE)
	}
	ecode := C.git_tag_create(oid.git_oid, repo.git_repository, cname, target.git_object, tagger.git_signature, cmessage, cforce)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

func (repo *Repository) CreateTagFromBuffer(buffer string, force bool) (*Oid, error) {
	oid := new(Oid)
	cbuffer := C.CString(buffer)
	defer C.free(unsafe.Pointer(cbuffer))
	cforce := C.int(c_FALSE)
	if force {
		cforce = C.int(c_TRUE)
	}
	ecode := C.git_tag_create_frombuffer(oid.git_oid, repo.git_repository, cbuffer, cforce)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

func (repo *Repository) CreateLightweightTag(name string, target Object, force bool) (*Oid, error) {
	oid := new(Oid)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cforce := C.int(c_FALSE)
	if force {
		cforce = C.int(c_TRUE)
	}
	ecode := C.git_tag_create_lightweight(oid.git_oid, repo.git_repository, cname, target.git_object, cforce)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

func (repo *Repository) DeleteTag(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ecode := C.git_tag_delete(repo.git_repository, cname)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (repo *Repository) TagList() ([]string, error) {
	var ctags C.git_strarray
	defer C.git_strarray_free(&ctags)
	ecode := C.git_tag_list(&ctags, repo.git_repository)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}

	// TODO: Find a safer way if one exists.
	var tagsSlice reflect.SliceHeader
	length := int(ctags.count)
	tagsSlice.Data = uintptr(unsafe.Pointer(ctags.strings))
	tagsSlice.Len = length
	tagsSlice.Cap = length
	ctagStrings := *(*[]*C.char)(unsafe.Pointer(&tagsSlice))

	tags := make([]string, length)
	for i := 0; i < len(ctagStrings); i++ {
		tags[i] = C.GoString(ctagStrings[i])
	}

	return tags, nil
}

func (repo *Repository) TagListMatch(pattern string) ([]string, error) {
	cpattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(cpattern))
	var ctags C.git_strarray
	defer C.git_strarray_free(&ctags)
	ecode := C.git_tag_list_match(&ctags, cpattern, repo.git_repository)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}

	// TODO: Find a safer way if one exists.
	var tagsSlice reflect.SliceHeader
	length := int(ctags.count)
	tagsSlice.Data = uintptr(unsafe.Pointer(ctags.strings))
	tagsSlice.Len = length
	tagsSlice.Cap = length
	ctagStrings := *(*[]*C.char)(unsafe.Pointer(&tagsSlice))

	tags := make([]string, length)
	for i := 0; i < len(ctagStrings); i++ {
		tags[i] = C.GoString(ctagStrings[i])
	}

	return tags, nil
}

func (repo *Repository) LookupTag(id *Oid) (*Tag, error) {
	tag := new(Tag)
	ecode := C.git_tag_lookup(&tag.git_tag, repo.git_repository, id.git_oid)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return tag, nil
}

func (repo *Repository) LookupTagPrefix(id *Oid, n uint) (*Tag, error) {
	tag := new(Tag)
	ecode := C.git_tag_lookup_prefix(&tag.git_tag, repo.git_repository, id.git_oid, C.uint(n))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return tag, nil
}
