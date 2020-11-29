package configuration

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/lancelotXie/proxy/proxy.lib/setting"
	"github.com/pkg/errors"
)

type configuration struct {
	l sync.Mutex
	setting.Setting
}

func newConfiguration() (c *configuration) {
	c = new(configuration)
	c.Setting = setting.New()
	return
}

// Load 加载配置文件
func (c *configuration) Load(fileName string) (err error) {
	c.l.Lock()
	defer c.l.Unlock()

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = json.Unmarshal(data, &c.Setting)
	return
}

// Save 保存配置到文件
func (c *configuration) Save(fileName string) (err error) {
	c.l.Lock()
	defer c.l.Unlock()

	data, err := json.MarshalIndent(c.Setting, "", "    ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}
