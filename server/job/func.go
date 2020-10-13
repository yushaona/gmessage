package job

import (
	"fmt"
	"strconv"

	"github.com/yushaona/gjson"
)

type CommonInterface interface {
	DoJob(param *gjson.GJSON) (result gjson.GJSON, err error)
}

var (
	mapfunc map[int]CommonInterface
)

func init() {
	mapfunc = make(map[int]CommonInterface)
}

func AddMap(id int, inter CommonInterface) {
	mapfunc[id] = inter
}

func DoJob(id int, param *gjson.GJSON) (result gjson.GJSON, err error) {
	if f, isok := mapfunc[id]; isok {
		return f.DoJob(param) // 这个地方可以优化,类似于我原来做的的那个通过reflect,找到函数,然后调用的方式
	}
	return result, fmt.Errorf("%s", strconv.Itoa(id)+" method not exists")
}

type BaseStruct struct {
}

func (t *BaseStruct) DoJob(param *gjson.GJSON) (result gjson.GJSON, err error) {
	fmt.Println("BaseStruct")

	return
}
