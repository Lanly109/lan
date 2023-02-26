/*
Copyright © 2022 Lanly

*/
package utils

import (
	"fmt"
	"time"
)

const TimeTemplate string = "2006-01-02 15:04:05"

type Data struct {
	Path       string
	Name       string
	Problem    string
	Extension  string
	FileName   string
	Md5        string
	Size       int64
	ModifyTime time.Time
}

func (s *Data) String() string {
	return fmt.Sprintf("[%s %s]: %s %.2fkb", s.Name, s.FileName, s.ModifyTime.Format(TimeTemplate), float64(s.Size) / 1024)
}

func (s *Data) Output() string {
	return fmt.Sprintf("%s的%s修改时间为%s，大小为%.2fkb;", s.Name, s.FileName, s.ModifyTime.Format(TimeTemplate), float64(s.Size) / 1024)
}
