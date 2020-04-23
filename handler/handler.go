package handler

import (
	"CloudDisk/meta"
	"CloudDisk/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// UploadHandler handle upload requests
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		const filename = "static/view/index.html"
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("read file %s failed, err: %s", filename, err.Error())
			io.WriteString(w, "server internal error")
		}

		w.Write(data)
	} else if r.Method == "POST" {
		src, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("get data failed, err: %s\n", err.Error())
			return
		}

		defer src.Close()

		fmeta := meta.FileMeta{
			Name:     header.Filename,
			Path:     "/tmp/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		dst, err := os.Create(fmeta.Path)
		if err != nil {
			fmt.Printf("create file %s failed, err: %s\n", fmeta.Path, err.Error())
			return
		}

		defer dst.Close()

		fmeta.Size, err = io.Copy(dst, src)
		if err != nil {
			fmt.Printf("write data into file %s failed, err: %s\n", fmeta.Path, err.Error())
			return
		}

		dst.Seek(0, 0)
		fmeta.Sha1 = util.FileSha1(dst)
		fmt.Println(fmeta.Sha1)
		// meta.Set(fmeta)
		_ = meta.SetToDB(fmeta)
		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

//UploadSuccHandler reply upload success rsp
func UploadSuccHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload success!")
}

//GetFileMetaHandler get file meta
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("parse form failed, err: %s\n", err.Error())
		return
	}

	hash := r.Form.Get("filehash")
	fmeta, ok := meta.Get(hash)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, _ := json.Marshal(fmeta)
	w.Write(data)
}

//FileDownloadHandler download file handler
func FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("parse form failed, err: %s\n", err.Error())
		return
	}

	hash := r.Form["filehash"][0]
	fmeta, ok := meta.Get(hash)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := os.Open(fmeta.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+fmeta.Name+"\"")
	w.Write(data)
}

//FileMetaUpdateHandler update file meta
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hash := r.Form.Get("filehash")
	fmeta, ok := meta.Get(hash)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filename := r.Form.Get("filename")
	fmeta.Name = filename
	meta.Set(fmeta)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "update success")
}

//FileDeleteHandler delete file
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hash := r.Form.Get("filehash")
	fmeta, ok := meta.Get(hash)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	meta.Delete(hash)
	os.Remove(fmeta.Path)
	w.WriteHeader(http.StatusOK)
}
