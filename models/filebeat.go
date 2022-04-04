package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var Filebeatlines []string

func FilebeatModify(nodeName string) error{

	file, err := os.Open("filebeat-template.yml")
	if err != nil {
		return err
	}

	defer file.Close()
	if  err != nil{
		return err
	}
	sc := bufio.NewScanner(file)
	Filebeatlines = make([]string, 0)
	
	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		line,err := filebeatTModify(sc.Text(),nodeName)
		if err != nil {
			return err
		}
		Filebeatlines = append(Filebeatlines, line)
	}

	err = filebeatCreate()
	if err != nil{
		return err
	}

	return nil
}

func filebeatCreate() error {
	file, err := os.Create("filebeat/filebeat.yml")
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range Filebeatlines {
        fmt.Fprintln(w, line)
    }
    w.Flush()
	return nil
}

func filebeatTModify(line string,nodeName string) (string,error) {
	if strings.Contains(line,"redroid_node_name"){
		return strings.ReplaceAll(line,"redroid_node_name",strings.ToLower(string(nodeName))),nil
	}
	return line,nil
}