package main

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/voxelbrain/goptions"
)

var (
	options = struct {
		Server string        `goptions:"-s, --server, description='Pixelpixel server to push to'"`
		Help   goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		Server: "localhost:8080",
	}
)

const (
	keyFile = ".key"
)

func main() {
	goptions.ParseAndFail(&options)

	key := loadKey()

	fs, err := createFs(".")
	if err != nil {
		log.Fatalf("Could not create filesystem: %s", err)
	}

	var code int
	var body string
	if key != "" {
		log.Printf("Found Pixel ID %s", key)
		code, body, err = makeApiCall("PUT", "/"+key, bytes.NewReader(fs))
	} else {
		log.Printf("Creating a new pixel")
		code, body, err = makeApiCall("POST", "/", bytes.NewReader(fs))
	}

	if code >= 300 || err != nil {
		log.Printf("%s", string(body))
		log.Fatalf("Could not upload new pixel. Status code: %d, Error: %s", code, err)
	}

	f, err := os.Create(keyFile)
	if err != nil {
		log.Fatalf("Could not save key file (key=%s): %s", body, err)
	}
	defer f.Close()
	io.WriteString(f, body)
}

func makeApiCall(method, subpath string, body io.Reader) (int, string, error) {
	url := path.Join(options.Server, "pixels") + "/"
	url = slashify(url + subpath)
	log.Printf("URL: %s", url)
	req, _ := http.NewRequest(method, url, body)
	resp, err := http.DefaultClient.Do(req)
	buf := &bytes.Buffer{}

	code := -1
	if resp != nil {
		defer resp.Body.Close()
		code = resp.StatusCode
		io.Copy(buf, resp.Body)
	}

	// TODO: Sadly, this seems to be necassary. Can't seem to get it
	// working without it. I always get `malformed HTTP response ""`
	http.DefaultTransport.(*http.Transport).CloseIdleConnections()

	return code, buf.String(), err
}

func createFs(path string) ([]byte, error) {
	buf := &bytes.Buffer{}
	fs := tar.NewWriter(buf)
	err := filepath.Walk(path, func(path string, info os.FileInfo, _ error) error {
		if strings.HasPrefix(path, ".") {
			return nil
		}
		log.Printf("Adding %s", path)
		if info.IsDir() {
			err := fs.WriteHeader(&tar.Header{
				Name:     info.Name(),
				Typeflag: tar.TypeDir,
			})
			return err
		}

		err := fs.WriteHeader(&tar.Header{
			Name:     info.Name(),
			Mode:     int64(info.Mode()),
			Size:     info.Size(),
			Typeflag: tar.TypeReg,
		})
		if err != nil {
			log.Printf("Could not start next file: %s", err)
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			log.Printf("Could not read %s: %s", path, err)
			return err
		}
		defer f.Close()
		io.Copy(fs, f)
		return nil
	})
	fs.Flush()
	fs.Close()
	return buf.Bytes(), err
}

func loadKey() string {
	f, err := os.Open(keyFile)
	if err != nil {
		return ""
	}
	defer f.Close()
	key, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(key))
}

func slashify(url string) string {
	url = path.Clean(url)
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return "http://" + strings.TrimPrefix(url, "/")
}
