package kubesphere

import (
	"context"
	"github.com/hchenc/devops-operator/pkg/syncer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type volumeInfo struct {
	kubeClient *kubernetes.Clientset
	ctx        context.Context
}

func (v volumeInfo) Create(obj interface{}) (interface{}, error) {
	volume := obj.(*v1.PersistentVolume)

	namespacePrefix := strings.Split(volume.Namespace, "-")[0]
	candidates := map[string]string{
		namespacePrefix + "-fat":     "fat",
		namespacePrefix + "-uat":     "uat",
		namespacePrefix + "-smoking": "smoking",
	}
	delete(candidates, volume.Namespace)

	for namespace := range candidates {
		volume := assembleVolume(volume, namespace)
		_, err := v.kubeClient.CoreV1().PersistentVolumes().Create(v.ctx, volume, metav1.CreateOptions{})
		if err == nil || errors.IsAlreadyExists(err) {
			continue
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (v volumeInfo) Update(objOld interface{}, objNew interface{}) error {
	panic("implement me")
}

func (v volumeInfo) Delete(obj interface{}) error {
	panic("implement me")
}

func (v volumeInfo) GetByName(name string) (interface{}, error) {
	panic("implement me")
}

func (v volumeInfo) GetByID(id int) (interface{}, error) {
	panic("implement me")
}

func (v volumeInfo) List(key string) (interface{}, error) {
	panic("implement me")
}

// TODO
func assembleVolume(volume *v1.PersistentVolume, namespace string) *v1.PersistentVolume {
	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        volume.Name,
			Namespace:   namespace,
			Labels:      volume.Labels,
			Annotations: volume.Annotations,
		},
		Spec: v1.PersistentVolumeSpec{},
	}
}

func NewVolumeGenerator(ctx context.Context, clientset *kubernetes.Clientset) syncer.Generator {
	return volumeInfo{
		kubeClient: clientset,
		ctx:        ctx,
	}
}
