package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"time"
	"unsafe"
)

type Commit struct {
	git_commit *C.git_commit
}

func (commit *Commit) Author() *Signature {
	sig := new(Signature)
	sig.git_signature = C.git_commit_author(commit.git_commit)
	if sig.git_signature == nil {
		return nil
	}
	return sig
}

func (commit *Commit) Committer() *Signature {
	sig := new(Signature)
	sig.git_signature = C.git_commit_committer(commit.git_commit)
	if sig.git_signature == nil {
		return nil
	}
	return sig
}

func (commit *Commit) Free() {
	C.git_commit_free(commit.git_commit)
}

func (commit *Commit) Id() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_commit_id(commit.git_commit)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (commit *Commit) Message() string {
	return C.GoString(C.git_commit_message(commit.git_commit))
}

func (commit *Commit) MessageEncoding() string {
	return C.GoString(C.git_commit_message_encoding(commit.git_commit))
}

func (commit *Commit) Parent(n uint) (*Commit, error) {
	parent := new(Commit)
	ecode := C.git_commit_parent(&parent.git_commit, commit.git_commit, C.uint(n))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return parent, nil
}

func (commit *Commit) ParentOid(n uint) (*Oid, error) {
	poid := new(Oid)
	poid.git_oid = C.git_commit_parent_oid(commit.git_commit, C.uint(n))
	if poid.git_oid == nil {
		return nil, gitError()
	}
	return poid, nil
}

func (commit *Commit) ParentCount() uint {
	return uint(C.git_commit_parentcount(commit.git_commit))
}

func (commit *Commit) Time() time.Time {
	t := time.Unix(int64(C.git_commit_time(commit.git_commit)), 0)
	l := time.FixedZone("", int(C.git_commit_time_offset(commit.git_commit))*60)
	return t.In(l)
}

func (commit *Commit) Tree() (*Tree, error) {
	tree := new(Tree)
	ecode := C.git_commit_tree(&tree.git_tree, commit.git_commit)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return tree, nil
}

func (commit *Commit) TreeOid() (*Oid, error) {
	toid := new(Oid)
	toid.git_oid = C.git_commit_tree_oid(commit.git_commit)
	if toid.git_oid == nil {
		return nil, gitError()
	}
	return toid, nil
}

func (repo *Repository) CreateCommit(ref string, author, committer *Signature, encoding, message string, tree *Tree, parents... *Commit) (*Oid, error) {
	oid := new(Oid)
	oid.git_oid = new(C.git_oid)
	cref := C.CString(ref)
	defer C.free(unsafe.Pointer(cref))
	cencoding := C.CString(encoding)
	defer C.free(unsafe.Pointer(cencoding))
	cmessage := C.CString(message)
	defer C.free(unsafe.Pointer(cmessage))

	parentCount := len(parents)
	cparentCount := C.int(parentCount)
	var cparentsSlice []*C.git_commit
	var cparents **C.git_commit
	if parentCount > 0 {
		cparentsSlice = make([]*C.git_commit, parentCount)
		for i := 0; i < parentCount; i++ {
			cparentsSlice[i] = parents[i].git_commit
		}
		cparents = &cparentsSlice[0]
	}

	ecode := C.git_commit_create(oid.git_oid, repo.git_repository, cref, author.git_signature, committer.git_signature, cencoding, cmessage, tree.git_tree, cparentCount, cparents)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

func (repo *Repository) LookupCommit(oid *Oid) (*Commit, error) {
	commit := new(Commit)
	ecode := C.git_commit_lookup(&commit.git_commit, repo.git_repository, oid.git_oid)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return commit, nil
}

func (repo *Repository) LookupCommitPrefix(oid *Oid, n uint) (*Commit, error) {
	commit := new(Commit)
	ecode := C.git_commit_lookup_prefix(&commit.git_commit, repo.git_repository, oid.git_oid, C.unsigned(n))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return commit, nil
}
