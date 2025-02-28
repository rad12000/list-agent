package file

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func DecompressTAR(filename string) (files []string, err error) {
	var closers []io.Closer
	defer func() {
		for _, c := range closers {
			if err := c.Close(); err != nil {
				fmt.Println(err)
			}
		}
	}()

	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	closers = append(closers, file)

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	closers = append(closers, gzipReader)

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate through the files in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// If the file is a directory, create it
		if header.Typeflag == tar.TypeDir {
			if err = os.MkdirAll(header.Name, 0755); err != nil {
				return nil, err
			}
			continue
		}

		file, err := os.Create(header.Name)
		if err != nil {
			return nil, err
		}
		files = append(files, file.Name())
		closers = append(closers, file)

		// Copy the contents of the file from the tar archive to the new file
		if _, err := io.Copy(file, tarReader); err != nil {
			return nil, err
		}
	}

	return
}
