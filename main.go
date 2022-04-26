package main

import (
	clientset "crd-practice/pkg/client/clientset/versioned"
	informers "crd-practice/pkg/client/informers/externalversions"
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"path/filepath"
	"time"
)

func main() {

	flag.Parse()
	_, cfg, err := initClient()
	if err != nil {
		klog.Fatalf("Error building kubernetes clientSet: %s", err.Error())
	}

	networkClient, err := clientset.NewForConfig(cfg)

	if err != nil {
		klog.Fatalf("Error building networkClient clientSet: %s", err.Error())
	}
	// informerFactory 工厂类， 这里注入我们通过代码生成的 client
	// clent 主要用于和 API Server 进行通信，实现 ListAndWatch
	informers.NewSharedInformerFactory(networkClient, 5*time.Second)

}

func initClient() (*kubernetes.Clientset, *rest.Config, error) {
	var err error
	var config *rest.Config
	// inCluster（Pod）、KubeConfig（kubectl）
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(可选) kubeconfig 文件的绝对路径")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig 文件的绝对路径")
	}
	flag.Parse()

	// 首先使用 inCluster 模式(需要去配置对应的 RBAC 权限，默认的sa是default->是没有获取deployments的List权限)
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 KubeConfig 文件创建集群配置 Config 对象
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			return nil, nil, err
		}
	}

	// 已经获得了 rest.Config 对象
	// 创建 Clientset 对象
	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, config, err
	}
	return kubeclient, config, nil
}
