package models

import (
	"fmt"
	"strings"
)

func RcloneChcek(nodeName string) {

	cmd, err := Cmd("[ -d '/oracle/screenshot/' ] && echo '1'")
	if err != nil {
		fmt.Printf("[%s] rclone check error\n", nodeName)
		fmt.Println(err.Error())
	}
	cmd = strings.Replace(cmd, "\n", "", -1)
	if cmd != "1" {
		fmt.Printf("[%s] rclone need to remount \n", nodeName)
		Cmd("sudo umount /oracle")
		_, err = Cmd("nohup rclone mount --allow-other --buffer-size 512m --dir-cache-time 72h --drive-chunk-size 128M --umask 002 --vfs-read-chunk-size 512M --vfs-read-chunk-size-limit off --daemon --use-mmap oracle:/ /oracle/ >> /root/rclonelog.log 2>&1 &")
		if err != nil {
			fmt.Printf("[%s] rclone start error \n", nodeName)
		}
	}

}
