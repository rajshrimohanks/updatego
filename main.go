package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	// Flags
	version := flag.String("version", "1.15.5", "the go version to fetch")
	goInstallDirectory := flag.String("go-dir", "/usr/local", "the directory inside which the go archive should be extracted")
	flag.Parse()

	fileName := "go" + (*version) + ".linux-amd64.tar.gz"
	newVersionURL := "https://golang.org/dl/" + fileName
	downloadedFile := filepath.Join("/tmp", fileName)

	// Download file
	err := downloadFile(downloadedFile, newVersionURL)
	if err != nil {
		log.Fatalf("an error occurred while downloading file: %s", err.Error())
	}
	log.Printf("file downloaded: %s", downloadedFile)

	// Open downloaded file
	r, err := os.Open(downloadedFile)
	if err != nil {
		log.Fatalf("an error occurred while opening downloaded file: %s", err.Error())
	}

	// Delete existing go directory
	currentGoDir := (*goInstallDirectory) + "/go"
	err = os.RemoveAll(currentGoDir)
	if err != nil {
		log.Fatalf("an error occurred while deleting current go directory: %s", err.Error())
	}
	log.Printf("directory deleted: %s", currentGoDir)

	// Extract contents
	err = untar(*goInstallDirectory, r)
	if err != nil {
		log.Fatalf("an error occurred while extracting downloaded file to destination: %s", err.Error())
	}
	log.Printf("extracted downloaded file to location: %s", *goInstallDirectory)

	// Delete downloaded file
	err = os.Remove(downloadedFile)
	if err != nil {
		log.Fatalf("an error occurred while deleting the downloaded file: %s", err.Error())
	}
	log.Printf("file deleted: %s", downloadedFile)
}

// downloadFile fetches a file from url and puts it at path
// shameless copy from: https://stackoverflow.com/a/33853856/4050218
func downloadFile(path string, url string) (err error) {
	log.Printf("begin downloading file to path: %s", path)

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to the file
	size, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	log.Printf("written: %d bytes", size)

	return nil
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
// Originally from: https://gist.github.com/sdomino/635a5ed4f32c93aad131/ but with
// added sanitization for zip-slip vulnerability
func untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)
		// https://snyk.io/research/zip-slip-vulnerability
		if !strings.HasPrefix(target, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("Illegal file path: %s", target)
		}

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
