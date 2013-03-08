package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
// #include <git2/attr.h>
// extern int go_attr_callback(char *name, char *value, void *payload);
// extern int goAttrForEach(git_repository *repo, uint32_t flags, const char *path, void *payload);
import "C"
import (
	"reflect"
	"unsafe"
)

type AttrFlag uint32

const (
	ATTR_CHECK_FILE_THEN_INDEX AttrFlag = iota
	ATTR_CHECK_INDEX_THEN_FILE
	ATTR_CHECK_INDEX_ONLY
)

func (repo *Repository) AddAttrMacro(name, values string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalues := C.CString(values)
	defer C.free(unsafe.Pointer(cvalues))
	ecode := C.git_attr_add_macro(repo.git_repository, cname, cvalues)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (repo *Repository) FlushAttrCache() {
	C.git_attr_cache_flush(repo.git_repository)
}

func (repo *Repository) ForEachAttr(flags AttrFlag, path string, callback AttrCallback, payload interface{}) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	data := unsafe.Pointer(&attrCallbackWrapper{callback, payload})
	ecode := C.goAttrForEach(repo.git_repository, C.uint32_t(flags), cpath, data)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

type AttrCallback func(name, value string, payload interface{}) error

type attrCallbackWrapper struct {
	f AttrCallback
	d interface{}
}

//export go_attr_callback
func go_attr_callback(name, value *C.char, payload unsafe.Pointer) C.int {
	wrap := (*attrCallbackWrapper)(payload)
	err := wrap.f(C.GoString(name), C.GoString(value), wrap.d)
	if err != nil {
		return C.int(git_SUCCESS - 1)
	}
	return C.int(git_SUCCESS)
}

// TODO: Use varargs, find other places to use it as well.
func (repo *Repository) GetAttr(path, name string, flags ...AttrFlag) (string, error) {
	var value [1]int8
	cvalue := (*C.char)(&value[0])
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var cflags C.uint32_t
	for _, flag := range flags {
		cflags |= C.uint32_t(flag)
	}

	ecode := C.git_attr_get(&cvalue, repo.git_repository, cflags, cpath, cname)
	if ecode != git_SUCCESS {
		return "", gitError()
	}
	return C.GoString(cvalue), nil
}

func (repo *Repository) GetManyAttrs(path string, names []string, flags ...AttrFlag) ([]string, error) {
	var values [1]int8
	cvalues := (*C.char)(&values[0])
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	length := len(names)
	clength := C.size_t(length)
	cnames := make([]*C.char, length)
	for i := 0; i < length; i++ {
		cnames[i] = C.CString(names[i])
		defer C.free(unsafe.Pointer(cnames[i]))
	}

	var cflags C.uint32_t
	for _, flag := range flags {
		cflags |= C.uint32_t(flag)
	}

	ecode := C.git_attr_get_many(&cvalues, repo.git_repository, cflags, cpath, clength, &cnames[0])
	if ecode != git_SUCCESS {
		return nil, gitError()
	}

	// TODO: Find a safer way if one exists.
	var valuesSlice reflect.SliceHeader
	valuesSlice.Data = uintptr(unsafe.Pointer(cvalues))
	valuesSlice.Len = length
	valuesSlice.Cap = length
	cvalueStrings := *(*[]*C.char)(unsafe.Pointer(&valuesSlice))

	valuesOut := make([]string, length)
	for i := 0; i < length; i++ {
		valuesOut[i] = C.GoString(cvalueStrings[i])
	}

	return valuesOut, nil
}
