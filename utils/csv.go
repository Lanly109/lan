/*
Copyright © 2022 Lanly

*/
package utils

import (
	"encoding/csv"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ReadCsv(filepath string) [][]string {
	//打开文件(只读模式)，创建io.read接口实例
	fp, err := os.Open(filepath)
	if err != nil {
		log.Error(err)
		cobra.CheckErr(err)
	}
	defer fp.Close()

	//创建csv读取接口实例
	readCsv := csv.NewReader(fp)

	//读取所有内容
	data, err := readCsv.ReadAll()

	cobra.CheckErr(err)

	return data
}

func WriteCsv(filepath string, datas [][]string) error {
	//打开文件(只读模式)，创建io.read接口实例
	fp, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		cobra.CheckErr(err)
	}
	defer fp.Close()

	//创建csv读取接口实例
	writeCsv := csv.NewWriter(fp)

	err = writeCsv.WriteAll(datas)

	cobra.CheckErr(err)

	return nil
}
