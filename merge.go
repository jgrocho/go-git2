package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"

func (repo *Repository) MergeBase(one, two Oid) (*Oid, error) {
	oid := new(Oid)
	ecode := C.git_merge_base(oid.git_oid, repo.git_repository, one.git_oid, two.git_oid)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}
