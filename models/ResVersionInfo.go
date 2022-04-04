package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"strings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ResVersionColl *mgm.Collection
var resVersion *ResVersionInfo

type ResVersionInfo struct {
	mgm.DefaultModel `bson:",inline"` //MongoDbModel
	ResVersion       string           `json:"resVersion"`
	ClientVersion    string           `json:"clientVersion"`
}

func ArkHookJSCheck(nodeName string,number int) {
	//	~/redroid/arkhookVersionCheck.info
	cmd, err := Cmd("/root/redroid/helper.sh arkhookCheck")
	if err != nil {
		Post(fmt.Sprintf("[%s] 检测arkhook.js 版本出错",nodeName))
		Post(err.Error())
		os.Exit(-1)
	}
	fmt.Println(string(cmd))

	oldArray, err := os.ReadFile("arkhookVersion.info")
	if err != nil {
		Post(fmt.Sprintf("[%s] 读取文件 arkhookVersion.info 出错",nodeName))
		Post(err.Error())
		os.Exit(-1)
	}
	newArray, err := os.ReadFile("arkhookVersionCheck.info")
	if err != nil {
		Post(fmt.Sprintf("[%s] 读取文件 arkhookVersionCheck.info 出错",nodeName))
		Post(err.Error())
		os.Exit(-1)
	}
	if string(oldArray) != string(newArray) {
		_, err := Cmd("/root/redroid/helper.sh arkhook")
		if err != nil {
			Post(fmt.Sprintf("[%s] 更新 arkhook.js 文件出错 ",nodeName))
			Post(err.Error())
			os.Exit(-1)
		}
		err = StopGame(number)
		if err != nil {
			Post(fmt.Sprintf("[%s] 结束游戏进程出错 : %s ",nodeName,err.Error()) )
			os.Exit(-1)
		}
	}

}

func ResVersionInit() {

	err := mgm.SetDefaultConfig(nil, "closure", options.Client().ApplyURI(ConnectStr))
	if err != nil {
		Post("MongoDB init failed in ArkResourceAutoUpdateBot")
		os.Exit(-1)
	}
	if !isFileExist("version.info"){
		resVersion = &ResVersionInfo{}
		ResVersionColl = mgm.Coll(resVersion)
		err = ResVersionColl.First(bson.M{}, resVersion)
		if err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				ResVersionColl.Create(resVersion)
			}
		}
		WriteResVersionFile(resVersion)
	}
}

func ResVersionCheck(nodeName string,number int) bool {

	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "closure", options.Client().ApplyURI(ConnectStr))
	if err != nil {
		Post("MongoDB init failed in ArkResourceAutoUpdateBot")
		os.Exit(-1)
	}

	if !isFileExist("version.info") {
		Post("[ArkResourceAutoUpdateBot] version.info 不存在!")
		os.Exit(-1)
	}

	resVersion = &ResVersionInfo{}
	ResVersionColl = mgm.Coll(resVersion)
	err = ResVersionColl.First(bson.M{}, resVersion)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			ResVersionColl.Create(resVersion)
		}
	}

	var fileVersion = ReadVersionFile()
	if fileVersion.ResVersion != resVersion.ResVersion {
		Post(fmt.Sprintf("[%s] 检测到热更新,将执行热更新shell脚本",nodeName))
		err = StopGame(number)
		if err != nil {
			Post(fmt.Sprintf("[%s] 错误 : %s ",nodeName,err.Error()) )
			os.Exit(-1)
		}
		_, err := Cmd("/root/redroid/helper.sh hotupdate")
		if err != nil {
			Post(fmt.Sprintf("[%s] 错误 : %s ",nodeName,err.Error()) )
			os.Exit(-1)
		}
		WriteResVersionFile(resVersion)
		return true
	}

	return false
}

func ResVersionUpdate() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "closure", options.Client().ApplyURI(ConnectStr))
	if err != nil {
		Post("MongoDB init failed in ArkResourceAutoUpdateBot")
		os.Exit(-1)
	}

	if !isFileExist("version.info") {
		Post("[ArkResourceAutoUpdateBot] version.info 不存在!")
		os.Exit(-1)
	}

	resVersion = &ResVersionInfo{}
	ResVersionColl = mgm.Coll(resVersion)
	err = ResVersionColl.First(bson.M{}, resVersion)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			ResVersionColl.Create(resVersion)
		}
	}
	var fileVersion = ReadVersionFile()
	if fileVersion.ResVersion != resVersion.ResVersion {
		resVersion.ResVersion = fileVersion.ResVersion
		resVersion.ClientVersion = fileVersion.ClientVersion
		err = ResVersionColl.Update(resVersion)
		if err != nil {
			Post("MongoDB update failed in Update Function(ArkResourceAutoUpdateBot)")
			os.Exit(-1)
		}
		Post("[ArkResourceAutoUpdateBot] version.info 与 mongoDB 同步成功")
		Post(fmt.Sprintf("[ArkResourceAutoUpdateBot] ResVersion : %s", fileVersion.ResVersion))
	}
}

func ReadVersionFile() *ResVersionInfo {
	jsonFile, err := os.Open("version.info")

	// 最好要处理以下错误
	if err != nil {
		fmt.Println(err)
	}

	// 要记得关闭
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		os.Exit(-1)
	}
	var res = &ResVersionInfo{}
	err = json.Unmarshal([]byte(byteValue), res)
	if err != nil {
		os.Exit(-1)
	}
	return res
}
