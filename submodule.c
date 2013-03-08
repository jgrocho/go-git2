#include <git2.h>

int go_submodule_callback2(const char *path, void *payload) {
	char *vpath = (char *)path;
	return go_submodule_callback(vpath, payload);
}

int goSubmoduleForEach(git_repository *repo, void *payload) {
	return git_submodule_foreach(repo, go_submodule_callback2, payload);
}
