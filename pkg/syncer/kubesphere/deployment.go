package kubesphere

import (
	"context"
	"github.com/hchenc/devops-operator/pkg/syncer"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type deploymentInfo struct {
	kubeClient *kubernetes.Clientset
	ctx        context.Context
}

func (d deploymentInfo) Create(obj interface{}) (interface{}, error) {
	deployment := obj.(*v1.Deployment)

	namespacePrefix := strings.Split(deployment.Namespace, "-")[0]
	candidates := map[string]string{
		namespacePrefix + "-fat":     "fat",
		namespacePrefix + "-uat":     "uat",
		namespacePrefix + "-smoking": "smoking",
	}
	delete(candidates, deployment.Namespace)

	for namespace := range candidates {
		deployment := assembleDeployment(deployment, namespace)
		_, err := d.kubeClient.AppsV1().Deployments(namespace).Create(d.ctx, deployment, metav1.CreateOptions{})
		if err == nil || errors.IsAlreadyExists(err) {
			continue
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (d deploymentInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (d deploymentInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (d deploymentInfo) GetByName(name string) (interface{}, error) {
	panic("implement me")
}

func (d deploymentInfo) GetByID(id int) (interface{}, error) {
	panic("implement me")
}

func (d deploymentInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

func assembleDeployment(deployment *v1.Deployment, namespace string) *v1.Deployment {
	return &v1.Deployment{
		TypeMeta: deployment.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:        deployment.Name,
			Namespace:   namespace,
			Labels:      deployment.Labels,
			Annotations: deployment.Annotations,
			Finalizers:  deployment.Finalizers,
			ClusterName: deployment.ClusterName,
		},
		Spec: v1.DeploymentSpec{
			Replicas: deployment.Spec.Replicas,
			Selector: deployment.Spec.Selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      deployment.Spec.Template.Labels,
					Annotations: deployment.Spec.Template.Annotations,
				},
				Spec: corev1.PodSpec{
					Containers:         deployment.Spec.Template.Spec.Containers,
					ServiceAccountName: deployment.Spec.Template.Spec.ServiceAccountName,
					Affinity:           deployment.Spec.Template.Spec.Affinity,
					InitContainers:     deployment.Spec.Template.Spec.InitContainers,
					Volumes:            deployment.Spec.Template.Spec.Volumes,
					ImagePullSecrets:   deployment.Spec.Template.Spec.ImagePullSecrets,
				},
			},
			Strategy: deployment.Spec.Strategy,
		},
	}
}

func NewDeploymentGenerator(ctx context.Context, kubeClient *kubernetes.Clientset) syncer.Generator {
	return deploymentInfo{
		kubeClient: kubeClient,
		ctx:        ctx,
	}
}
