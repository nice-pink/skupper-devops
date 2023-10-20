package kynetes

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// variables

var (
	KubeConfigPath string = "/home/vscode/go/src/github.com/nice-pink/skupper-devops/.kube/config"
	IsInCluster    bool   = false
)

func getClientSet() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if IsInCluster {
		// get config from cluster
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
		// get config from file
		config, err = clientcmd.BuildConfigFromFlags("", KubeConfigPath)
		if err != nil {
			return nil, err
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func getClientSetDynamic() (*dynamic.DynamicClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", KubeConfigPath)
	if err != nil {
		return nil, err
	}
	// create the clientset
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// deployment

func GetDeployment(name string, namespace string) error {
	// client set
	clientset, err := getClientSet()
	if err != nil {
		panic(err)
	}

	// get deployment
	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Deployment %s in namespace %s not found\n", name, namespace)
		return err
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Deployment %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message)
		return err
	} else if err != nil {
		fmt.Printf("Error getting Deployment %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message)
		return err
	}

	fmt.Printf("Found Deployment %s in namespace %s\n", name, namespace)
	return nil
}

func DeleteDeployment(name string, namespace string) error {
	return DeleteResource(name, namespace, "deployments", "v1", "apps")
}

func RestartDeployment(name string, namespace string) error {
	// client set
	clientset, err := getClientSet()
	if err != nil {
		panic(err)
	}

	// get deployment
	deploy, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// init annotation if not present
	if deploy.Spec.Template.ObjectMeta.Annotations == nil {
		deploy.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	// update restartedAt annotation
	deploy.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	// update deployment
	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	return err
}

// pod

func ListPods(namespace string) {
	// client set
	clientset, err := getClientSet()
	if err != nil {
		panic(err)
	}

	// get pod
	// if namespace is "", then all pods are listed
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the namespace %s\n", len(pods.Items), namespace)
	for _, pod := range pods.Items {
		fmt.Println("- " + pod.Name)
	}
}

func GetPod(name string, namespace string) {
	// client set
	clientset, err := getClientSet()
	if err != nil {
		panic(err)
	}

	// get pod
	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", name, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", name, namespace)
	}
}

// secret

func ListSecrets(namespace string) {
	// client set
	clientset, err := getClientSet()
	if err != nil {
		panic(err)
	}

	// get pod
	// if namespace is "", then all pods are listed
	secrets, err := clientset.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d secrets in the namespace %s\n", len(secrets.Items), namespace)
	for _, secret := range secrets.Items {
		fmt.Println("- " + secret.Name)
	}
}

func GetSecret(name string, namespace string, print bool) *v1.Secret {
	// client set
	clientset, err := getClientSet()
	if err != nil {
		panic(err)
	}

	// get secret
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Secret %s in namespace %s not found\n", name, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting secret %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found secret %s in namespace %s\n", name, namespace)
	}

	if print {
		fmt.Println(secret)
	}

	return secret
}

func CreateSecret(namespace string, content []byte) error {
	object, err := UnmarshalYAML(content)
	if err != nil {
		return err
	}

	// create resource from interface
	secretResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}
	secret := &unstructured.Unstructured{
		Object: object,
	}

	// client
	client, err := getClientSetDynamic()
	if err != nil {
		return err
	}

	// create resource
	resource, err := client.Resource(secretResource).Namespace(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created secret %q.\n", resource.GetName())

	return nil
}

func SecretHasData(name string, namespace string) bool {
	secret := GetSecret(name, namespace, false)
	return len(secret.Data) > 0
}

func DeleteSecret(name string, namespace string) error {
	return DeleteResource(name, namespace, "secrets", "v1", "")
}

// config map

func GetConfigMap(name string, namespace string, print bool) *v1.ConfigMap {
	// client set
	clientset, err := getClientSet()
	if err != nil {
		panic(err)
	}

	// get configMap
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("ConfigMap %s in namespace %s not found\n", name, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting ConfigMap %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found ConfigMap %s in namespace %s\n", name, namespace)
	}

	if print {
		fmt.Println(configMap)
	}

	return configMap
}

func ConfigMapHasData(name string, namespace string) bool {
	configMap := GetConfigMap(name, namespace, false)
	return len(configMap.Data) > 0
}

func ConfigMapGetData(name string, namespace string) map[string]string {
	configMap := GetConfigMap(name, namespace, false)
	return configMap.Data
}

func ConfigMapGetResourceVersion(name string, namespace string) string {
	configMap := GetConfigMap(name, namespace, false)
	return configMap.ResourceVersion
}

// resource

func ReadResource(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	return byteValue, err
}

func DumpResource(name string, namespace string, resource schema.GroupVersionResource, path string) error {
	// client
	client, err := getClientSetDynamic()
	if err != nil {
		return err
	}

	// get resource
	resourceDef, err := client.Resource(resource).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// dump to file
	file, err := os.Create(path)
	yamlPrinter := printers.YAMLPrinter{}
	defer file.Close()
	yamlPrinter.PrintObj(resourceDef, file)

	return nil
}

func DeleteResource(name string, namespace string, resourceType string, version string, group string) error {
	// client
	client, err := getClientSetDynamic()
	if err != nil {
		return err
	}

	// delete
	fmt.Println("Deleting deployment...")
	resource := schema.GroupVersionResource{Group: group, Version: version, Resource: resourceType}
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := client.Resource(resource).Namespace(namespace).Delete(context.TODO(), name, deleteOptions); err != nil {
		return err
	}
	fmt.Println("Deleted secret.")

	return nil
}

// https://github.com/elastic/beats/blob/6435194af9f42cbf778ca0a1a92276caf41a0da8/libbeat/common/mapstr.go

type MapStr map[string]interface{}

func UnmarshalYAML(data []byte) (MapStr, error) {
	var result map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	ms := cleanUpInterfaceMap(result)
	return ms, nil
}

func cleanUpInterfaceArray(in []interface{}) []interface{} {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpInterfaceMap(in map[interface{}]interface{}) MapStr {
	result := make(MapStr)
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanUpInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanUpInterfaceMap(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
