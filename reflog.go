package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"unsafe"
)

type Reflog struct {
	git_reflog *C.git_reflog
}

func (reflog *Reflog) Count() uint {
	return uint(C.git_reflog_entrycount(reflog.git_reflog))
}

func (reflog *Reflog) EntryByIndex(idx uint) *ReflogEntry {
	entry := new(ReflogEntry)
	entry.git_reflog_entry = C.git_reflog_entry_byindex(reflog.git_reflog, C.uint(idx))
	return entry
}

type ReflogEntry struct {
	git_reflog_entry *C.git_reflog_entry
}

func (entry *ReflogEntry) Committer() *Signature {
	sig := new(Signature)
	sig.git_signature = C.git_reflog_entry_committer(entry.git_reflog_entry)
	return sig
}

func (entry *ReflogEntry) Msg() string {
	return C.GoString(C.git_reflog_entry_msg(entry.git_reflog_entry))
}

func (entry *ReflogEntry) NewOid() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_reflog_entry_oidnew(entry.git_reflog_entry)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (entry *ReflogEntry) OldOid() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_reflog_entry_oidold(entry.git_reflog_entry)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (ref *Reference) DeleteReflog() error {
	ecode := C.git_reflog_delete(ref.git_reference)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (ref *Reference) ReadReflog() (*Reflog, error) {
	reflog := new(Reflog)
	ecode := C.git_reflog_read(&reflog.git_reflog, ref.git_reference)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return reflog, nil
}

func (ref *Reference) RenameReflog(newName string) error {
	cnewName := C.CString(newName)
	defer C.free(unsafe.Pointer(cnewName))
	ecode := C.git_reflog_rename(ref.git_reference, cnewName)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (ref *Reference) WriteReflog(oldOid *Oid, committer *Signature, msg string) error {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	ecode := C.git_reflog_write(ref.git_reference, oldOid.git_oid, committer.git_signature, cmsg)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}
