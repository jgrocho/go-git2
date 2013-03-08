#include <git2.h>

int go_cfg_foreach_callback2(const char *name, const char *value, void *payload) {
	char *vname = (char *)name;
	char *vvalue = (char *)value;
	return go_cfg_foreach_callback(vname, vvalue, payload);
}

int goCfgForEach(git_config *cfg, void *payload) {
	return git_config_foreach(cfg, go_cfg_foreach_callback2, payload);
}

int go_cfg_multivar_callback2(const char *value, void *data) {
	char *vvalue = (char *)value;
	return go_cfg_multivar_callback(vvalue, data);
}

int goCfgGetMultivar(git_config *cfg, const char *name, const char *regexp, void *data) {
	return git_config_get_multivar(cfg, name, regexp, go_cfg_multivar_callback2, data);
}
