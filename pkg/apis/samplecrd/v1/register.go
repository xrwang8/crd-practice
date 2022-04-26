package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = "samplecrd.k8s.io"

//SchemeGroupVersion 注册自己的自定义资源
var SchemeGroupVersion = schema.GroupVersion{
	Group: GroupName, Version: "v1",
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Kind takes an unqualified resource and returns a Group qualified GroupKind
func Kind(resource string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(resource).GroupKind()
}

func addKnowTypes(scheme *runtime.Scheme) error {

	scheme.AddKnownTypes(SchemeGroupVersion,
		&Network{},
		&NetworkList{},
	)
	// register the type in the scheme
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
