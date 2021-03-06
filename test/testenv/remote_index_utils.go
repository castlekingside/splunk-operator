package testenv

import (
	"encoding/json"

	enterprisev1 "github.com/splunk/splunk-operator/pkg/apis/enterprise/v1beta1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// dataIndexesResponse struct for /data/indexes response
type dataIndexesResponse struct {
	Entry []struct {
		Name string `json:"name"`
	} `json:"entry"`
}

// GetIndexOnPod get list of indexes on given pod
func GetIndexOnPod(deployment *Deployment, podName string, indexName string) bool {
	stdin := "curl -ks -u admin:$(cat /mnt/splunk-secrets/password) https://localhost:8089/services/data/indexes?output_mode=json"
	command := []string{"/bin/sh"}
	stdout, stderr, err := deployment.PodExecCommand(podName, command, stdin, false)
	if err != nil {
		logf.Log.Error(err, "Failed to execute command on pod", "pod", podName, "command", command)
		return false
	}
	logf.Log.Info("Command executed on pod", "pod", podName, "command", command, "stdin", stdin, "stdout", stdout, "stderr", stderr)
	restResponse := dataIndexesResponse{}
	err = json.Unmarshal([]byte(stdout), &restResponse)
	if err != nil {
		logf.Log.Error(err, "Failed to parse data/indexes response")
		return false
	}
	indexFound := false
	for _, entry := range restResponse.Entry {
		if entry.Name == indexName {
			indexFound = true
			break
		}
	}
	return indexFound
}

// RestartSplunk Restart splunk inside the container
func RestartSplunk(deployment *Deployment, podName string) bool {
	stdin := "/opt/splunk/bin/splunk restart -auth admin:$(cat /mnt/splunk-secrets/password)"
	command := []string{"/bin/sh"}
	stdout, stderr, err := deployment.PodExecCommand(podName, command, stdin, false)
	if err != nil {
		logf.Log.Error(err, "Failed to execute command on pod", "pod", podName, "command", command)
		return false
	}
	logf.Log.Info("Command executed on pod", "pod", podName, "command", command, "stdin", stdin, "stdout", stdout, "stderr", stderr)
	return true
}

// RollHotToWarm rolls hot buckets to warm for a given index and pod
func RollHotToWarm(deployment *Deployment, podName string, indexName string) bool {
	stdin := "/opt/splunk/bin/splunk _internal call /data/indexes/" + indexName + "/roll-hot-buckets admin:$(cat /mnt/splunk-secrets/password)"
	command := []string{"/bin/sh"}
	stdout, stderr, err := deployment.PodExecCommand(podName, command, stdin, false)
	if err != nil {
		logf.Log.Error(err, "Failed to execute command on pod", "pod", podName, "command", command)
		return false
	}
	logf.Log.Info("Command executed on pod", "pod", podName, "command", command, "stdin", stdin, "stdout", stdout, "stderr", stderr)
	return true
}

// GenerateIndexVolumeSpec return VolumeSpec struct with given values
func GenerateIndexVolumeSpec(volumeName string, endpoint string, Path string, secretRef string) enterprisev1.VolumeSpec {
	return enterprisev1.VolumeSpec{
		Name:      volumeName,
		Endpoint:  endpoint,
		Path:      testIndexesS3Bucket,
		SecretRef: secretRef,
	}
}

// GenerateIndexSpec return VolumeSpec struct with given values
func GenerateIndexSpec(indexName string, volName string) enterprisev1.IndexSpec {
	return enterprisev1.IndexSpec{
		Name:       indexName,
		RemotePath: indexName,
		IndexAndGlobalCommonSpec: enterprisev1.IndexAndGlobalCommonSpec{
			VolName: volName,
		},
	}
}
