package utils

//var (
//	Config        *restclient.Config
//	Clientset     *kubernetes.Clientset
//	DynamicClient dynamic.Interface
//	GitLabClient  *git.Client
//	once          sync.Once
//	err           error
//)

//func init() {
//	var kubeconfig *string
//	if home := homeDir(); home != "" {
//		kubeconfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "Config"), "(optional) absolute path to the kubeconfig file")
//	} else {
//		kubeconfig = flag.String("kubeConfig", "", "absolute path to the kubeconfig file")
//	}
//	// use the current context in kubeconfig
//	Config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
//
//	// create the ClientSet
//	Clientset, err = kubernetes.NewForConfig(Config)
//
//	// create the DynamicClientSet
//	DynamicClient = dynamic.NewForConfigOrDie(Config)
//
//}

//func SetupGitLabClient(host, port, user, password, token string) error {
//	url := "http://" + host + ":" + port
//	if token != "" {
//		GitLabClient, err = git.NewClient(token, git.WithBaseURL(url))
//		return err
//	} else if user != "" && password != "" {
//		GitLabClient, err = git.NewBasicAuthClient(user, password, git.WithBaseURL(url))
//		return err
//	} else {
//		return errors.New("gitlab certification not provided")
//	}
//}

//func SetupClientSet(kubeconfig *string) error {
//	Config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
//	if err != nil {
//		return err
//	}
//	// create the Clientset
//	Clientset, err = kubernetes.NewForConfig(Config)
//	if err != nil {
//		return err
//	}
//	DynamicClient = dynamic.NewForConfigOrDie(Config)
//	return nil
//}

//func GetClientSet() (*restclient.Config, *kubernetes.Clientset, error) {
//	once.Do(func() {
//		var kubeconfig *string
//		if home := homeDir(); home != "" {
//			kubeconfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "Config"), "(optional) absolute path to the kubeconfig file")
//		} else {
//			kubeconfig = flag.String("kubeConfig", "", "absolute path to the kubeconfig file")
//		}
//		// use the current context in kubeconfig
//		Config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
//
//		// create the Clientset
//		Clientset, err = kubernetes.NewForConfig(Config)
//	})
//	return Config, Clientset, err
//}

//func homeDir() string {
//	if h := os.Getenv("HOME"); h != "" {
//		return h
//	}
//	return os.Getenv("USERPROFILE") // windows
//}
