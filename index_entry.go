package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
// #include <git2/index.h>
import "C"

type IndexEntry struct {
	git_index_entry *C.git_index_entry
}

func (entry *IndexEntry) Stage() int {
	return int(C.git_index_entry_stage(entry.git_index_entry))
}

type IndexEntryUnmerged struct {
	git_index_entry_unmerged *C.git_index_entry_unmerged
}
