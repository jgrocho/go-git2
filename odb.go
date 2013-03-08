package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"unsafe"
)

type OdbBackend struct {
	git_odb_backend *C.git_odb_backend
}

type OdbStream struct {
	git_odb_stream *C.git_odb_stream
}

func NewOdb() (*Odb, error) {
	odb := new(Odb)
	ecode := C.git_odb_new(&odb.git_odb)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return odb, nil
}

func OpenOdb(dir string) (*Odb, error) {
	odb := new(Odb)
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	ecode := C.git_odb_open(&odb.git_odb, cdir)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return odb, nil
}

type Odb struct {
	git_odb *C.git_odb
}

func (odb *Odb) AddBackend(backend *OdbBackend, priority int) error {
	ecode := C.git_odb_add_backend(odb.git_odb, backend.git_odb_backend, C.int(priority))
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (odb *Odb) AddAlternate(backend *OdbBackend, priority int) error {
	ecode := C.git_odb_add_alternate(odb.git_odb, backend.git_odb_backend, C.int(priority))
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (odb *Odb) Exists(oid *Oid) bool {
	return C.git_odb_exists(odb.git_odb, oid.git_oid) != c_FALSE
}

func (odb *Odb) Free() {
	C.git_odb_free(odb.git_odb)
}

func (odb *Odb) Hash(data []byte, form ObjectType) (*Oid, error) {
	oid := new(Oid)
	cdata := unsafe.Pointer(&data[0])
	defer C.free(cdata)
	length := C.size_t(len(data))
	ecode := C.git_odb_hash(oid.git_oid, cdata, length, C.git_otype(form))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

func (odb *Odb) HashFile(path string, form ObjectType) (*Oid, error) {
	oid := new(Oid)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_odb_hashfile(oid.git_oid, cpath, C.git_otype(form))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

func (odb *Odb) OpenRStream(oid *Oid) (*OdbStream, error) {
	stream := new(OdbStream)
	ecode := C.git_odb_open_rstream(&stream.git_odb_stream, odb.git_odb, oid.git_oid)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return stream, nil
}

func (odb *Odb) OpenWStream(size int, form ObjectType) (*OdbStream, error) {
	stream := new(OdbStream)
	ecode := C.git_odb_open_wstream(&stream.git_odb_stream, odb.git_odb, C.size_t(size), C.git_otype(form))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return stream, nil
}

type OdbObject struct {
	git_odb_object *C.git_odb_object
}

func (obj *OdbObject) Data() []byte {
	data := C.git_odb_object_data(obj.git_odb_object)
	length := C.git_odb_object_size(obj.git_odb_object)
	return C.GoBytes(data, C.int(length))
}

func (obj *OdbObject) Free() {
	C.git_odb_object_free(obj.git_odb_object)
}

func (obj *OdbObject) Id() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_odb_object_id(obj.git_odb_object)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (obj *OdbObject) Size() int {
	return int(C.git_odb_object_size(obj.git_odb_object))
}

func (obj *OdbObject) Type() ObjectType {
	return ObjectType(C.git_odb_object_type(obj.git_odb_object))
}

func (odb *Odb) Read(oid *Oid) (*OdbObject, error) {
	obj := new(OdbObject)
	ecode := C.git_odb_read(&obj.git_odb_object, odb.git_odb, oid.git_oid)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return obj, nil
}

func (odb *Odb) ReadHeader(oid *Oid) (int, ObjectType, error) {
	var clen C.size_t
	var ctype C.git_otype
	ecode := C.git_odb_read_header(&clen, &ctype, odb.git_odb, oid.git_oid)
	if ecode != git_SUCCESS {
		return int(clen), ObjectType(ctype), gitError()
	}
	return int(clen), ObjectType(ctype), nil
}

func (odb *Odb) ReadPrefix(oid *Oid, n uint) (*OdbObject, error) {
	obj := new(OdbObject)
	ecode := C.git_odb_read_prefix(&obj.git_odb_object, odb.git_odb, oid.git_oid, C.uint(n))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return obj, nil
}

func (odb *Odb) Write(data []byte, form ObjectType) (*Oid, error) {
	oid := new(Oid)
	cdata := unsafe.Pointer(&data[0])
	defer C.free(cdata)
	length := C.size_t(len(data))
	ecode := C.git_odb_write(oid.git_oid, odb.git_odb, cdata, length, C.git_otype(form))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}
