package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"
import (
	"unsafe"
)

func NewIndexer(packname string) (*Indexer, error) {
	idxr := new(Indexer)
	cpackname := C.CString(packname)
	defer C.free(unsafe.Pointer(cpackname))
	ecode := C.git_indexer_new(&idxr.git_indexer, cpackname)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return idxr, nil
}

type Indexer struct {
	git_indexer *C.git_indexer
}

func (idxr *Indexer) Free() {
	C.git_indexer_free(idxr.git_indexer)
}

func (idxr *Indexer) Hash() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_indexer_hash(idxr.git_indexer)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}

func (idxr *Indexer) Run(stats *IndexerStats) error {
	ecode := C.git_indexer_run(idxr.git_indexer, stats.git_indexer_stats)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (idxr *Indexer) Write() error {
	ecode := C.git_indexer_write(idxr.git_indexer)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

type IndexerStats struct {
	git_indexer_stats *C.git_indexer_stats
}

type IndexerStream struct {
	git_indexer_stream *C.git_indexer_stream
}

func NewIndexerStream(dir string) (*IndexerStream, error) {
	stream := new(IndexerStream)
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	ecode := C.git_indexer_stream_new(&stream.git_indexer_stream, cdir)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return stream, nil
}

func (stream *IndexerStream) Add(data []byte, stats *IndexerStats) error {
	cdata := unsafe.Pointer(&data[0])
	defer C.free(cdata)
	length := C.size_t(len(data))
	ecode := C.git_indexer_stream_add(stream.git_indexer_stream, cdata, length, stats.git_indexer_stats)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (stream *IndexerStream) Finalize(stats IndexerStats) error {
	ecode := C.git_indexer_stream_finalize(stream.git_indexer_stream, stats.git_indexer_stats)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (stream *IndexerStream) Free() {
	C.git_indexer_stream_free(stream.git_indexer_stream)
}

func (stream *IndexerStream) Hash() *Oid {
	oid := new(Oid)
	oid.git_oid = C.git_indexer_stream_hash(stream.git_indexer_stream)
	if oid.git_oid == nil {
		return nil
	}
	return oid
}
