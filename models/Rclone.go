package models

import (
	"fmt"
	"strings"
)

func RcloneChcek(nodeName string) {

	cmd, err := Cmd("[ -d '/oracle/screenshot/' ] && echo '1'")
	if err != nil {
		Post(fmt.Sprintf("[%s] rclone check error", nodeName))
		Post(err.Error())
	}
	cmd = strings.Replace(cmd, "\n", "", -1)
	if cmd != "1" {
		Post(fmt.Sprintf("[%s] rclone need to remount", nodeName))
		Cmd("sudo umount /oracle")
		_, err = Cmd("nohup rclone mount --allow-other --buffer-size 512m --dir-cache-time 72h --drive-chunk-size 128M --umask 002 --vfs-read-chunk-size 512M --vfs-read-chunk-size-limit off --daemon --use-mmap oracle:/ /oracle/ >> /root/rclonelog.log 2>&1 &")
		if err != nil {
			Post(fmt.Sprintf("[%s] rclone start error", nodeName))
			Post(err.Error())
		}
	}

}
