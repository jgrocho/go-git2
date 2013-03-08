package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"unsafe"
)

func OpenIndex(path string) (*Index, error) {
	idx := new(Index)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_index_open(&idx.git_index, cpath)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return idx, nil
}

type Index struct {
	git_index *C.git_index
}

func (idx *Index) Add(name string, stage int) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ecode := C.git_index_add(idx.git_index, cname, C.int(stage))
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (idx *Index) AddEntry(entry *IndexEntry) error {
	ecode := C.git_index_add2(idx.git_index, entry.git_index_entry)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (idx *Index) Append(name string, stage int) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ecode := C.git_index_append(idx.git_index, cname, C.int(stage))
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (idx *Index) AppendEntry(entry *IndexEntry) error {
	ecode := C.git_index_append2(idx.git_index, entry.git_index_entry)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (idx *Index) Clear() {
	C.git_index_clear(idx.git_index)
}

func (idx *Index) EntryCount() uint {
	return uint(C.git_index_entrycount(idx.git_index))
}

func (idx *Index) UnmergedEntryCount() uint {
	return uint(C.git_index_entrycount_unmerged(idx.git_index))
}

func (idx *Index) Find(path string) int {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return int(C.git_index_find(idx.git_index, cpath))
}

func (idx *Index) Free() {
	C.git_index_free(idx.git_index)
}

func (idx *Index) Get(n uint) *IndexEntry {
	entry := new(IndexEntry)
	entry.git_index_entry = C.git_index_get(idx.git_index, C.uint(n))
	return entry
}

func (idx *Index) GetUnmergedByIndex(n uint) *IndexEntryUnmerged {
	entry := new(IndexEntryUnmerged)
	entry.git_index_entry_unmerged = C.git_index_get_unmerged_byindex(idx.git_index, C.uint(n))
	if entry.git_index_entry_unmerged == nil {
		return nil
	}
	return entry
}

func (idx *Index) GetUnmergedByPath(path string) *IndexEntryUnmerged {
	entry := new(IndexEntryUnmerged)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	entry.git_index_entry_unmerged = C.git_index_get_unmerged_bypath(idx.git_index, cpath)
	if entry.git_index_entry_unmerged == nil {
		return nil
	}
	return entry
}

func (idx *Index) Read() error {
	ecode := C.git_index_read(idx.git_index)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (idx *Index) ReadTree(tree *Tree) error {
	ecode := C.git_index_read_tree(idx.git_index, tree.git_tree)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (idx *Index) Remove(n int) error {
	ecode := C.git_index_remove(idx.git_index, C.int(n))
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (idx *Index) Unique() {
	C.git_index_uniq(idx.git_index)
}

func (idx *Index) Write() error {
	ecode := C.git_index_write(idx.git_index)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}
