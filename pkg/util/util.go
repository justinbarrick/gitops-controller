package util

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	Scheme    = runtime.NewScheme()
	defaulter = runtime.ObjectDefaulter(Scheme)
	Log       = logf.Log.WithName("git-controller")
)

func init() {
	logf.SetLogger(logf.ZapLogger(false))
	corev1.AddToScheme(Scheme)
}

// Return the metadata of an object.
func GetMeta(obj runtime.Object) metav1.Object {
	meta, _ := meta.Accessor(obj)
	return meta
}

// Get the Group Version Kind of an object.
func GetType(obj runtime.Object) schema.GroupVersionKind {
	return obj.GetObjectKind().GroupVersionKind()
}

func Kind(kind, group, version string) runtime.Object {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Kind: kind,
		Group: group,
		Version: version,
	})
	return obj
}
