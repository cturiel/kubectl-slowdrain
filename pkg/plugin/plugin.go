package plugin

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cturiel/kubectl-slowdrain/pkg/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewKubeClient create a new Kubernetes client
func NewKubeClient(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creando cliente de Kubernetes: %w", err)
	}
	return clientset, nil
}

// DrainNode drains a Kubernetes node by deleting application pods one by one with a delay
func DrainNode(ctx context.Context, clientset *kubernetes.Clientset, nodeName string, delaySeconds int, autoConfirm bool, logLevel string, infraPrefixes []string) error {
	log := logger.NewLogger(logLevel)

	// 1.- Cordon the node
	log.Info("Mark the node %s as 'cordon'", nodeName)
	if err := cordonNode(ctx, clientset, nodeName); err != nil {
		return err
	}

	// 2️.- Get the node pods and classify them into infra vs apps
	log.Info("Getting pods in the node %s...", nodeName)
	infraPods, appPods, err := getNodePods(ctx, clientset, nodeName, infraPrefixes)
	if err != nil {
		return err
	}

	log.Info("Infrastructure pods detected: %d", len(infraPods))
	log.Info("Applications pods detected: %d", len(appPods))

	fmt.Printf("Node %s has the following pods scheduled:\n\n", nodeName)
	printPodsTable(infraPods, appPods)

	// 3️.- Show the list of application pods to be deleted and ask for confirmation
	if len(appPods) > 0 {
		fmt.Printf("The following applications pods will be deleted (NAMESPACE/POD):\n\n")
		for _, pod := range appPods {
			fmt.Printf(" - %s/%s\n", pod.Namespace, pod.Name)
		}
		fmt.Println()

		// Confirm the deletion of the pods
		if autoConfirm {
			log.Warn("Auto-confirmation enabled (-y). Proceeding with draining node %s.", nodeName)
		} else {
			fmt.Print("Confirm pod deletion? [y/N]: ")
			var confirm string
			fmt.Scanln(&confirm)

			confirm = strings.ToLower(strings.TrimSpace(confirm))

			if confirm != "yes" && confirm != "y" {
				log.Warn("Node %s drain cancelled", nodeName)
				return nil
			}
		}
	} else {
		log.Info("There are no application pods on the node %s", nodeName)
		return nil
	}

	// 4️.- Delete application pods one by one with a delay and in a random order
	appPods = shufflePods(appPods)

	for _, pod := range appPods {
		log.Info("Deleting pod %s/%s...", pod.Namespace, pod.Name)
		err := clientset.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})

		if err != nil {
			log.Error("Error deleting pod %s/%s: %v", pod.Namespace, pod.Name, err)
		} else {
			log.Info("Pod %s/%s successfully deleted", pod.Namespace, pod.Name)
		}

		log.Info("Waiting %d seconds before deleting the next pod...", delaySeconds)
		time.Sleep(time.Duration(delaySeconds) * time.Second)
	}

	log.Info("Node %s drain completed", nodeName)
	return nil
}

// cordonNode applies the 'cordon' action to a node, marking it as unschedulable
func cordonNode(ctx context.Context, clientset *kubernetes.Clientset, nodeName string) error {
	node, err := clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting node %s: %w", nodeName, err)
	}

	if node.Spec.Unschedulable {
		fmt.Printf("Warn: Node %s was already cordoned.\n", nodeName)
		return nil
	}

	node.Spec.Unschedulable = true
	_, err = clientset.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error applying cordon in the node %s: %w", nodeName, err)
	}

	fmt.Printf("Node %s marked as 'cordon'.\n", nodeName)
	return nil
}

// getNodePods Gets the pods running on a node and classifies them as infrastructure or application pods
func getNodePods(ctx context.Context, clientset *kubernetes.Clientset, nodeName string, infraPrefixes []string) ([]v1.Pod, []v1.Pod, error) {
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error listing pods of the node %s: %w", nodeName, err)
	}

	var infraPods []v1.Pod
	var appPods []v1.Pod

	for _, pod := range pods.Items {
		if isInfraPod(pod, infraPrefixes) {
			infraPods = append(infraPods, pod)
		} else {
			appPods = append(appPods, pod)
		}
	}

	return infraPods, appPods, nil
}

// isInfraPod Determines if a pod is an infrastructure pod based on its namespace
func isInfraPod(pod v1.Pod, infraPrefixes []string) bool {
	for _, prefix := range infraPrefixes {
		if len(pod.Namespace) >= len(prefix) && pod.Namespace[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func printPodsTable(infraPods, appPods []v1.Pod) {
	w := tabwriter.NewWriter(os.Stdout, 4, 8, 2, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "NAMESPACE\tPOD NAME")
	fmt.Fprintln(w, "---------\t--------")

	// Print infrastructure pods
	for _, pod := range infraPods {
		fmt.Fprintf(w, "%s\t%s\n", pod.Namespace, pod.Name)
	}

	// Print application pods
	for _, pod := range appPods {
		fmt.Fprintf(w, "%s\t%s\n", pod.Namespace, pod.Name)
	}

	fmt.Fprintln(w)

	w.Flush()
}

// shufflePods shuffles a slice of pods to avoid deleting pods in a predictable order
func shufflePods(pods []v1.Pod) []v1.Pod {
	shuffled := make([]v1.Pod, len(pods))
	copy(shuffled, pods) // Copy the original slice to avoid modifying it

	// Shuffle the slice of pods using the Fisher-Yates algorithm
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
	return shuffled
}
