package ipcam

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/r2d2-ai/core/action"
	"github.com/r2d2-ai/core/support/test"
	"github.com/r2d2-ai/core/trigger"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "flogo-cam",
	"ref": "github.com/r2d2-ai/cam/trigger/ipcam",
	"settings": {},
	"handlers": [
	  {
		"action": {
		  "id": "dummy"
		},
		"settings": {
		  "protocol": "RSTP",
		  "host": "192.168.50.195",
		  "user": "admin",
		  "password": "P$rolaMeaCAM",
		  "videoUri": "11",
		  "groupId": "group-1",
		  "cameraId": "cam-1"
		}
	  }
	]
  }`

func TestIPCamTrigger_Initializer(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)
	actions := map[string]action.Action{"dummy": test.NewDummyAction((func() {}))}

	trg, err := test.InitTrigger(f, config, actions)

	assert.Nil(t, err)
	assert.NotNil(t, trg)

	err = trg.Start()
	assert.Nil(t, err)
	time.Sleep(time.Second * 1)
	err = trg.Stop()
	assert.Nil(t, err)
}
