package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type Direction int

const (
	DIR_FETCH Direction = iota
	DIR_PUSH
)

func SupportedRemoteUrl(url string) bool {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	return C.git_remote_supported_url(curl) != c_FALSE
}

func ValidRemoteUrl(url string) bool {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	return C.git_remote_valid_url(curl) != c_FALSE
}

type Remote struct {
	git_remote *C.git_remote
}

func (remote *Remote) Connect(direction Direction) error {
	ecode := C.git_remote_connect(remote.git_remote, C.int(direction))
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (remote *Remote) Connected() bool {
	return C.git_remote_connected(remote.git_remote) != c_FALSE
}

func (remote *Remote) Disconnect() {
	C.git_remote_disconnect(remote.git_remote)
}

func (remote *Remote) Download(bytes *int64, stats *IndexerStats) error {
	ecode := C.git_remote_download(remote.git_remote, (*C.git_off_t)(bytes), stats.git_indexer_stats)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (remote *Remote) Fetchspec() *Refspec {
	refspec := new(Refspec)
	refspec.git_refspec = C.git_remote_fetchspec(remote.git_remote)
	if refspec.git_refspec == nil {
		return nil
	}
	return refspec
}

func (remote *Remote) Free() {
	C.git_remote_free(remote.git_remote)
}

/* TODO: Implement
func (remote *Remote) Ls(callback HeadListCallback, payload interface{}) error {
}
*/

func (remote *Remote) Name() string {
	return C.GoString(C.git_remote_name(remote.git_remote))
}

func (remote *Remote) Pushspec() *Refspec {
	refspec := new(Refspec)
	refspec.git_refspec = C.git_remote_pushspec(remote.git_remote)
	if refspec.git_refspec == nil {
		return nil
	}
	return refspec
}

func (remote *Remote) SetFetchspec(spec string) error {
	cspec := C.CString(spec)
	defer C.free(unsafe.Pointer(cspec))
	ecode := C.git_remote_set_fetchspec(remote.git_remote, cspec)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (remote *Remote) SetPushspec(spec string) error {
	cspec := C.CString(spec)
	defer C.free(unsafe.Pointer(cspec))
	ecode := C.git_remote_set_pushspec(remote.git_remote, cspec)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (remote *Remote) Save() error {
	ecode := C.git_remote_save(remote.git_remote)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

/* TODO: Implement
func (remote *Remote) UpdateTips(callback UpdateTipCallback) error {
}
*/

func (remote *Remote) Url() string {
	return C.GoString(C.git_remote_url(remote.git_remote))
}

func (repo *Repository) AddRemote(name, url string) (*Remote, error) {
	remote := new(Remote)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	ecode := C.git_remote_add(&remote.git_remote, repo.git_repository, cname, curl)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return remote, nil
}

func (repo *Repository) ListRemotes() ([]string, error) {
	var cremotes C.git_strarray
	defer C.git_strarray_free(&cremotes)
	ecode := C.git_remote_list(&cremotes, repo.git_repository)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}

	// TODO: Find a safer way if one exists.
	var remotesSlice reflect.SliceHeader
	length := int(cremotes.count)
	remotesSlice.Data = uintptr(unsafe.Pointer(cremotes.strings))
	remotesSlice.Len = length
	remotesSlice.Cap = length
	cremoteStrings := *(*[]*C.char)(unsafe.Pointer(&remotesSlice))

	remotes := make([]string, length)
	for i := 0; i < len(cremoteStrings); i++ {
		remotes[i] = C.GoString(cremoteStrings[i])
	}

	return remotes, nil
}

func (repo *Repository) LoadRemote(name string) (*Remote, error) {
	remote := new(Remote)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ecode := C.git_remote_load(&remote.git_remote, repo.git_repository, cname)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return remote, nil
}

func (repo *Repository) NewRemote(name, url, fetch string) (*Remote, error) {
	remote := new(Remote)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	cfetch := C.CString(fetch)
	defer C.free(unsafe.Pointer(cfetch))
	ecode := C.git_remote_new(&remote.git_remote, repo.git_repository, cname, curl, cfetch)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return remote, nil
}
