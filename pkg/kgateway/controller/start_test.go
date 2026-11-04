package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	apisettings "github.com/kgateway-dev/kgateway/v2/api/settings"
	"github.com/kgateway-dev/kgateway/v2/pkg/deployer"
	"github.com/kgateway-dev/kgateway/v2/pkg/kgateway/wellknown"
)

func TestGetDefaultClassInfoAppliesParametersRefs(t *testing.T) {
	t.Parallel()

	const standardClass = "kgateway"

	ns := gwv1.Namespace("control-plane")
	settings := &apisettings.Settings{
		GatewayClassParametersRefs: apisettings.GatewayClassParametersRefs{
			standardClass: {
				Name:      "standard-gwp",
				Namespace: &ns,
			},
		},
	}

	classInfos := GetDefaultClassInfo(
		settings,
		standardClass,
		"ctrl.kgateway.dev",
		nil,
	)

	require.NotNil(t, classInfos[standardClass])
	require.NotNil(t, classInfos[standardClass].ParametersRef)
	require.Equal(t, "standard-gwp", classInfos[standardClass].ParametersRef.Name)
	require.Equal(t, gwv1.Namespace("control-plane"), *classInfos[standardClass].ParametersRef.Namespace)
	require.Equal(t, gwv1.Group(wellknown.GatewayParametersGVK.Group), classInfos[standardClass].ParametersRef.Group)
	require.Equal(t, gwv1.Kind(wellknown.GatewayParametersGVK.Kind), classInfos[standardClass].ParametersRef.Kind)
}

func TestGetDefaultClassInfoManagedToggle(t *testing.T) {
	t.Parallel()

	const standardClass = "kgateway"

	managed := GetDefaultClassInfo(&apisettings.Settings{ManageDefaultGatewayClass: true}, standardClass, "ctrl.kgateway.dev", nil)
	require.NotNil(t, managed[standardClass])
	require.False(t, managed[standardClass].Unmanaged, "default class should be managed when enabled")

	unmanaged := GetDefaultClassInfo(&apisettings.Settings{ManageDefaultGatewayClass: false}, standardClass, "ctrl.kgateway.dev", nil)
	require.NotNil(t, unmanaged[standardClass], "class info is kept so status is still reported")
	require.True(t, unmanaged[standardClass].Unmanaged, "default class should be unmanaged when disabled")
	require.NotEmpty(t, unmanaged[standardClass].SupportedFeatures, "supported features are still reported")
}

func TestReconcileGatewayClassSkipsUnmanaged(t *testing.T) {
	t.Parallel()

	// No client is configured, so any apply attempt would panic.
	r := &gatewayClassReconciler{defaultControllerName: "ctrl.kgateway.dev"}
	require.NoError(t, r.reconcileGatewayClass("kgateway", &deployer.GatewayClassInfo{Unmanaged: true}))
}
