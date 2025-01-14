package ck8s

import (
	"context"
	"fmt"

	"github.com/daniarmas/notes/internal/clog"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Ck8s struct {
	Clientset *kubernetes.Clientset
}

var ck8s Ck8s

func NewClient() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		clog.Error(context.Background(), "error creating in-cluster config", err)
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		clog.Error(context.Background(), "error creating kubernetes clientset", err)
		return err
	}

	// set the client
	ck8s = Ck8s{
		Clientset: clientset,
	}

	return nil
}

// CreateJob creates a job in the k8s cluster
func CreateJob(ctx context.Context, jobName, namespace, imageName string, args []string) error {
	// Create a job spec
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  jobName,
							Image: imageName,
							Args:  args,
						},
					},
				},
			},
		},
	}

	// Create the job
	job, err := ck8s.Clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		clog.Error(context.Background(), "error creating job", err)
		return err
	}

	clog.Info(context.Background(), fmt.Sprintf("job created successful: %s", job.Name), nil)
	return nil
}
