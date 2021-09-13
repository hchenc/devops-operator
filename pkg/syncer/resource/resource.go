package resource

type assembleResourceFunc func(obj interface{}, namespace string) interface{}


func setupResource(obj interface{}, namespace string, resourceFunc assembleResourceFunc) interface{} {
	return resourceFunc(obj, namespace)
}