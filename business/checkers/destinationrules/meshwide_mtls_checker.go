package destinationrules

import (
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/models"
)

type MeshWideMTLSChecker struct {
	DestinationRule kubernetes.IstioObject
	MTLSDetails     kubernetes.MTLSDetails
}

func (m MeshWideMTLSChecker) Check() ([]*models.IstioCheck, bool) {
	validations := make([]*models.IstioCheck, 0)

	// if DestinationRule doesn't enable mTLS, stop validation with any check
	if enabled, _ := kubernetes.DestinationRuleHasMeshWideMTLSEnabled(m.DestinationRule); !enabled {
		return validations, true
	}

	// otherwise, check among MeshPeerAuthentications for a rule enabling mesh-wide mTLS
	// ServiceMeshPolicies are a clone of MeshPeerAuthentications but used in Maistra scenarios
	// MeshPeerAuthentications and ServiceMeshPolicies won't co-exist, only ony array will be populated
	mPolicies := m.MTLSDetails.MeshPeerAuthentications
	checkerId := "destinationrules.mtls.meshpolicymissing"
	if m.MTLSDetails.ServiceMeshPolicies != nil {
		mPolicies = m.MTLSDetails.ServiceMeshPolicies
		checkerId = "destinationrules.mtls.servicemeshpolicymissing"
	}
	for _, mp := range mPolicies {
		if enabled, _ := kubernetes.PeerAuthnHasMTLSEnabled(mp); enabled {
			return validations, true
		}
	}

	check := models.Build(checkerId, "spec/trafficPolicy/tls/mode")
	validations = append(validations, &check)

	return validations, false
}
