#include <git2.h>

int go_tree_walk_callback(char *root, git_tree_entry *entry, void *payload) {
	return 0;
}

int go_tree_walk_callback2(const char *root, git_tree_entry *entry, void *payload) {
	char *vroot = (char *)root;
	return go_tree_walk_callback(vroot, entry, payload);
}

int goTreeWalk(git_tree *tree, int mode, void *payload) {
	return git_tree_walk(tree, go_tree_walk_callback2, mode, payload);
}

int go_treebuilder_filter(git_tree_entry *entry, void *payload) {
	return 0;
}

int go_treebuilder_filter2(const git_tree_entry *entry, void *payload) {
	git_tree_entry *ventry = (git_tree_entry *)entry;
	return go_treebuilder_filter(ventry, payload);
}

void goTreeBuilderFilter(git_treebuilder *builder, void *payload) {
	git_treebuilder_filter(builder, go_treebuilder_filter2, payload);
}
