package model

import (
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

var Languages = map[string]string{
	"zh-CN": "简体中文",
	"zh-TW": "繁體中文",
	"en-US": "English",
	"es-ES": "Español",
}

var Themes = map[string]string{
	"default":       "Default",
	"daynight":      "JackieSung DayNight",
	"mdui":          "Neko Mdui",
	"hotaru":        "Hotaru",
	"angel-kanade":  "AngelKanade",
	"server-status": "SeverStatus",
	"custom":        "Custom(local)",
}

var DashboardThemes = map[string]string{
	"default": "Default",
	"custom":  "Custom(local)",
}

const (
	ConfigTypeGitHub  = "github"
	ConfigTypeGitee   = "gitee"
	ConfigTypeGitlab  = "gitlab"
	ConfigTypeJihulab = "jihulab"
	ConfigTypeGitea   = "gitea"
)

const (
	ConfigCoverAll = iota
	ConfigCoverIgnoreAll
)

type AgentConfig struct {
	HardDrivePartitionAllowlist []string
	NICAllowlist                map[string]bool
	v                           *viper.Viper
}

// Read 从给定的文件目录加载配置文件
func (c *AgentConfig) Read(path string) error {
	c.v = viper.New()
	c.v.SetConfigFile(path)
	err := c.v.ReadInConfig()
	if err != nil {
		return err
	}
	err = c.v.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *AgentConfig) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(c.v.ConfigFileUsed(), data, os.ModePerm)
}

// Config 站点配置
type Config struct {
	Debug    bool   // debug模式开关
	Language string // 系统语言，默认 zh-CN
	Site     struct {
		Brand          string // 站点名称
		CookieName     string // 浏览器 Cookie 名称
		Theme          string
		DashboardTheme string
		CustomCode     string
		ViewPassword   string // 前台查看密码
	}
	Oauth2 struct {
		Type         string
		Admin        string // 管理员用户名列表
		ClientID     string
		ClientSecret string
		Endpoint     string
	}
	HTTPPort      uint
	GRPCPort      uint
	GRPCHost      string
	ProxyGRPCPort uint
	TLS           bool

	EnablePlainIPInNotification bool // 通知信息IP不打码

	// IP变更提醒
	EnableIPChangeNotification bool
	IPChangeNotificationTag    string
	Cover                      uint8  // 覆盖范围（0:提醒未被 IgnoredIPNotification 包含的所有服务器; 1:仅提醒被 IgnoredIPNotification 包含的服务器;）
	IgnoredIPNotification      string // 特定服务器IP（多个服务器用逗号分隔）

	Location string // 时区，默认为 Asia/Shanghai

	v                              *viper.Viper
	IgnoredIPNotificationServerIDs map[uint64]bool // [ServerID] -> bool(值为true代表当前ServerID在特定服务器列表内）
	MaxTCPPingValue                int32
	AvgPingCount                   int

	// 动态域名解析更新
	DDNS struct {
		Enable             bool
		Provider           string
		AccessID           string
		AccessSecret       string
		WebhookURL         string
		WebhookMethod      string
		WebhookRequestBody string
		WebhookHeaders     string
		MaxRetries         uint32
	}
}

// Read 读取配置文件并应用
func (c *Config) Read(path string) error {
	c.v = viper.New()
	c.v.SetConfigFile(path)
	err := c.v.ReadInConfig()
	if err != nil {
		return err
	}

	err = c.v.Unmarshal(c)
	if err != nil {
		return err
	}

	if c.Site.Theme == "" {
		c.Site.Theme = "default"
	}
	if c.Site.DashboardTheme == "" {
		c.Site.DashboardTheme = "default"
	}
	if c.Language == "" {
		c.Language = "zh-CN"
	}
	if c.GRPCPort == 0 {
		c.GRPCPort = 5555
	}
	if c.EnableIPChangeNotification && c.IPChangeNotificationTag == "" {
		c.IPChangeNotificationTag = "default"
	}
	if c.Location == "" {
		c.Location = "Asia/Shanghai"
	}
	if c.MaxTCPPingValue == 0 {
		c.MaxTCPPingValue = 1000
	}
	if c.AvgPingCount == 0 {
		c.AvgPingCount = 2
	}
	if c.DDNS.Provider == "" {
		c.DDNS.Provider = "webhook"
	}
	if c.DDNS.WebhookMethod == "" {
		c.DDNS.WebhookMethod = "POST"
	}
	if c.DDNS.MaxRetries == 0 {
		c.DDNS.MaxRetries = 3
	}

	c.updateIgnoredIPNotificationID()
	return nil
}

// updateIgnoredIPNotificationID 更新用于判断服务器ID是否属于特定服务器的map
func (c *Config) updateIgnoredIPNotificationID() {
	c.IgnoredIPNotificationServerIDs = make(map[uint64]bool)
	splitedIDs := strings.Split(c.IgnoredIPNotification, ",")
	for i := 0; i < len(splitedIDs); i++ {
		id, _ := strconv.ParseUint(splitedIDs[i], 10, 64)
		if id > 0 {
			c.IgnoredIPNotificationServerIDs[id] = true
		}
	}
}

// Save 保存配置文件
func (c *Config) Save() error {
	c.updateIgnoredIPNotificationID()
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(c.v.ConfigFileUsed(), data, os.ModePerm)
}
