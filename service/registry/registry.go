package registry

import (
	"github.com/furutachiKurea/kb-adapter-rbdplugin/internal/log"
	"github.com/furutachiKurea/kb-adapter-rbdplugin/service/adapter"
	"github.com/furutachiKurea/kb-adapter-rbdplugin/service/builder"
	"github.com/furutachiKurea/kb-adapter-rbdplugin/service/coordinator"
)

// Cluster 在这里注册 Block Mechanica 支持的数据库集群
var Cluster = map[string]adapter.ClusterAdapter{
	"postgresql": _postgresql,
	"mysql":      _mysql,
	"redis":      _redis,
	"rabbitmq":   _rabbitmq,
	"mongodb":    _mongodb,
	// ... new types here
}

var (
	_postgresql = adapter.ClusterAdapter{
		Builder:     &builder.PostgreSQL{},
		Coordinator: &coordinator.PostgreSQL{},
	}

	_mysql = adapter.ClusterAdapter{
		Builder:     &builder.MySQL{},
		Coordinator: &coordinator.MySQL{},
	}

	_redis = adapter.ClusterAdapter{
		Builder:     &builder.Redis{},
		Coordinator: &coordinator.Redis{},
	}
	_rabbitmq = adapter.ClusterAdapter{
		Builder:     &builder.RabbitMQ{},
		Coordinator: &coordinator.RabbitMQ{},
	}
	// TODO Rainbond 目前不支持为组件创建多个同端口的 service，
	// 因此无法做到为 MongoDB 分别创建用于读写（主）和只读（从）的 service，
	// Rainbond 只提供主节点访问的 service
	_mongodb = adapter.ClusterAdapter{
		Builder:     &builder.MongoDB{},
		Coordinator: &coordinator.MongoDB{},
	}
)

// init 函数进行注册表验证
func init() {
	validateClusterRegistry()
}

// validateClusterRegistry 验证集群注册表的完整性
func validateClusterRegistry() {
	for dbType, clusterAdapter := range Cluster {
		if err := clusterAdapter.Validate(); err != nil {
			log.Fatal("Critical validation error", log.String("DB Type", dbType), log.Err(err))
		}
		log.Info("Database validation passed", log.String("DB Type", dbType))
	}
	log.Info("All database validation passed")
}
