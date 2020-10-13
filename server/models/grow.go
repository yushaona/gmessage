package models

import (
	"fmt"

	"github.com/yushaona/gjson"
	"github.com/yushaona/gmessage/server/job"
)

type Grow struct {
	job.BaseStruct
}

func (t *Grow) DoJob(param *gjson.GJSON) (result gjson.GJSON, err error) {
	fmt.Println("Grow")

	result.SetString("grow", "11111")
	return
}
