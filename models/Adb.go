package models

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	adb "github.com/zach-klippenstein/goadb"
)

var redroidhost = "127.0.0.1"

func connectRedroid(number int) ([]*adb.Device, error) {
	var result = make([]*adb.Device, 0)
	client, err := adb.New()
	if err != nil {
		return nil, err
	}

	client.StartServer()
	for i := 0; i < number; i++ {
		fmt.Println("try to conncet adb host " + redroidhost)
		err = client.Connect(redroidhost, 5555+i)
		if err != nil {
			fmt.Println("Connect error")
			fmt.Println(err.Error())
			return nil, err
		}
	}

	serials, err := client.ListDeviceSerials()
	if err != nil {
		fmt.Println("ListDeviceSerials error")
		return nil, err
	}
	fmt.Println("try to get serials")
	for _, serial := range serials {
		if strings.Contains(serial, redroidhost) {
			device := client.Device(adb.DeviceWithSerial(serial))
			result = append(result, device)
		}
	}
	result = rev(result)
	fmt.Println(result)
	return result, nil
}

func RedroidInit(number int) error {
	var err error
	flag.Parse()

	devices, err := connectRedroid(number)
	if err != nil {
		return err
	}
	for _, device := range devices {
		state, err := device.State()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Println(state)
		//nohup netcat -s 127.0.0.1 -p 36889 -L /system/bin/sh > /dev/null 2>&1 &
		err = runRedroidCMD(device, "nohup netcat -s 127.0.0.1 -p 36889 -L /system/bin/sh > /dev/null 2>&1 &")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = runRedroidCMD(device, "cp /data/local/host /system/etc/")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = runRedroidCMD(device, "pm install /data/local/1.apk")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = runRedroidCMD(device, "settings put secure immersive_mode_confirmations confirmed")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = runRedroidCMD(device, "pm grant com.hypergryph.arknights android.permission.READ_PHONE_STATE")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = runRedroidCMD(device, "pm grant com.hypergryph.arknights android.permission.WRITE_EXTERNAL_STORAGE")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = runRedroidCMD(device, "chmod 777 /data/local/frida-inject")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = runRedroidCMD(device, "su 0 killall frida-inject")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func StopGame(number int) error {
	devices, err := connectRedroid(number)
	if err != nil {
		return err
	}
	for _, device := range devices {
		err = runRedroidCMD(device, "am force-stop com.hypergryph.arknights")
		if err != nil {
			Post(err.Error())
			os.Exit(-1)
		}
		time.Sleep(time.Second * 5)
	}
	return nil
}

func Ping(number int, needScreenshot int, nodeName string) error {
	devices, err := connectRedroid(number)
	if err != nil {
		return err
	}
	for i := 0; i < 6; i++ {
		for _, device := range devices {
			cmdOutput, err := device.RunCommand("exec nc 127.0.0.1 36891")
			Post("测试" + nodeName + "连通性 : " + cmdOutput)
			if strings.Contains(cmdOutput, "pong") {
				continue
			}
			if err != nil {
				Post("测试" + nodeName + "连通性 exec nc 127.0.0.1 36891 遭遇错误 " + err.Error())
				continue
			}
			//need to stop
			err = runRedroidCMD(device, "killall com.hypergryph.arknights")
			if err != nil {
				Post("测试" + nodeName + "连通性 killall com.hypergryph.arknights 遭遇错误 " + err.Error())
				continue
			}
			time.Sleep(time.Second * 5)
			//and restart
			Surveillance(number, needScreenshot, nodeName)

		}

	}
	return nil
}

func Surveillance(number int, needScreenshot int, nodeName string) error {

	for i := 0; i < 6; i++ {
		devices, err := connectRedroid(number)
		if err != nil {
			return err
		}
		for _, device := range devices {
			cmdOutput, err := device.RunCommand("pidof netcat")
			if err != nil {
				return err
			}
			if cmdOutput == "" {
				err = runRedroidCMD(device, "nohup netcat -s 127.0.0.1 -p 36889 -L /system/bin/sh > /dev/null 2>&1 &")
				if err != nil {
					Post(fmt.Sprintf("[%s] 启动netcat出错 ", nodeName))
					Post(err.Error())
					os.Exit(-1)
				}
			}
		}

		for _, device := range devices {
			cmdOutput, err := runRedroidCMD_OutPut(device, "timeout 2 sh -c 'exec nc 127.0.0.1 36891'")
			PostToMe("测试" + nodeName + "连通性 : " + cmdOutput)
			if strings.Contains(cmdOutput, "pong") {
				continue
			}
			if err != nil {
				Post("测试" + nodeName + "连通性 exec nc 127.0.0.1 36891 遭遇错误 " + err.Error())
				continue
			}
			//need to stop
			err = runRedroidCMD(device, "killall com.hypergryph.arknights")
			if err != nil {
				Post("测试" + nodeName + "连通性 killall com.hypergryph.arknights 遭遇错误 " + err.Error())
				continue
			}

		}
		time.Sleep(time.Second * 2)
		for index, device := range devices {
			cmdOutput, err := device.RunCommand("pidof com.hypergryph.arknights")
			if err != nil {
				return err
			}
			if cmdOutput == "" {
				ahargs := AhArgs{
					NodeName:       nodeName,
					InstanceId:     index + 1,
					NeedScreenshot: needScreenshot,
				}

				buf := new(bytes.Buffer)
				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(false)
				_ = enc.Encode(ahargs)

				if err != nil {
					fmt.Println(err)
					return err
				}
				_, err = device.RunCommand("mkdir -p /data/local/logs")
				if err != nil {
					fmt.Println(err)
					return err
				}

				cmdOutput, err = Cmd("/root/redroid/helper.sh GameStart " + strconv.Itoa(index+1) + " " + nodeName + " " + strconv.Itoa(needScreenshot))
				if err != nil {
					Post(fmt.Sprintf("[%s] 启动游戏出错 redroid ID : %s ", nodeName, strconv.Itoa(index+1)))
					Post(fmt.Sprintf("CMD命令返回 : %s ", cmdOutput))
					Post(fmt.Sprintf("err : %s ", err.Error()))
				}
				time.Sleep(time.Second * 10)
			}
		}

		time.Sleep(time.Second * 9)
	}

	return nil
}

func runRedroidCMD(device *adb.Device, cmd string) error {
	cmdOutput, err := device.RunCommand(cmd)
	fmt.Println(cmdOutput)
	if err != nil {
		return err
	}
	return nil
}

func runRedroidCMD_OutPut(device *adb.Device, cmd string) (string, error) {
	cmdOutput, err := device.RunCommand(cmd)
	fmt.Println(cmdOutput)
	if err != nil {
		return "", err
	}
	return cmdOutput, nil
}

func rev(slice []*adb.Device) []*adb.Device {
	fmt.Println(slice)
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
