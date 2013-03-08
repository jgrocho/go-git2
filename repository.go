package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"strings"
	"unsafe"
)

func Discover(start string, across_fs bool, ceilings []string) (string, error) {
	cstart := C.CString(start)
	defer C.free(unsafe.Pointer(cstart))

	var path [git_PATH_MAX]int8
	cpath := (*C.char)(&path[0])

	cacross_fs := C.int(c_FALSE)
	if across_fs {
		cacross_fs = C.int(c_TRUE)
	}

	ceiling_dirs := strings.Join(ceilings, git_PATH_LIST_SEPARATOR)
	cceiling_dirs := C.CString(ceiling_dirs)
	defer C.free(unsafe.Pointer(cceiling_dirs))

	ecode := C.git_repository_discover(cpath, git_PATH_MAX, cstart, cacross_fs, nil)
	if ecode != git_SUCCESS {
		return "", gitError()
	}
	return C.GoString(cpath), nil
}

func InitRepository(path string, bare bool) (*Repository, error) {
	repo := new(Repository)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	cbare := C.unsigned(c_FALSE)
	if bare {
		cbare = C.unsigned(c_TRUE)
	}
	ecode := C.git_repository_init(&repo.git_repository, cpath, cbare)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return repo, nil
}

func Open(path string) (*Repository, error) {
	repo := new(Repository)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_repository_open(&repo.git_repository, cpath)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return repo, nil
}

type Repository struct {
	git_repository *C.git_repository
}

func (repo *Repository) Config() (*Config, error) {
	config := new(Config)
	ecode := C.git_repository_config(&config.git_config, repo.git_repository)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return config, nil
}

func (repo *Repository) Free() {
	C.git_repository_free(repo.git_repository)
}

func (repo *Repository) SetConfig(config *Config) {
	C.git_repository_set_config(repo.git_repository, config.git_config)
}

func (repo *Repository) Path() string {
	return C.GoString(C.git_repository_path(repo.git_repository))
}

func (repo *Repository) Head() (*Reference, error) {
	ref := new(Reference)
	ecode := C.git_repository_head(&ref.git_reference, repo.git_repository)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return ref, nil
}

func (repo *Repository) Detached() (bool, error) {
	detached := C.git_repository_head_detached(repo.git_repository)
	if detached == c_TRUE {
		return true, nil
	} else if detached == c_FALSE {
		return false, nil
	}
	return false, gitError()
}

func (repo *Repository) Orphan() (bool, error) {
	orphan := C.git_repository_head_orphan(repo.git_repository)
	if orphan == c_TRUE {
		return true, nil
	} else if orphan == c_FALSE {
		return false, nil
	}
	return false, gitError()
}

func (repo *Repository) Index() (*Index, error) {
	index := new(Index)
	ecode := C.git_repository_index(&index.git_index, repo.git_repository)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return index, nil
}

func (repo *Repository) SetIndex(index *Index) {
	C.git_repository_set_index(repo.git_repository, index.git_index)
}

func (repo *Repository) Bare() bool {
	return bool(C.git_repository_is_bare(repo.git_repository) == 1)
}

func (repo *Repository) Empty() bool {
	return bool(C.git_repository_is_empty(repo.git_repository) == 1)
}

func (repo *Repository) Odb() (*Odb, error) {
	odb := new(Odb)
	ecode := C.git_repository_odb(&odb.git_odb, repo.git_repository)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return odb, nil
}

func (repo *Repository) SetOdb(odb *Odb) {
	C.git_repository_set_odb(repo.git_repository, odb.git_odb)
}

func (repo *Repository) Workdir() string {
	return C.GoString(C.git_repository_workdir(repo.git_repository))
}

func (repo *Repository) SetWorkdir(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_repository_set_workdir(repo.git_repository, cpath)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}
