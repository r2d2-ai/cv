package ipcam

import "github.com/r2d2-ai/core/data/coerce"

type Settings struct {
}

type HandlerSettings struct {
	Protocol string `md:"protocol,required,allowed(RSTP,ONVIF)"`
	Host     string `md:"host,required"`
	User     string `md:"user"`
	Password string `md:"password"`
	VideoURI string `md:"videoUri"`
	GroupId  string `md:"groupId"`
	CameraId string `md:"cameraId"`
}

type Output struct {
	Image    interface{} `md:"image"`
	GroupdId string      `md:"groupId"`
	CameraId string      `md:"cameraId"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"image":    o.Image,
		"groupId":  o.GroupdId,
		"cameraId": o.CameraId,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Image = values["image"]

	o.GroupdId, err = coerce.ToString(values["groupId"])
	if err != nil {
		return err
	}

	o.CameraId, err = coerce.ToString(values["cameraId"])
	if err != nil {
		return err
	}

	return nil
}
