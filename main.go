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
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

// SetupSignalHandler 注册 SIGTERM 和 SIGINT 信号
// 返回一个 stop channel，该通道在捕获到第一个信号时被关闭
// 如果捕捉到第二个信号，程序将直接退出
func setupSignalHandler() (stopCh <-chan struct{}) {
	// 当调用两次的时候 panics
	close(onlyOneSignalHandler)
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	// Notify 函数让 signal 包将输入信号转发到 c
	// 如果没有列出要传递的信号，会将所有输入信号传递到 c；否则只传递列出的输入信号
	// 参考文档：https://cloud.tencent.com/developer/article/1645996
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // 第二个信号，直接退出
	}()

	return stop
}

func main() {
	//设置一个信号处理，应用于优雅关闭
	stopCh := setupSignalHandler()
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

	// 实例化自定义的控制器

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
