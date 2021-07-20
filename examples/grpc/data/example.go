package data

import (
	"fmt"

	"github.com/mike955/zrpc/log"
)

type ExampleData struct {
	logger *log.Entry
	// dao    *dao.LayoutDao
}

func NewExampleData(logger *log.Entry) *ExampleData {
	return &ExampleData{
		logger: log.Helper(logger, map[string]interface{}{"data": "example"}),
		// dao:    dao.NewLayoutDao(logger),
	}
}

func (s *ExampleData) Hello(data string) (res string, err error) {
	s.logger.Infof("func:Hello|request:%+v", res)
	res = fmt.Sprintf("hello %s", data)
	return
}
