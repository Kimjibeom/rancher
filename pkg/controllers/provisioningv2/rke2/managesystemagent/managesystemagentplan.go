package managesystemagent

import (
	"github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	rancherv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	namespaces "github.com/rancher/rancher/pkg/namespace"
	"github.com/rancher/rancher/pkg/settings"
	"github.com/rancher/wrangler/pkg/name"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (h *handler) OnChangeInstallSUC(cluster *rancherv1.Cluster, status rancherv1.ClusterStatus) ([]runtime.Object, rancherv1.ClusterStatus, error) {
	if cluster.Spec.RKEConfig == nil {
		return nil, status, nil
	}

	mcc := &v3.MultiClusterChart{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      name.SafeConcatName(cluster.Name, "managed", "system-upgrade-controller"),
		},
		Spec: v3.MultiClusterChartSpec{
			DefaultNamespace: namespaces.System,
			RepoName:         "rancher-charts",
			Chart:            "system-upgrade-controller",
			Version:          "0.7.001",
			Values: &v1alpha1.GenericMap{
				Data: map[string]interface{}{
					"global": map[string]interface{}{
						"systemDefaultRegistry": settings.SystemDefaultRegistry.Get(),
					},
				},
			},
			Targets: []v1alpha1.BundleTarget{
				{
					ClusterSelector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "metadata.name",
								Operator: metav1.LabelSelectorOpIn,
								Values: []string{
									cluster.Name,
								},
							},
							{
								Key:      "provisioning.cattle.io/unmanaged-system-agent",
								Operator: metav1.LabelSelectorOpDoesNotExist,
							},
						},
					},
				},
			},
		},
	}

	return []runtime.Object{
		mcc,
	}, status, nil
}
