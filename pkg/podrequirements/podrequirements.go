package podrequirements

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// PodRequirements binds together information about a pod
// and its aggregated resource requirements
type PodRequirements struct {
	// The name of the pod. Here we could use the
	// full pod specification, however we deliberately
	// choose to just use the pod name. This will
	// make JSON output a bit more readable
	Pod string
	// An standard resource requirements struct
	// that represents aggregated metrics for that pod,
	// with contributions from all containers in that pod
	Requirements corev1.ResourceRequirements
}

// List represents a list of PodRequirements
type List struct {
	// List of PodRequirements
	Items []PodRequirements
}

// ToJSON returns a JSON string representation
// of the list of pod requirements
func (list *List) ToJSON() (string, error) {
	e, err := json.Marshal(list)
	if err != nil {
		return "", err
	}
	return string(e), nil
}

// GetPodsCPUAndMemoryRequirements converts the given list of pods into
// a similar list where each pod is attached with its aggregated
// resource cpu and memory requirements
func GetPodsCPUAndMemoryRequirements(pods *corev1.PodList) *List {
	return getPodsRequirements(pods, "cpu", "memory")
}

// getPodsRequirements converts the given list of pods into
// a similar list where each pod is attached with its aggregated
// resource requirements. The set of
func getPodsRequirements(pods *corev1.PodList, resourceNames ...corev1.ResourceName) *List {

	items := []PodRequirements{}

	for _, pod := range pods.Items {

		// Build a vanilla resource requirements for the pod
		reqs := &corev1.ResourceRequirements{
			Limits:   make(map[corev1.ResourceName]resource.Quantity),
			Requests: make(map[corev1.ResourceName]resource.Quantity),
		}

		// aggregate cpu and memory requirements from all
		// init containers
		for _, container := range pod.Spec.InitContainers {
			for _, resourceName := range resourceNames {
				aggregateRequirements(reqs, container.Resources, resourceName)
			}
		}

		// aggregate cpu and memory requirements from all
		// standard containers
		for _, container := range pod.Spec.Containers {
			for _, resourceName := range resourceNames {
				aggregateRequirements(reqs, container.Resources, resourceName)
			}
		}

		// aggregate cpu and memory requirements from all
		// ephemeral containers
		for _, container := range pod.Spec.EphemeralContainers {
			for _, resourceName := range resourceNames {
				aggregateRequirements(reqs, container.Resources, resourceName)
			}
		}

		// attach the pod to its aggregated resource requirements
		// and add it to our list
		items = append(items, PodRequirements{
			Pod:          pod.Name,
			Requirements: *reqs,
		})
	}

	return &List{
		Items: items,
	}

}

// aggregateRequirements conveniently adds the requests and limits from the
// given reqs2 resource requirements, for the given resource name (eg. cpu or memory),
// into the given resource requirements struct
func aggregateRequirements(reqs *corev1.ResourceRequirements, reqs2 corev1.ResourceRequirements, resourceName corev1.ResourceName) {
	// aggregate limits
	aggregatedLimits := reqs.Limits[resourceName]
	aggregatedLimits.Add(reqs2.Limits[resourceName])
	reqs.Limits[resourceName] = aggregatedLimits

	// aggregate requests
	aggregatedRequests := reqs.Requests[resourceName]
	aggregatedRequests.Add(reqs2.Requests[resourceName])
	reqs.Requests[resourceName] = aggregatedRequests
}
