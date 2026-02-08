package builder

import (
	"github.com/furutachiKurea/kb-adapter-rbdplugin/internal/model"
	"github.com/furutachiKurea/kb-adapter-rbdplugin/service/adapter"

	kbappsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
)

type MongoDB struct {
	Builder
}

var _ adapter.ClusterBuilder = &MongoDB{}

// BuildCluster 构建 Cluster struct
func (b *MongoDB) BuildCluster(input model.ClusterInput) (*kbappsv1.Cluster, error) {
	cluster, err := b.Builder.BuildCluster(input)
	if err != nil {
		return nil, err
	}

	cluster.Spec.Topology = "replicaset"

	if cluster.Spec.Backup != nil {
		cluster.Spec.Backup.Method = "datafile"
	}

	return cluster, nil
}
