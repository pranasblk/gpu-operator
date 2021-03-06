package clusterpolicy

import (
	"context"
	"fmt"
	"os"

	gpuv1 "github.com/NVIDIA/gpu-operator/pkg/apis/nvidia/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type controlFunc []func(n ClusterPolicyController) (gpuv1.State, error)

func ServiceAccount(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ServiceAccount.DeepCopy()
	logger := log.WithValues("ServiceAccount", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func Role(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].Role.DeepCopy()
	logger := log.WithValues("Role", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func RoleBinding(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].RoleBinding.DeepCopy()
	logger := log.WithValues("RoleBinding", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func ClusterRole(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ClusterRole.DeepCopy()
	logger := log.WithValues("ClusterRole", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func ClusterRoleBinding(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ClusterRoleBinding.DeepCopy()
	logger := log.WithValues("ClusterRoleBinding", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func ConfigMap(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ConfigMap.DeepCopy()
	logger := log.WithValues("ConfigMap", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func kernelFullVersion(n ClusterPolicyController) (string, string) {
	logger := log.WithValues("Request.Namespace", "default", "Request.Name", "Node")
	// We need the node labels to fetch the correct container
	opts := []client.ListOption{
		client.MatchingLabels{"feature.node.kubernetes.io/pci-10de.present": "true"},
	}

	list := &corev1.NodeList{}
	err := n.rec.client.List(context.TODO(), list, opts...)
	if err != nil {
		logger.Info("Could not get NodeList", "ERROR", err)
		return "", ""
	}

	if len(list.Items) == 0 {
		// none of the nodes matched a pci-10de label
		// either the nodes do not have GPUs, or NFD is not running
		logger.Info("Could not get any nodes to match pci-0302_10de.present=true label", "ERROR", "")
		return "", ""
	}

	// Assuming all nodes are running the same kernel version,
	// One could easily add driver-kernel-versions for each node.
	node := list.Items[0]
	labels := node.GetLabels()

	var ok bool
	kernelFullVersion, ok := labels["feature.node.kubernetes.io/kernel-version.full"]
	if ok {
		logger.Info(kernelFullVersion)
	} else {
		err := errors.NewNotFound(schema.GroupResource{Group: "Node", Resource: "Label"},
			"feature.node.kubernetes.io/kernel-version.full")
		logger.Info("Couldn't get kernelVersion, did you run the node feature discovery?", err)
		return "", ""
	}

	osName, ok := labels["feature.node.kubernetes.io/system-os_release.ID"]
	if !ok {
		return kernelFullVersion, ""
	}
	osVersion, ok := labels["feature.node.kubernetes.io/system-os_release.VERSION_ID"]
	if !ok {
		return kernelFullVersion, ""
	}
	osTag := fmt.Sprintf("%s%s", osName, osVersion)

	return kernelFullVersion, osTag
}

func getDcgmExporter() string {
	dcgmExporter := os.Getenv("NVIDIA_DCGM_EXPORTER")
	if dcgmExporter == "" {
		log.Info(fmt.Sprintf("ERROR: Could not find environment variable NVIDIA_DCGM_EXPORTER"))
		os.Exit(1)
	}
	return dcgmExporter
}

func preProcessDaemonSet(obj *appsv1.DaemonSet, n ClusterPolicyController) {
	transformations := map[string]func(*appsv1.DaemonSet, *gpuv1.ClusterPolicySpec, ClusterPolicyController) error{
		"nvidia-driver-daemonset":            TransformDriver,
		"nvidia-container-toolkit-daemonset": TransformToolkit,
		"nvidia-device-plugin-daemonset":     TransformDevicePlugin,
		"nvidia-dcgm-exporter":               TransformDCGMExporter,
	}

	t, ok := transformations[obj.Name]
	if !ok {
		log.Info(fmt.Sprintf("No transformation for Daemonset '%s'", obj.Name))
		return
	}

	err := t(obj, &n.singleton.Spec, n)
	if err != nil {
		log.Info(fmt.Sprintf("Failed to apply transformation '%s' with error: '%v'", obj.Name, err))
		os.Exit(1)
	}
}

func TransformDriver(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	kvers, osTag := kernelFullVersion(n)
	if osTag == "" {
		return fmt.Errorf("ERROR: Could not find kernel full version: ('%s', '%s')", kvers, osTag)
	}

	img := fmt.Sprintf("%s-%s", config.Driver.ImagePath(), osTag)
	obj.Spec.Template.Spec.Containers[0].Image = img

	if osTag != "rhel" && osTag != "centos" {
		return nil
	}

	entitlementPath := "/etc/pki/entitlements"
	if _, err := os.Stat(entitlementPath); os.IsNotExist(err) {
		log.Info(fmt.Sprintf("ERROR: Could not find RedHat entitlement on current node at path %s", entitlementPath))
		os.Exit(1)
	}

	volName, volSecretName := "openshift-entitlements", "entitlement"
	volMount := corev1.VolumeMount{Name: volName, ReadOnly: true, MountPath: entitlementPath}
	obj.Spec.Template.Spec.Containers[0].VolumeMounts = append(obj.Spec.Template.Spec.Containers[0].VolumeMounts, volMount)

	vol := corev1.Volume{Name: volName, VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: volSecretName}}}
	obj.Spec.Template.Spec.Volumes = append(obj.Spec.Template.Spec.Volumes, vol)

	return nil
}

func TransformToolkit(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	obj.Spec.Template.Spec.Containers[0].Image = config.Toolkit.ImagePath()
	runtime := string(config.Operator.DefaultRuntime)

	setContainerEnv(&(obj.Spec.Template.Spec.Containers[0]), "RUNTIME", runtime)
	if runtime == "docker" {
		setContainerEnv(&(obj.Spec.Template.Spec.Containers[0]), "RUNTIME_ARGS",
			"--socket /var/run/docker.sock")
	}

	return nil
}

func TransformDevicePlugin(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	obj.Spec.Template.Spec.Containers[0].Image = config.DevicePlugin.ImagePath()

	return nil
}

func TransformDCGMExporter(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	obj.Spec.Template.Spec.Containers[0].Image = config.DCGMExporter.ImagePath()

	kvers, osTag := kernelFullVersion(n)
	if osTag == "" {
		return fmt.Errorf("ERROR: Could not find kernel full version: ('%s', '%s')", kvers, osTag)
	}

	if osTag != "rhel" && osTag != "centos" {
		return nil
	}

	initContainerImage, initContainerName, initContainerCmd := "ubuntu:18.04", "init-pod-nvidia-metrics-exporter", "/bin/entrypoint.sh"
	obj.Spec.Template.Spec.InitContainers[0].Image = initContainerImage
	obj.Spec.Template.Spec.InitContainers[0].Name = initContainerName
	obj.Spec.Template.Spec.InitContainers[0].Command[0] = initContainerCmd

	volMountSockName, volMountSockPath := "pod-gpu-resources", "/var/lib/kubelet/pod-resources"
	volMountSock := corev1.VolumeMount{Name: volMountSockName, MountPath: volMountSockPath}
	obj.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(obj.Spec.Template.Spec.InitContainers[0].VolumeMounts, volMountSock)

	volMountConfigName, volMountConfigPath, volMountConfigSubPath := "init-config", "/bin/entrypoint.sh", "entrypoint.sh"
	volMountConfig := corev1.VolumeMount{Name: volMountConfigName, ReadOnly: true, MountPath: volMountConfigPath, SubPath: volMountConfigSubPath}
	obj.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(obj.Spec.Template.Spec.InitContainers[0].VolumeMounts, volMountConfig)

	volMountConfigKey, volMountConfigDefaultMode := "nvidia-dcgm-exporter", int32(0700)
	initVol := corev1.Volume{Name: volMountConfigName, VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: volMountConfigKey}, DefaultMode: &volMountConfigDefaultMode}}}
	obj.Spec.Template.Spec.Volumes = append(obj.Spec.Template.Spec.Volumes, initVol)

	return nil
}

func setContainerEnv(c *corev1.Container, key, value string) {
	for i, val := range c.Env {
		if val.Name != key {
			continue
		}

		c.Env[i].Value = value
		return
	}

	log.Info(fmt.Sprintf("Info: Could not find environment variable %s in container %s, appending it", key, c.Name))
	c.Env = append(c.Env, corev1.EnvVar{Name: key, Value: value})
}

func isDaemonSetReady(name string, n ClusterPolicyController) gpuv1.State {
	opts := []client.ListOption{
		client.MatchingLabels{"app": name},
	}
	log.Info("DEBUG: DaemonSet", "LabelSelector", fmt.Sprintf("app=%s", name))
	list := &appsv1.DaemonSetList{}
	err := n.rec.client.List(context.TODO(), list, opts...)
	if err != nil {
		log.Info("Could not get DaemonSetList", err)
	}
	log.Info("DEBUG: DaemonSet", "NumberOfDaemonSets", len(list.Items))
	if len(list.Items) == 0 {
		return gpuv1.NotReady
	}

	ds := list.Items[0]
	log.Info("DEBUG: DaemonSet", "NumberUnavailable", ds.Status.NumberUnavailable)

	if ds.Status.NumberUnavailable != 0 {
		return gpuv1.NotReady
	}

	return isPodReady(name, n, "Running")
}

func DaemonSet(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].DaemonSet.DeepCopy()

	preProcessDaemonSet(obj, n)
	logger := log.WithValues("DaemonSet", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return isDaemonSetReady(obj.Name, n), nil
}

// The operator starts two pods in different stages to validate
// the correct working of the DaemonSets (driver and dp). Therefore
// the operator waits until the Pod completes and checks the error status
// to advance to the next state.
func isPodReady(name string, n ClusterPolicyController, phase corev1.PodPhase) gpuv1.State {
	opts := []client.ListOption{&client.MatchingLabels{"app": name}}

	log.Info("DEBUG: Pod", "LabelSelector", fmt.Sprintf("app=%s", name))
	list := &corev1.PodList{}
	err := n.rec.client.List(context.TODO(), list, opts...)
	if err != nil {
		log.Info("Could not get PodList", err)
	}
	log.Info("DEBUG: Pod", "NumberOfPods", len(list.Items))
	if len(list.Items) == 0 {
		return gpuv1.NotReady
	}

	pd := list.Items[0]

	if pd.Status.Phase != phase {
		log.Info("DEBUG: Pod", "Phase", pd.Status.Phase, "!=", phase)
		return gpuv1.NotReady
	}
	log.Info("DEBUG: Pod", "Phase", pd.Status.Phase, "==", phase)
	return gpuv1.Ready
}

func Pod(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].Pod.DeepCopy()
	logger := log.WithValues("Pod", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return isPodReady(obj.Name, n, "Succeeded"), nil
}

func SecurityContextConstraints(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].SecurityContextConstraints.DeepCopy()
	logger := log.WithValues("SecurityContextConstraints", obj.Name, "Namespace", "default")

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func Service(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].Service.DeepCopy()
	logger := log.WithValues("Service", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func ServiceMonitor(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ServiceMonitor.DeepCopy()
	logger := log.WithValues("ServiceMonitor", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}
