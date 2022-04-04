package models

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)
var lines []string
var Filelines []string

func dockerFileModify(line string,width int,height int,fps int) (string,error) {
	if strings.Contains(line,"redroidwidth"){
		return strings.ReplaceAll(line,"redroidwidth",strconv.Itoa(width)),nil
	}
	if strings.Contains(line,"redroidheight"){
		return strings.ReplaceAll(line,"redroidheight", strconv.Itoa(height)),nil
	}
	if strings.Contains(line,"redroidfps"){
		return strings.ReplaceAll(line,"redroidfps", strconv.Itoa(fps)),nil
	}
	if strings.Contains(line,"image-arch-replace"){
		if runtime.GOARCH ==  "arm64" {
			return strings.ReplaceAll(line,"image-arch-replace", "redroid/redroid:11.0.0-arm64"),nil
		}
		if runtime.GOARCH ==  "amd64" {
			return strings.ReplaceAll(line,"image-arch-replace", "redroid/redroid:11.0.0-amd64"), nil
		}
		fmt.Println(runtime.GOARCH)
		return "", errors.New("unknown runtime") 
	}
	return line,nil

}

func dockerFileCreate(number int) error {

	Filelines := make([]string, 0)
	Filelines = append(Filelines, "services:")
	
	i := 0
	for i < number {
		for _, val := range lines {
			if strings.Contains(val,"redroid_id"){
				val = strings.ReplaceAll(val,"redroid_id","redroid" + strconv.Itoa(i + 1))
			}
			if strings.Contains(val,"redroidport"){
				val = strings.ReplaceAll(val,"redroidport",strconv.Itoa(5555 + i))
			}
			if strings.Contains(val,"_id"){
				val = strings.ReplaceAll(val,"_id",strconv.Itoa(i + 1))
			}
			Filelines = append(Filelines, val)
		}
		i = i + 1
	}

	file, err := os.Create("docker-compose.yml")
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range Filelines {
        fmt.Fprintln(w, line)
    }
    w.Flush()
	return nil
}


func Modify(number int,width int,height int,fps int) error{

	file, err := os.Open("docker-compose-template.yml")
	if err != nil {
		return err
	}

	defer file.Close()
	if  err != nil{
		return err
	}
	sc := bufio.NewScanner(file)
	lines = make([]string, 0)
	
	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		line,err := dockerFileModify(sc.Text(),width,height,fps)
		if err != nil {
			return err
		}
		lines = append(lines, line)
	}
	
	err = dockerFileCreate(number)
	if err != nil{
		return err
	}


	return nil
}

