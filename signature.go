package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"time"
	"unsafe"
)

func NewSignature(name, email string, when time.Time) (*Signature, error) {
	sig := new(Signature)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cemail := C.CString(email)
	defer C.free(unsafe.Pointer(cemail))
	ctime := C.git_time_t(when.Unix())
	_, offset := when.Zone()
	coffset := C.int(offset / 60)
	ecode := C.git_signature_new(&sig.git_signature, cname, cemail, ctime, coffset)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return sig, nil
}

func SignatureNow(name, email string) (*Signature, error) {
	sig := new(Signature)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cemail := C.CString(email)
	defer C.free(unsafe.Pointer(cemail))
	ecode := C.git_signature_now(&sig.git_signature, cname, cemail)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return sig, nil
}

type Signature struct {
	git_signature *C.git_signature
}

func (sig *Signature) Duplicate() *Signature {
	newSig := new(Signature)
	newSig.git_signature = C.git_signature_dup(sig.git_signature)
	return newSig
}

func (sig *Signature) Free() {
	C.git_signature_free(sig.git_signature)
}
