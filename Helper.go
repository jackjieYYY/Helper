package main

import (
	"Helper/models"
	"fmt"
	"os"
	"strconv"
)

func main() {

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			models.UnLock()
			models.Post("[Go] Helper 捕获到异常: " + err.(string))
		}
	}()

	if len(os.Args) == 1 {
		os.Exit(-1)
	}

	// usage : ./Helper ResVersionUpdate
	// for ArkResourceAutoUpdateBot
	if os.Args[1] == "ResVersionUpdate" {
		if len(os.Args) != 2 {
			os.Exit(-1)
		}
		models.ResVersionUpdate()
	}

	// usage : ./Helper "nodereport","test","1"
	if os.Args[1] == "nodereport" {
		if len(os.Args) != 4 {
			os.Exit(-1)
		}
		var number int
		number, err := strconv.Atoi(os.Args[3])
		if err != nil {
			os.Exit(-1)
		}
		models.NodeReportStart(os.Args[2], number)
	}

	// usage : ./Helper "fileModify","nodename", 4","90","45","20"
	if os.Args[1] == "fileModify" {
		//   sys.argv[1] = type
		//   sys.argv[2] = nodename
		//   sys.argv[3] = number
		//   sys.argv[4] = redroid.width
		//   sys.argv[5] = redroid.height
		//   sys.argv[6] = redroid.fps
		if len(os.Args) != 7 {
			os.Exit(-1)
		}

		nodeName := os.Args[2]
		number, err := strconv.Atoi(os.Args[3])
		if err != nil {
			os.Exit(-1)
		}
		width, err := strconv.Atoi(os.Args[4])
		if err != nil {
			os.Exit(-1)
		}
		height, err := strconv.Atoi(os.Args[5])
		if err != nil {
			os.Exit(-1)
		}
		fps, err := strconv.Atoi(os.Args[6])
		if err != nil {
			os.Exit(-1)
		}
		models.FilebeatModify(nodeName)
		models.Modify(number, width, height, fps)
	}

	if os.Args[1] == "init" {
		//   sys.argv[1] = type
		//   sys.argv[2] = number
		if len(os.Args) != 3 {
			os.Exit(-1)
		}
		var number int
		number, err := strconv.Atoi(os.Args[2])
		if err != nil {
			os.Exit(-1)
		}
		models.ResVersionInit()
		models.RedroidInit(number)
	}
	//	usage : ./Helper	surveillance,	test,		1,			1			4
	//	usage : ./Helper	type,			nodeName,	nodeNumber,	taskType,	number
	//	sys.argv[1] = type
	//	sys.argv[2] = nodeName
	//	sys.argv[3] = nodeNumber
	//	sys.argv[4] = taskType
	//	sys.argv[5] = number
	if os.Args[1] == "surveillance" {
		fmt.Println("surveillance start ")
		if len(os.Args) != 6 || models.SurveillanceIsLock() {
			fmt.Println(os.Args)
			fmt.Println(models.SurveillanceIsLock())
			fmt.Println("len(os.Args) != 6 || models.IsLock() ")
			os.Exit(-1)
		}

		models.Lock()

		nodeName := os.Args[2]
		nodeNumber, err := strconv.Atoi(os.Args[3])
		if err != nil {
			os.Exit(-1)
		}
		needScreenshot, err := strconv.Atoi(os.Args[4])
		if err != nil {
			os.Exit(-1)
		}
		if needScreenshot < 0 || needScreenshot > 1 {
			models.Post("截图参数错误,needScreenshot 应为1 (需要截图) 或 0 (不需要截图)")
			os.Exit(-1)
		}

		number, err := strconv.Atoi(os.Args[5])
		if err != nil {
			os.Exit(-1)
		}
		models.NodeReportStart(nodeName, nodeNumber)
		models.ResVersionCheck(nodeName, number)
		models.ArkHookJSCheck(nodeName, number)
		models.Surveillance(number, needScreenshot, nodeName)
		models.UnLock()
	}
}
