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
}

type Output struct {
	Image interface{} `md:"image"`
	FPS   float64     `md:"fps"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"image": o.Image,
		"fps":   o.FPS,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Image = values["image"] //, err = coerce.ToBytes(values["image"])
	// if err != nil {
	// 	return err
	// }
	o.FPS, err = coerce.ToFloat64(values["fps"])
	if err != nil {
		return err
	}
	return nil
}
