package ddns

import (
	"context"
	"os"
	"strconv"

	"github.com/XIU2/CloudflareSpeedTest/utils"
	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"gopkg.in/yaml.v3"
)

type CloudflareConfig struct {
	ApiToken string `yaml:"api-token"`
	ZoneId   string `yaml:"zone-id"`
	Domain   string `yaml:"domain"`
	Comment  string `yaml:"comment"`
}

type SpeedConfig struct {
	TestUrl       string `yaml:"test-url"`
	DownloadSpeed int    `yaml:"download-speed"`
}

type Config struct {
	Speed      SpeedConfig      `yaml:"speed"`
	Cloudflare CloudflareConfig `yaml:"cloudflare"`
}

var (
	// Cloudflare 客户端
	client *cloudflare.Client
	// Cloudflare 配置
	cloudflareConfig CloudflareConfig
	// 测速配置
	speedConfig SpeedConfig
)

func init() {
	// 初始化 Cloudflar 客户端
	client = cloudflare.NewClient(
		option.WithAPIToken("EhrMft9N5uuTzuU_1o-uB_qAuBor-IDGjfy_1bx6"),
	)
	// 打开配置文件
	file, err := os.Open("config.yaml")
	if err != nil {
		utils.Red.Printf("打开 config.yaml 配置文件时发生异常: %s\n", err)
		return
	}
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	config := &Config{}
	err = decoder.Decode(&config)
	if err != nil {
		utils.Red.Printf("解析 config.yaml 配置文件时发生异常: %s\n", err)
		return
	}
	cloudflareConfig = config.Cloudflare
	speedConfig = config.Speed
}

// GetSpeedConfig 获取测速配置
func GetSpeedConfig() SpeedConfig {
	return speedConfig
}

// UpdateDNSRecord 更新 DNS 记录
func UpdateDNSRecord(speedSet utils.DownloadSpeedSet) {
	speed := speedSet[0]
	downloadSpeed := strconv.FormatFloat(speed.DownloadSpeed/1024/1024, 'f', 2, 32)
	delay := strconv.FormatFloat(speed.Delay.Seconds()*1000, 'f', 2, 32)
	if speed.DownloadSpeed < 15 {
		utils.Cyan.Printf("跳过更新 DNS 记录, 最大下载速度为 %s; 平均延迟 %sms\n", downloadSpeed, delay)
		return
	}
	utils.Cyan.Printf("开始更新 DNS 解析记录 IP: %s; 域名: %s; 下载速度: %sMB/s; 平均延迟: %sms; 备注: %s\n", speed.IP, cloudflareConfig.Domain, downloadSpeed, delay, cloudflareConfig.Comment)
	// 查询记录
	record, err := findDomainRecord()
	if err != nil {
		return
	}
	// 更新记录
	domainRecord, err := updateDomainRecord(record.ID, speed.IP.String())
	if err != nil {
		utils.Red.Printf("更新域名解析记录失败 %s\n", err)
		return
	}
	utils.Cyan.Printf("DNS 解析记录更新成功  IP: %s; 域名: %s;\n", domainRecord.Content, domainRecord.Name)
}

// 校验是否存在 域名解析记录
func findDomainRecord() (dns.RecordResponse, error) {
	page, err := client.DNS.Records.List(context.TODO(), dns.RecordListParams{
		ZoneID: cloudflare.F(cloudflareConfig.ZoneId),
		Name: cloudflare.F(dns.RecordListParamsName{
			Exact: cloudflare.F(cloudflareConfig.Domain),
		}),
	})
	if err != nil {
		utils.Red.Printf("获取域名解析记录时发生异常: %s\n", err)
		return dns.RecordResponse{}, err
	}
	result := page.Result
	if len(page.Result) != 1 {
		utils.Yellow.Printf("没有找到对应的域名[%s]解析记录\n", cloudflareConfig.Domain)
	}
	return result[0], nil
}

// 更新域名解析记录
func updateDomainRecord(recordId string, content string) (*dns.RecordResponse, error) {
	res, err := client.DNS.Records.Edit(
		context.TODO(),
		recordId,
		dns.RecordEditParams{
			ZoneID: cloudflare.F(cloudflareConfig.ZoneId),
			Body: dns.ARecordParam{
				Name:    cloudflare.F(cloudflareConfig.Domain),
				TTL:     cloudflare.F(dns.TTL1),
				Type:    cloudflare.F(dns.ARecordTypeA),
				Content: cloudflare.F(content),
				Comment: cloudflare.F(cloudflareConfig.Comment),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
