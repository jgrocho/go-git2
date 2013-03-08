#include <git2.h>

int go_status_callback2(const char *path, unsigned int flags, void *payload) {
	char *vpath = (char *)path;
	return go_status_callback(vpath, flags, payload);
}

int goStatusForEach(git_repository *repo, void *payload) {
	return git_status_foreach(repo, go_status_callback2, payload);
}

int goStatusForEachExt(git_repository *repo, git_status_options *opts, void *payload) {
	return git_status_foreach_ext(repo, opts, go_status_callback2, payload);
}
