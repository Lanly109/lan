/*
Copyright © 2022 Lanly

*/
package utils

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	server string = "moss.stanford.edu"
	port   int    = 7690
)

type Moss struct {
	UserId       string
	Language     string
	MaxLimit     int64
	Comment      string
	Experimental int
	NumberResult int64
	Files        []string
}

func (moss *Moss) AddFile(file string) {
	moss.Files = append(moss.Files, file)
}

func (moss *Moss) uploadFile(session net.Conn, file string, index int) error {
	name := strings.ReplaceAll(strings.ReplaceAll(file, " ", "_"), "\\", "/")

	fileSize, err := FileSize(file)
	if err != nil {
		return err
	}

	fmt.Fprintf(session, "file %d %s %d %s\n", index, moss.Language, fileSize, name)

	fileData, err := ReadFile(file)
	if err != nil {
		return err
	}
	fmt.Fprint(session, fileData)
	return nil
}

func (moss *Moss) Review() string {
	session, err := net.Dial("tcp", fmt.Sprintf("%v:%v", server, port))

	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()

	fmt.Fprintf(session, "moss %s\n", moss.UserId)
	fmt.Fprintf(session, "directory 0\n")
	fmt.Fprintf(session, "X %d\n", moss.Experimental)
	fmt.Fprintf(session, "maxmatches %d\n", moss.MaxLimit)
	fmt.Fprintf(session, "show %d\n", moss.NumberResult)
	fmt.Fprintf(session, "language %s\n", moss.Language)

	response, err := bufio.NewReader(session).ReadString('\n')

	if err != nil || strings.TrimSpace(response) == "no" {
		if err != nil {
			log.Fatal(err)
		} else {
			log.Fatalf("Unrecognized language %s", moss.Language)
		}
	}

	for index, file := range moss.Files {
		err := moss.uploadFile(session, file, index+1)
		if err != nil {
			log.Errorf("upload [%s] failed: %s", file, err)
		} else {
			log.Infof("uploaded file %s", file)
		}
	}

	fmt.Fprintf(session, "query 0 %s\n", moss.Comment)

	log.Info("Query submitted.  Waiting for the server's response.")

	response, err = bufio.NewReader(session).ReadString('\n')

	if err != nil {
		log.Fatal(err)
	} else if strings.HasPrefix(response, "Error") {
		log.Fatal(response)
	}

	return response
}
