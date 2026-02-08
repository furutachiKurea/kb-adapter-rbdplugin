package coordinator

import (
	"fmt"

	"github.com/furutachiKurea/kb-adapter-rbdplugin/internal/model"
	"github.com/furutachiKurea/kb-adapter-rbdplugin/service/adapter"
	"gopkg.in/yaml.v3"
)

type MongoDB struct {
	Coordinator
}

var _ adapter.Coordinator = (*MongoDB)(nil)

// TargetPort 返回 KubeBlocks Cluster 的连接端口，
// 用于配置 KubeBlocksComponent 将连接转发至 Cluster 的 service
func (c *MongoDB) TargetPort() int {
	return 27017
}

// GetSecretName 返回该数据库类型的 Secret 命名格式
func (c *MongoDB) GetSecretName(clusterName string) string {
	return fmt.Sprintf("%s-mongodb-account-root", clusterName)
}

// GetBackupMethod 返回该数据库类型支持的备份方法
func (c *MongoDB) GetBackupMethod() string {
	return "datafile"
}

// GetParametersConfigMap 返回该类型的 Cluster 用于储存参数配置的 ConfigMap 名称，
// 并非所有的数据库类型都支持参数配置，不支持则返回 nil
func (c *MongoDB) GetParametersConfigMap(clusterName string) *string {
	cmName := fmt.Sprintf("%s-mongodb-mongodb-config", clusterName)
	return &cmName
}

// ParseParameters 解析 ConfigMap 中的配置文件参数
// configData 为 ConfigMap 的 data 字段，包含各种配置文件内容
func (c *MongoDB) ParseParameters(configData map[string]string) ([]model.ParameterEntry, error) {
	const configKey = "mongodb.conf"
	content, ok := configData[configKey]
	if !ok {
		return nil, fmt.Errorf("config key %s not found", configKey)
	}

	var config map[string]any
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	var params []model.ParameterEntry
	flattenConfig("", config, &params)
	return params, nil
}

func flattenConfig(prefix string, config map[string]any, params *[]model.ParameterEntry) {
	for k, v := range config {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		if nested, ok := v.(map[string]any); ok {
			flattenConfig(key, nested, params)
		} else {
			*params = append(*params, model.ParameterEntry{
				Name:  key,
				Value: v,
			})
		}
	}
}
