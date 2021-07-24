package ipcam

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
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"image": o.Image,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	//var err error
	o.Image = values["image"] //, err = coerce.ToBytes(values["image"])
	// if err != nil {
	// 	return err
	// }
	return nil
}
