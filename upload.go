package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	page := `
<!DOCTYPE html>
<html>
<head>
    <title>Upload Files</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no">
</head>
<body>
    <form id="form1" enctype="multipart/form-data" method="post" action="Upload.aspx">
        <div class="row">
            <label for="fileToUpload">Select a File to Upload</label><br />
            <input type="file" name="files" id="fileToUpload" onchange="fileSelected();"/>
        </div>
            <div id="fileName"></div>
            <div id="fileSize"></div>
            <div id="fileType"></div>
            <div class="row">
            <input type="button" onclick="uploadFile()" value="Upload" />
        </div>
        <div id="progressNumber"></div>
    </form>

    <script type="text/javascript">
    function fileSelected() {
        var file = document.getElementById('fileToUpload').files[0];
        if (file) {
            var fileSize = 0;
            if (file.size > 1024 * 1024)
                fileSize = (Math.round(file.size * 100 / (1024 * 1024)) / 100).toString() + 'MB';
            else fileSize = (Math.round(file.size * 100 / 1024) / 100).toString() + 'KB';
            document.getElementById('fileName').innerHTML = 'Name: ' + file.name;
            document.getElementById('fileSize').innerHTML = 'Size: ' + fileSize;
            document.getElementById('fileType').innerHTML = 'Type: ' + file.type;
        }
    }
    function uploadFile() {
        var fd = new FormData();
        fd.append("files", document.getElementById('fileToUpload').files[0]);
        var xhr = new XMLHttpRequest();
        xhr.upload.addEventListener("progress", uploadProgress, false);
        xhr.addEventListener("load", uploadComplete, false);
        xhr.addEventListener("error", uploadFailed, false);
        xhr.addEventListener("abort", uploadCanceled, false);
        xhr.open("POST", "#");
        xhr.send(fd);
    }
    function uploadProgress(evt) {
        if (evt.lengthComputable) {
            var percentComplete = Math.round(evt.loaded * 100 / evt.total);
            document.getElementById('progressNumber').innerHTML = percentComplete.toString() + '%';
        } else document.getElementById('progressNumber').innerHTML = 'unable to compute';
    }
    function uploadComplete(evt) {
        alert(evt.target.responseText);
    }
    function uploadFailed(evt) {
        alert("There was an error attempting to upload the file.");
    }
    function uploadCanceled(evt) {
        alert("The upload has been canceled by the user or the browser dropped the connection.");
    }
    </script>
</body>
</html>
`
	if r.Method == http.MethodGet {
		_, _ = io.WriteString(w, page)
		log.Printf("upload: get page")
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("upload: method error")
		return
	}

	path := r.URL.Path
	if !check(path) {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("upload: path invalid")
		return
	}

	err := r.ParseMultipartForm(1024000)
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("upload: %s", err.Error())
		return
	}

	files := r.MultipartForm.File["files"]
	log.Printf("upload: %s", string(len(files)))
	if len(files) > 0 {
		log.Printf("upload: %s", files[0].Filename)
	}
	result := true
	fileNames := ""
	for i, _ := range files {
		file := files[i]
		fileNames += " " + file.Filename
		err = _saveFile(w, file, path)
		if err != nil {
			result = false
			http.Error(w, "404", http.StatusNotFound)
			log.Printf("upload: %s", err.Error())
			break
		}
	}
	if result {
		_, _ = w.Write([]byte("success"))
		log.Printf("upload: success," + fileNames)
	}
}

func _saveFile(w http.ResponseWriter, fileHeader *multipart.FileHeader, path string) error {
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		return err
	}
	dst, err := os.Create("." + path + "/" + fileHeader.Filename)
	defer dst.Close()
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
