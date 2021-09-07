package syncer

import (
	app "github.com/hchenc/application/pkg/client/clientset/versioned"
	"github.com/hchenc/devops-operator/pkg/models"
	harbor "github.com/hchenc/go-harbor"
	pager "github.com/hchenc/pager/pkg/client/clientset/versioned"
	git "github.com/xanzy/go-gitlab"
	"k8s.io/client-go/kubernetes"
)

const (
	DEVOPS_NAMESPACE = "devops-system"
)

type ClientSet struct {
	Client *kubernetes.Clientset

	AppClient *app.Clientset

	PagerClient *pager.Clientset

	GitlabClient *git.Client

	HarborClient *harbor.APIClient

	Config *models.Config
}

type GenerateService interface {
	//add obj to target service
	Add(obj interface{}) (interface{}, error)

	//update obj
	Update(objOld interface{}, objNew interface{}) error

	Delete(obj interface{}) error
}

func NewGenerateService(g Generator) GenerateService {
	return generator{g: g}
}

type Generator interface {
	Create(obj interface{}) (interface{}, error)
	//update obj
	Update(objOld interface{}, objNew interface{}) error

	Delete(obj interface{}) error

	GetByName(name string) (interface{}, error)

	GetByID(id int) (interface{}, error)

	List(key string) (interface{}, error)
}

type generator struct {
	g Generator
}

func (g generator) Add(obj interface{}) (interface{}, error) {
	return g.g.Create(obj)
}

func (g generator) Update(objOld interface{}, objNew interface{}) error {
	return g.g.Update(objOld, objNew)
}

func (g generator) Delete(obj interface{}) error {
	return g.g.Delete(obj)
}
