package main

import (
	"io"
	"log"
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
		log.Printf("upload %s: path invalid", path)
		return
	}

	file, _, err := r.FormFile("files")
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("upload form file %s error: path invalid", path)
		return
	}
	defer file.Close()

	fW, err := os.Create("." + path)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		log.Printf("upload %s: %s", path, err.Error())
		return
	}
	defer fW.Close()
	_, err = io.Copy(fW, file)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		log.Printf("upload %s: %s", path, err.Error())
		return
	}
	_, _ = w.Write([]byte("success"))
	log.Printf("upload %s success", path)
}
