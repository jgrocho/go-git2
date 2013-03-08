#include <git2.h>
#include <git2/attr.h>

int go_attr_callback2(const char *name, const char *value, void *payload) {
	char *vname = (char *)name;
	char *vvalue = (char *)value;
	return go_attr_callback(vname, vvalue, payload);
}

int goAttrForEach(git_repository *repo, uint32_t flags, const char *path, void *payload) {
	return git_attr_foreach(repo, flags, path, go_attr_callback2, payload);
}
