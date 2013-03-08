package git2

// #cgo pkg-config: libgit2
// #include <git2.h>
// extern int go_cfg_foreach_callback(char *name, char *value, void *payload);
// extern int goCfgForEach(git_config *cfg, void *payload);
// extern int go_cfg_multivar_callback(char *value, void *data);
// extern int goCfgGetMultivar(git_config *cfg, const char *name, const char *regexp, void *data);
import "C"
import (
	"unsafe"
)

func FindGlobalConfig() string {
	var path [git_PATH_MAX]int8
	cpath := (*C.char)(&path[0])
	ecode := C.git_config_find_global(cpath, C.size_t(git_PATH_MAX))
	if ecode != git_SUCCESS {
		return ""
	}
	return C.GoString(cpath)
}

func FindSystemConfig() string {
	var path [git_PATH_MAX]int8
	cpath := (*C.char)(&path[0])
	ecode := C.git_config_find_system(cpath, C.size_t(git_PATH_MAX))
	if ecode != git_SUCCESS {
		return ""
	}
	return C.GoString(cpath)
}

func OpenConfigOnDisk(path string) (*Config, error) {
	cfg := new(Config)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_config_open_ondisk(&cfg.git_config, cpath)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return cfg, nil
}

func OpenGlobalConfig() (*Config, error) {
	cfg := new(Config)
	ecode := C.git_config_open_global(&cfg.git_config)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return cfg, nil
}

type Config struct {
	git_config *C.git_config
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	ecode := C.git_config_new(&cfg.git_config)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	return cfg, nil
}

func (cfg *Config) AddFile(file *ConfigFile, priority int) error {
	ecode := C.git_config_add_file(cfg.git_config, file.git_config_file, C.int(priority))
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (cfg *Config) AddFileOnDisk(path string, priority int) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ecode := C.git_config_add_file_ondisk(cfg.git_config, cpath, C.int(priority))
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (cfg *Config) Delete(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ecode := C.git_config_delete(cfg.git_config, cname)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (cfg *Config) ForEach(callback ConfigForEachCallback, payload interface{}) error {
	data := unsafe.Pointer(&cfgForEachCallbackWrapper{callback, payload})
	ecode := C.goCfgForEach(cfg.git_config, data)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

//export go_cfg_foreach_callback
func go_cfg_foreach_callback(name, value *C.char, payload unsafe.Pointer) C.int {
	wrap := (*cfgForEachCallbackWrapper)(payload)
	err := wrap.f(C.GoString(name), C.GoString(value), wrap.d)
	if err != nil {
		// In v0.17.0 this does nothing, I believe it is fixed in HEAD (as of 2013-03-05).
		return C.int(git_SUCCESS - 1)
	}
	return C.int(git_SUCCESS)
}

type ConfigForEachCallback func(name, value string, payload interface{}) error

type cfgForEachCallbackWrapper struct {
	f ConfigForEachCallback
	d interface{}
}

func (cfg *Config) Free() {
	C.git_config_free(cfg.git_config)
}

func (cfg *Config) GetBool(name string) (bool, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var cval C.int
	ecode := C.git_config_get_bool(&cval, cfg.git_config, cname)
	if ecode != git_SUCCESS {
		return false, gitError()
	}
	return (cval != c_FALSE), nil
}

func (cfg *Config) SetBool(name string, value bool) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.int(c_FALSE)
	if value {
		cvalue = C.int(c_TRUE)
	}
	ecode := C.git_config_set_bool(cfg.git_config, cname, cvalue)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (cfg *Config) GetInt32(name string) (int32, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var cval C.int32_t
	ecode := C.git_config_get_int32(&cval, cfg.git_config, cname)
	if ecode != git_SUCCESS {
		return 0, gitError()
	}
	return int32(cval), nil
}

func (cfg *Config) SetInt32(name string, value int32) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.int32_t(value)
	ecode := C.git_config_set_int32(cfg.git_config, cname, cvalue)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (cfg *Config) GetInt64(name string) (int64, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var cval C.int64_t
	ecode := C.git_config_get_int64(&cval, cfg.git_config, cname)
	if ecode != git_SUCCESS {
		return 0, gitError()
	}
	return int64(cval), nil
}

func (cfg *Config) SetInt64(name string, value int64) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.int64_t(value)
	ecode := C.git_config_set_int64(cfg.git_config, cname, cvalue)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (cfg *Config) GetString(name string) (string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	// TODO: Is there a better way to pass a string pointer to C?
	var val [1]int8
	cval := (*C.char)(&val[0])
	ecode := C.git_config_get_string(&cval, cfg.git_config, cname)
	if ecode != git_SUCCESS {
		return "", gitError()
	}
	return C.GoString(cval), nil
}

func (cfg *Config) SetString(name, value string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	ecode := C.git_config_set_string(cfg.git_config, cname, cvalue)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

func (cfg *Config) GetMultivar(name, regexp string, callback ConfigMultivarCallback, data interface{}) error {
	payload := unsafe.Pointer(&cfgMultivarCallbackWrapper{callback, data})
	var cname *C.char
	// Treat an empty string as nil(NULL); avoids using regex in libgit2
	if name != "" {
		cname = C.CString(name)
		defer C.free(unsafe.Pointer(cname))
	}
	cregexp := C.CString(regexp)
	defer C.free(unsafe.Pointer(cregexp))
	ecode := C.goCfgGetMultivar(cfg.git_config, cname, cregexp, payload)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

//export go_cfg_multivar_callback
func go_cfg_multivar_callback(value *C.char, data unsafe.Pointer) C.int {
	wrap := (*cfgMultivarCallbackWrapper)(data)
	err := wrap.f(C.GoString(value), wrap.d)
	if err != nil {
		return C.int(git_SUCCESS - 1)
	}
	return C.int(git_SUCCESS)
}

type ConfigMultivarCallback func(value string, data interface{}) error

type cfgMultivarCallbackWrapper struct {
	f ConfigMultivarCallback
	d interface{}
}

func (cfg *Config) SetMultivar(name, regexp, value string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cregexp := C.CString(regexp)
	defer C.free(unsafe.Pointer(cregexp))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	ecode := C.git_config_set_multivar(cfg.git_config, cname, cregexp, cvalue)
	if ecode != git_SUCCESS {
		return gitError()
	}
	return nil
}

type ConfigFile struct {
	git_config_file *C.git_config_file
}

func NewConfigFile(path string) (*ConfigFile, error) {
	cfgFile := new(ConfigFile)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	// libgit2 awkwardly uses a struct here, we cast it later
	var cfgFileStruct *C.struct_git_config_file
	ecode := C.git_config_file__ondisk(&cfgFileStruct, cpath)
	if ecode != git_SUCCESS {
		return nil, gitError()
	}
	cfgFile.git_config_file = (*C.git_config_file)(cfgFileStruct)
	return cfgFile, nil
}
