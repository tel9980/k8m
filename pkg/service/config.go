package service

import (
	"github.com/fatih/color"
	"github.com/weibaohui/k8m/internal/dao"
	"github.com/weibaohui/k8m/pkg/comm/utils"
	"github.com/weibaohui/k8m/pkg/flag"
	"github.com/weibaohui/k8m/pkg/models"
	"gorm.io/gorm"
	"k8s.io/klog"
)

type configService struct {
	db *gorm.DB
}

func NewConfigService() *configService {
	return &configService{db: dao.DB()}
}

func (s *configService) GetConfig() (*models.Config, error) {
	var config models.Config
	if err := s.db.First(&config).Error; err != nil {
		return nil, err
	}

	return &config, nil
}

func (s *configService) UpdateConfig(config *models.Config) error {

	err := s.db.Save(config).Error
	if err != nil {
		return err
	}
	// 保存后，让其生效
	err = s.UpdateFlagFromDBConfig()
	if err != nil {
		return err
	}
	return nil
}

// UpdateFlagFromDBConfig 从数据库中加载配置，更新Flag方法中的值
func (s *configService) UpdateFlagFromDBConfig() error {
	cfg := flag.Init()
	m, err := s.GetConfig()
	if err != nil {
		return err
	}
	if cfg.PrintConfig {
		klog.Infof("已开启配置信息打印选项。下面是数据库配置的回显.\n%s:\n %+v\n%s\n", color.RedString("↓↓↓↓↓↓生产环境请务必关闭↓↓↓↓↓↓"), utils.ToJSON(m), color.RedString("↑↑↑↑↑↑生产环境请务必关闭↑↑↑↑↑↑"))
		cfg.ShowConfigCloseMethod()
	}

	cfg.AnySelect = m.AnySelect
	// cfg.Debug = m.Debug
	// cfg.InCluster = m.InCluster

	cfg.ApiKey = m.ApiKey
	cfg.ApiModel = m.ApiModel
	cfg.ApiURL = m.ApiURL
	cfg.UseBuiltInModel = m.UseBuiltInModel
	// if m.KubeConfig != "" {
	// 	cfg.KubeConfig = m.KubeConfig
	// }
	if m.KubectlShellImage != "" {
		cfg.KubectlShellImage = m.KubectlShellImage
	}
	if m.NodeShellImage != "" {
		cfg.NodeShellImage = m.NodeShellImage
	}
	// if m.Port > 0 {
	// 	cfg.Port = m.Port
	// }
	// if m.SqlitePath != "" {
	// 	cfg.SqlitePath = m.SqlitePath
	// }
	// if m.MCPServerPort > 0 {
	// 	cfg.MCPServerPort = m.MCPServerPort
	// }
	// if m.LogV > 0 {
	// 	cfg.LogV = m.LogV
	// }

	cfg.PrintConfig = m.PrintConfig
	cfg.EnableAI = m.EnableAI
	// cfg.ConnectCluster = m.ConnectCluster

	// JwtTokenSecret 暂不启用，因为前端也要处理
	// cfg.JwtTokenSecret = m.JwtTokenSecret
	// LoginType 暂不启用，因为就一种password
	// cfg.LoginType = m.LoginType

	return nil
}
