package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"unsafe"
)

type Note struct {
	git_note *C.git_note
}

func (note *Note) Free() {
	C.git_note_free(note.git_note)
}

func (note *Note) Message() string {
	return C.GoString(C.git_note_message(note.git_note))
}

func (note *Note) Oid() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_note_oid(note.git_note)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (repo *Repository) CreateNote(author, committer *Signature, ref string, oid *Oid, note string) (*Oid, error) {
	out := new(Oid)
	var cref *C.char
	if ref != "" {
		cref = C.CString(ref)
		defer C.free(unsafe.Pointer(cref))
	}
	var cnote *C.char
	if note != "" {
		cnote = C.CString(note)
		defer C.free(unsafe.Pointer(cnote))
	}
	ecode := C.git_note_create(out.git_oid, repo.git_repository, author.git_signature, committer.git_signature, cref, oid.git_oid, cnote)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return out, nil
}

func (repo *Repository) DefaultNoteRef() (string, error) {
	var ref [1]int8
	cref := (*C.char)(&ref[0])
	ecode := C.git_note_default_ref(&cref, repo.git_repository)
	if ecode != git_SUCCESS {
		return "", gitError()
	}
	return C.GoString(cref), nil
}

func (repo *Repository) ReadNote(ref string, oid *Oid) (*Note, error) {
	note := new(Note)
	cref := C.CString(ref)
	defer C.free(unsafe.Pointer(cref))
	ecode := C.git_note_read(&note.git_note, repo.git_repository, cref, oid.git_oid)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return note, nil
}

func (repo *Repository) RemoveNote(ref string, author, committer *Signature, oid *Oid) error {
	cref := C.CString(ref)
	defer C.free(unsafe.Pointer(cref))
	ecode := C.git_note_remove(repo.git_repository, cref, author.git_signature, committer.git_signature, oid.git_oid)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}
