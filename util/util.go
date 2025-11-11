package util

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (s *Sha1Stream) Update(data []byte) {
	if s._sha1 == nil {
		s._sha1 = sha1.New()
	}
	s._sha1.Write(data)
}

func (s *Sha1Stream) Sum() string {
	return hex.EncodeToString(s._sha1.Sum(nil))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum(nil))
}
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func Md5(data []byte) string {
	_md5 := sha1.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum(nil))
}
func FMd5Sha1(file *os.File) string {
	_md5 := sha1.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		result = info.Size()
		return nil
	})
	return result
}
