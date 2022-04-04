package models

type AhArgs struct {
	NeedScreenshot int    `json:"needScreenshot"`
	NodeName       string `json:"nodeName"`
	InstanceId     int    `json:"instanceId"`
}
