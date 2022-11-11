/*
Copyright © 2022 Lanly

*/
package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	mapset "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
)

func ResolvePath(path string) (Data, error) {
	// GD-00510/expr/expr.cpp
	s := strings.Split(path, string(os.PathSeparator))
	if len(s) != 3 {
		return Data{}, fmt.Errorf("Invalid len[%d] path: %s", len(s), path)
	}
	fileSuffix := filepath.Ext(s[2])
	fileNameOnly := strings.TrimSuffix(s[2], fileSuffix)
	data := Data{
		Path:      path,
		Name:      s[0],
		Problem:   fileNameOnly,
		Extension: fileSuffix,
		FileName:  s[2],
	}

	if s[1] != data.Problem {
		return data, fmt.Errorf("Inconsistent folder names and file names: %s", path)
	}

	return data, nil
}

func Md5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Copy(src string, dest string) error {
	fileInfo, err := os.Stat(src)
	if err != nil {
		log.Error(err)
		return err
	}

	if fileInfo.IsDir() {
		if err = os.MkdirAll(dest, 0775); err != nil {
			log.Error(err)
			return err
		}
		return nil
	}
	parentDir := filepath.Dir(dest)
	if _, err = os.Stat(parentDir); err != nil {
		log.Error(err)
		log.Warn("Well, it shouldn't reach here")
		if err = os.MkdirAll(parentDir, 0775); err != nil {
			log.Error(err)
			return err
		}
	}

	srcFile, err := os.Open(src)
	if err != nil {
		log.Error(err)
		return err
	}

	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)

	if err != nil {
		log.Error(err)
		return err
	}

	defer srcFile.Close()
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	if err != nil {
		log.Error(err)
		return err
	}

	stat, err := srcFile.Stat()
	if err != nil {
		log.Error(err)
		return err
	}
	err = os.Chtimes(dest, stat.ModTime(), stat.ModTime())
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func WriteFile(fileName string, data string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(data)
	write.Flush()
	return nil
}

func ReadFile(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func FileSize(fileName string) (int, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return 0, err
	}
	return int(fileInfo.Size()), nil
}

func SetToOrderStringSlice(set mapset.Set) []string {
	var data []string

	for _, tmp := range set.ToSlice() {
		data = append(data, fmt.Sprint(tmp))
	}

	sort.Strings(data)

	return data
}
