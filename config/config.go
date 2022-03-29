package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// EnableFlagArgs 根据启动参数加载配置文件
func EnableFlagArgs(defaultPath string) {
	path := flag.String("c", defaultPath, "默认值")
	flag.Parse()
	defaultConfig = Load(*path)
}

func Load(filepath string) *MapConfig {
	cfg := new(MapConfig)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("\n配置文件读取失败：\n %s\n %s\n\n", filepath, err)
		os.Exit(-1)
	}
	yaml.Unmarshal(content, cfg)
	return cfg
}

// 公共方法
var defaultConfig *MapConfig

func Get(key string) interface{} {
	return defaultConfig.Get(key)
}
func GetStruct(key string, bind interface{}) error {
	return defaultConfig.GetStruct(key, bind)
}
func GetString(key string) string {
	return defaultConfig.GetString(key)
}
func GetStrings(key string) []string {
	return defaultConfig.GetStrings(key)
}
func GetInt(key string) int {
	return defaultConfig.GetInt(key)
}
func GetInt8(key string) int8 {
	return int8(defaultConfig.GetInt(key))
}
func GetInt16(key string) int16 {
	return int16(defaultConfig.GetInt(key))
}
func GetInt32(key string) int32 {
	return int32(defaultConfig.GetInt(key))
}
func GetInt64(key string) int64 {
	return defaultConfig.GetInt64(key)
}
func GetMap(key string) map[string]string {
	return defaultConfig.GetMap(key)
}
func GetBool(key string) bool {
	return defaultConfig.GetBool(key)
}
func GetBools(key string) []bool {
	return defaultConfig.GetBools(key)
}

type MapConfig map[string]interface{}

func (c *MapConfig) GetStruct(key string, bind interface{}) error {
	str, err := yaml.Marshal(c.Get(key))
	if err == nil {
		err = yaml.Unmarshal(str, bind)
	}
	return err
}
func (c *MapConfig) Get(key string) interface{} {
	if c == nil {
		fmt.Printf("\n未指定配置文件或指定配置文件不存在\n\n")
		os.Exit(-1)
	}
	keys := strings.Split(key, ".")
	var conf interface{}
	conf = *c
	for _, k := range keys {
		if temp, ok := conf.(MapConfig); ok {
			conf = temp[k]
		} else {
			return nil
		}
	}
	return conf
}
func (c *MapConfig) GetString(key string) string {
	value := c.Get(key)
	return InterfaceToString(value)
}
func (c *MapConfig) GetStrings(key string) (result []string) {
	if str, err := yaml.Marshal(c.Get(key)); err == nil {
		yaml.Unmarshal(str, &result)
	}
	return
}
func (c *MapConfig) GetInt(key string) int {
	value := c.Get(key)
	switch vt := value.(type) {
	case int:
		return vt
	default:
		t := fmt.Sprint(vt)
		tInt, _ := strconv.Atoi(t)
		return tInt
	}
}
func (c *MapConfig) GetInts(key string) (result []int) {
	if str, err := yaml.Marshal(c.Get(key)); err == nil {
		yaml.Unmarshal(str, &result)
	}
	return
}
func (c *MapConfig) GetInt8(key string) int8 {
	return int8(c.GetInt(key))
}
func (c *MapConfig) GetInt32(key string) int32 {
	return int32(c.GetInt(key))
}
func (c *MapConfig) GetInt64(key string) int64 {
	value := c.Get(key)
	switch vt := value.(type) {
	case int64:
		return vt
	default:
		t := fmt.Sprint(vt)
		tInt, _ := strconv.ParseInt(t, 10, 64)
		return tInt
	}
}
func (c *MapConfig) GetMap(key string) (result map[string]string) {
	if str, err := yaml.Marshal(c.Get(key)); err == nil {
		yaml.Unmarshal(str, &result)
	}
	return
}
func (c *MapConfig) GetBool(key string) bool {
	value := c.Get(key)
	switch vt := value.(type) {
	case bool:
		return vt
	case int, int8, int16, int64, uint, uint8, uint16, uint64:
		if vt == 0 {
			return false
		}
		return true
	default:
		t := fmt.Sprint(vt)
		if t == "true" {
			return true
		}
		return false
	}
}
func (c *MapConfig) GetBools(key string) (result []bool) {
	if str, err := yaml.Marshal(c.Get(key)); err == nil {
		yaml.Unmarshal(str, &result)
	}
	return
}

func InterfaceToString(i interface{}) string {
	switch vt := i.(type) {
	case int, int8, int16, int64, uint, uint8, uint16, uint64:
		return fmt.Sprintf("%d", vt)
	case float32, float64:
		return fmt.Sprintf("%f", vt)
	case []byte:
		return string(vt)
	case string:
		return vt
	default:
		return ""
	}
}
