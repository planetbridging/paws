package main


import (
    "fmt"
    "log"
    "net/http"
    "flag"
    
	"image/jpeg"
    "bytes"
	"image"
	"strconv"

)

var root = flag.String("root", ".", "file system path")

func start_hosting() {
    http.HandleFunc("/", handler)
	http.Handle("/web/", http.FileServer(http.Dir(*root)))
	log.Println("Serving at localhost:9123...")
	log.Fatal(http.ListenAndServe(":9123", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page: ", r.URL.Path)

    if r.URL.Path == "/devices"{
		fmt.Fprint(w,get_devices())
        //get_devices()
	}

}

func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		//log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		//log.Println("unable to write image.")
	}
}


func get_devices() string{
    tbl := ""
    for d := range lst_devices {
        //tbl += lst_devices[d].name
        num := strconv.Itoa(d)
        tbl += `<div class='dropdown-divider'></div>
        <a class="dropdown-item" onclick="Select_Device('/dev/`+num+`')">
        <table class='table table-dark'>
        <tbody>
          <tr>
            <td></td>
            <td>Name:</td>
            <td>`+lst_devices[d].name+`</td>
          </tr>
          <tr>
            <td></td>
            <td>description:</td>
            <td>`+lst_devices[d].description+`</td>
          </tr>`
        
        for a := range lst_devices[d].addresses {
            tbl += `<tr>
            <td></td>
            <td>ip:</td>
            <td>`+lst_devices[d].addresses[a].ip.String()+`</td>
            </tr>`
        }

        tbl += `</tbody>
        </table></a>`
    }
    return tbl
}