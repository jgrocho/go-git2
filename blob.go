package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import "unsafe"

type Blob struct {
	git_blob *C.git_blob
}

func (blob *Blob) Content() []byte {
	size := C.git_blob_rawsize(blob.git_blob)
	content := C.git_blob_rawcontent(blob.git_blob)
	return C.GoBytes(unsafe.Pointer(content), C.int(size))
}

func (blob *Blob) Free() {
	C.git_blob_free(blob.git_blob)
}

func (repo *Repository) LookupBlob(oid *Oid) (*Blob, error) {
	blob := new(Blob)
	ecode := C.git_blob_lookup(&blob.git_blob, repo.git_repository, oid.git_oid)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return blob, nil
}

func (repo *Repository) LookupBlobPrefix(oid *Oid, n uint) (*Blob, error) {
	blob := new(Blob)
	ecode := C.git_blob_lookup_prefix(&blob.git_blob, repo.git_repository, oid.git_oid, C.uint(n))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return blob, nil
}

// Create a blob from byte slice.
func (repo *Repository) CreateBlob(buffer []byte) (*Oid, error) {
	oid := new(Oid)
	cbuffer := unsafe.Pointer(&buffer[0])
	defer C.free(cbuffer)
	length := C.size_t(len(buffer))
	ecode := C.git_blob_create_frombuffer(oid.git_oid, repo.git_repository, cbuffer, length)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

// Create a blob from a file.
func (repo *Repository) CreateBlobFromFile(path string) (*Oid, error) {
	oid := new(Oid)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_blob_create_fromdisk(oid.git_oid, repo.git_repository, cpath)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

// Create a blob from a file, relative to the repo's workdir.
func (repo *Repository) CreateBlobFromWorkdir(path string) (*Oid, error) {
	oid := new(Oid)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_blob_create_fromfile(oid.git_oid, repo.git_repository, cpath)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}
