package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
// extern int go_tree_walk_callback(char *root, git_tree_entry *entry, void *payload);
// extern int goTreeWalk(git_tree *tree, int mode, void *payload);
// extern int go_treebuilder_filter(git_tree_entry *entry, void *payload);
// extern int goTreeBuilderFilter(git_treebuilder *builder, void *payload);
import "C"
import (
	"unsafe"
)

type TreeWalkMode int

const (
	TREEWALK_PRE TreeWalkMode = iota
	TREEWALK_POST
)

type FileMode uint

const (
	FILEMODE_NEW             FileMode = 0000000
	FILEMODE_TREE                     = 0040000
	FILEMODE_BLOB                     = 0100644
	FILEMODE_BLOB_EXECUTABLE          = 0100755
	FILEMODE_LINK                     = 0120000
	FILEMODE_COMMIT                   = 0160000
)

type Tree struct {
	git_tree *C.git_tree
}

func (tree *Tree) CreateTreeBuilder() (*TreeBuilder, error) {
	builder := new(TreeBuilder)
	ecode := C.git_treebuilder_create(&builder.git_treebuilder, tree.git_tree)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return builder, nil
}

func (tree *Tree) EntryByIndex(idx uint) *TreeEntry {
	entry := new(TreeEntry)
	entry.git_tree_entry = C.git_tree_entry_byindex(tree.git_tree, C.uint(idx))
	return entry
}

func (tree *Tree) EntryByName(name string) *TreeEntry {
	entry := new(TreeEntry)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	entry.git_tree_entry = C.git_tree_entry_byname(tree.git_tree, cname)
	return entry
}

func (tree *Tree) EntryCount() uint {
	return uint(C.git_tree_entrycount(tree.git_tree))
}

func (tree *Tree) Free() {
	C.git_tree_free(tree.git_tree)
}

func (tree *Tree) Id() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_tree_id(tree.git_tree)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (tree *Tree) Subtree(path string) (*Tree, error) {
	subtree := new(Tree)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_tree_get_subtree(&subtree.git_tree, tree.git_tree, cpath)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return subtree, nil
}

/* TODO: Resolve error: Go type not supported in export: [0]byte
func (tree *Tree) Walk(callback TreeWalkCallback, mode TreeWalkMode, payload interface{}) error {
	data := unsafe.Pointer(&treeWalkCallbackWrapper{ callback, payload })
	ecode := C.goTreeWalk(tree.git_tree, C.int(mode), data)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

//export go_tree_walk_callback
func go_tree_walk_callback(root *C.char, centry *C.git_tree_entry, payload unsafe.Pointer) C.int {
	wrap := (*treeWalkCallbackWrapper)(payload)
	entry := &TreeEntry{centry}
	err := wrap.f(C.GoString(root), entry, wrap.d)
	if err != nil {
		return C.int(git_SUCCESS-1)
	}
	return C.int(git_SUCCESS)
}

type TreeWalkCallback func(root string, entry *TreeEntry, payload interface{})

type treeWalkCallbackWrapper struct {
	f TreeWalkCallback
	d interface{}
}
*/

type TreeBuilder struct {
	git_treebuilder *C.git_treebuilder
}

func (builder *TreeBuilder) Clear() {
	C.git_treebuilder_clear(builder.git_treebuilder)
}

/* TODO: Resolve error: Go type not supported in export: [0]byte
func (builder *TreeBuilder) Filter(filter TreeBuilderFilter, payload unsafe.Pointer) error {
	data := unsafe.Pointer(treeBuilderWrapper{ filter, payload })
	C.goTreeBuilderFilter(builder.git_treebuilder, data)
}

//export go_treebuilder_filter
func go_treebuilder_filter(centry *C.git_tree_entry, payload unsafe.Pointer) C.int {
	wrap := (*treeBuilderWrapper)(payload)
	entry := &TreeEntry{centry}
	err := wrap.f(entry, wrap.d)
	if err != nil {
		return C.int(git_SUCCESS-1)
	}
	return C.int(git_SUCCESS)
}

type TreeBuilderFilter func(entry *TreeEntry, payload interface{})

type treeBuilderWrapper struct {
	f TreeBuilderFilter
	d interface{}
}
*/

func (builder *TreeBuilder) Free() {
	C.git_treebuilder_free(builder.git_treebuilder)
}

func (builder *TreeBuilder) Get(filename string) *TreeEntry {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	entry := new(TreeEntry)
	entry.git_tree_entry = C.git_treebuilder_get(builder.git_treebuilder, cfilename)
	if entry.git_tree_entry == nil {
		return nil
	}
	return entry
}

func (builder *TreeBuilder) Insert(filename string, id *Oid, attributes FileMode) (*TreeEntry, error) {
	entry := new(TreeEntry)
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	ecode := C.git_treebuilder_insert(&entry.git_tree_entry, builder.git_treebuilder, cfilename, id.git_oid, C.uint(attributes))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return entry, nil
}

func (builder *TreeBuilder) Remove(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	ecode := C.git_treebuilder_remove(builder.git_treebuilder, cfilename)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (builder *TreeBuilder) Write(repo *Repository) (*Oid, error) {
	oid := new(Oid)
	ecode := C.git_treebuilder_write(oid.git_oid, repo.git_repository, builder.git_treebuilder)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

type TreeEntry struct {
	git_tree_entry *C.git_tree_entry
}

func (entry *TreeEntry) Attributes() uint {
	return uint(C.git_tree_entry_attributes(entry.git_tree_entry))
}

func (entry *TreeEntry) Id() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_tree_entry_id(entry.git_tree_entry)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (entry *TreeEntry) Name() string {
	return C.GoString(C.git_tree_entry_name(entry.git_tree_entry))
}

func (entry *TreeEntry) Type() ObjectType {
	return ObjectType(C.git_tree_entry_type(entry.git_tree_entry))
}

func (idx *Index) CreateTree() (*Oid, error) {
	oid := new(Oid)
	oid.git_oid = new(C.git_oid)
	ecode := C.git_tree_create_fromindex(oid.git_oid, idx.git_index)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return oid, nil
}

func (repo *Repository) LookupTree(id *Oid) (*Tree, error) {
	tree := new(Tree)
	ecode := C.git_tree_lookup(&tree.git_tree, repo.git_repository, id.git_oid)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return tree, nil
}

func (repo *Repository) LookupTreePrefix(id *Oid, n uint) (*Tree, error) {
	tree := new(Tree)
	ecode := C.git_tree_lookup_prefix(&tree.git_tree, repo.git_repository, id.git_oid, C.uint(n))
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return tree, nil
}
