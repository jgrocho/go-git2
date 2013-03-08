package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import ()

type DiffOptions struct {
	git_diff_options *C.git_diff_options
}
