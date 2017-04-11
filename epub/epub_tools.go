package epub

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
)

// MoveEpub move epub file to special dir, and relative url path will be reset also
func MoveEpub(destPath string, epub *Epub) error {
	destFullPath := path.Join(destPath, epub.FileDir)
	if isFileExist(destFullPath) == true {
		err := os.RemoveAll(destFullPath)
		if err != nil {
			return err
		}
	}
	err := os.Rename(epub.FileDir, destFullPath)
	if err != nil {
		return err
	}

	err = epub.Navigations.Each(func(nav *NavigationPoint) error {
		nav.URL = path.Join(destPath, nav.URL)
		return nil
	})
	if err != nil {
		return err
	}

	if len(epub.MetaInfo["coverage"]) > 0 {
		epub.MetaInfo["coverage"] = path.Join(destPath, epub.MetaInfo["coverage"])
	}

	return nil
}

func caculateMD5Value(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// exists returns whether the given file or directory exists or not
func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
