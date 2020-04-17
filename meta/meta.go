package meta

//FileMeta file meta data
type FileMeta struct {
	Sha1     string
	Name     string
	Size     int64
	Path     string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

//Set set or update filemetas
func Set(meta FileMeta) {
	fileMetas[meta.Sha1] = meta
}

//Get get filemeta from sha1
func Get(sha1 string) (FileMeta, bool) {
	m, ok := fileMetas[sha1]
	return m, ok
}

//Delete delete filemeta by sha1
func Delete(sha1 string) {
	_, ok := fileMetas[sha1]
	if !ok {
		return
	}

	delete(fileMetas, sha1)
}
