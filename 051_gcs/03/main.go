package skyhdd

import (
	"net/http"
"google.golang.org/appengine/log"
	"google.golang.org/appengine"
	"strconv"
	"io"
)

const gcsBucket = "learning-1130.appspot.com"

func init() {
	http.HandleFunc("/", handler)
}

func handler(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	html := `
	    <form method="POST" enctype="multipart/form-data">
		<input type="file" name="dahui">
		<input type="submit">
	    </form>
	`

	if req.Method == "POST" {

		mpf, hdr, err := req.FormFile("dahui")
		if err != nil {
			log.Errorf(ctx, "ERROR handler req.FormFile: ", err)
			http.Error(res, "We were unable to upload your file\n", http.StatusInternalServerError)
		}
		defer mpf.Close()


		fname, err := uploadFile(req, mpf, hdr)
		if err != nil {
			log.Errorf(ctx, "ERROR handler uploadFile: ", err)
			http.Error(res, "We were unable to accept your file\n" + err.Error(), http.StatusUnsupportedMediaType)
		}

		fnames, err := putCookie(res, req, fname)
		if err != nil {
			log.Errorf(ctx, "ERROR handler putCookie: ", err)
			http.Error(res, "We were unable to accept your file\n" + err.Error(), http.StatusUnsupportedMediaType)
		}

		html += `<h1>Files</h1>`
		for i, v := range fnames {
			html += `<h3>` + strconv.Itoa(i) + ` - ` + v + `</h3>`
		}
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(res, html)
}
