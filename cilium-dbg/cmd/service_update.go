// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cilium/cilium/api/v1/models"
	cmtypes "github.com/cilium/cilium/pkg/clustermesh/types"
	"github.com/cilium/cilium/pkg/loadbalancer"
	"github.com/cilium/cilium/pkg/loadbalancer/legacy/service"
)

var (
	k8sExternalIPs      bool
	k8sNodePort         bool
	k8sHostPort         bool
	k8sLoadBalancer     bool
	k8sExtTrafficPolicy string
	k8sIntTrafficPolicy string
	k8sClusterInternal  bool
	localRedirect       bool
	idU                 uint64
	frontend            string
	protocol            string
	backends            []string
	backendStates       []string
	backendWeights      []uint
)

func warnIdTypeDeprecation() {
	fmt.Printf("Deprecation warning: --id parameter will change from int to string in v1.14\n")
}

// serviceUpdateCmd represents the service_update command
var serviceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a service",
	Run: func(cmd *cobra.Command, args []string) {
		updateService(cmd, args)
	},
}

func init() {
	ServiceCmd.AddCommand(serviceUpdateCmd)
	serviceUpdateCmd.Flags().Uint64VarP(&idU, "id", "", 0, "Identifier")
	serviceUpdateCmd.Flags().BoolVarP(&k8sExternalIPs, "k8s-external", "", false, "Set service as a k8s ExternalIPs")
	serviceUpdateCmd.Flags().BoolVarP(&k8sNodePort, "k8s-node-port", "", false, "Set service as a k8s NodePort")
	serviceUpdateCmd.Flags().BoolVarP(&k8sLoadBalancer, "k8s-load-balancer", "", false, "Set service as a k8s LoadBalancer")
	serviceUpdateCmd.Flags().BoolVarP(&k8sHostPort, "k8s-host-port", "", false, "Set service as a k8s HostPort")
	serviceUpdateCmd.Flags().BoolVarP(&localRedirect, "local-redirect", "", false, "Set service as Local Redirect")
	serviceUpdateCmd.Flags().StringVarP(&k8sExtTrafficPolicy, "k8s-ext-traffic-policy", "", "Cluster", "Set service with k8s externalTrafficPolicy as {Local,Cluster}")
	serviceUpdateCmd.Flags().StringVarP(&k8sIntTrafficPolicy, "k8s-int-traffic-policy", "", "Cluster", "Set service with k8s internalTrafficPolicy as {Local,Cluster}")
	serviceUpdateCmd.Flags().BoolVarP(&k8sClusterInternal, "k8s-cluster-internal", "", false, "Set service as cluster-internal for externalTrafficPolicy=Local xor internalTrafficPolicy=Local")
	serviceUpdateCmd.Flags().StringVarP(&frontend, "frontend", "", "", "Frontend address")
	serviceUpdateCmd.Flags().StringVarP(&protocol, "protocol", "", "tcp", "Protocol for service (e.g. TCP, UDP)")
	serviceUpdateCmd.Flags().StringSliceVarP(&backends, "backends", "", []string{}, "Backend address or addresses (<IP:Port>)")
	serviceUpdateCmd.Flags().StringSliceVarP(&backendStates, "states", "", []string{}, "Backend state(s) as {active(default),terminating,quarantined,maintenance}")
	serviceUpdateCmd.Flags().UintSliceVarP(&backendWeights, "backend-weights", "", []uint{}, "Backend weights (100 default, 0 means maintenance state, only for maglev mode)")
}

func parseAddress(l4Protocol, address string) (ip net.IP, port int, proto string, err error) {
	switch proto = strings.ToLower(l4Protocol); proto {
	case "tcp":
		var tcpAddr *net.TCPAddr
		tcpAddr, err = net.ResolveTCPAddr(proto, address)
		if err != nil {
			return
		}
		ip = tcpAddr.IP
		port = tcpAddr.Port
	case "udp":
		var udpAddr *net.UDPAddr
		udpAddr, err = net.ResolveUDPAddr(proto, address)
		if err != nil {
			return
		}
		ip = udpAddr.IP
		port = udpAddr.Port
	default:
		err = fmt.Errorf("unrecognized protocol %q", l4Protocol)
	}
	return
}

func parseFrontendAddress(l4Protocol, address string) *models.FrontendAddress {
	ip, port, proto, err := parseAddress(l4Protocol, address)
	if err != nil {
		Fatalf("Unable to parse frontend address: %s\n", err)
	}

	scope := models.FrontendAddressScopeExternal
	if k8sClusterInternal {
		scope = models.FrontendAddressScopeInternal
	}

	return &models.FrontendAddress{
		IP:       ip.String(),
		Port:     uint16(port),
		Protocol: proto,
		Scope:    scope,
	}
}

func boolToInt(set bool) int {
	if set {
		return 1
	}
	return 0
}

func updateService(cmd *cobra.Command, args []string) {
	warnIdTypeDeprecation()

	id := int64(idU)
	maxID := int64(service.MaxSetOfServiceID)
	if id > maxID {
		Fatalf("Service ID %d exceeds the maximum limit of %d", id, maxID)
	}

	fa := parseFrontendAddress(protocol, frontend)
	skipFrontendCheck := false

	var spec *models.ServiceSpec
	svc, err := client.GetServiceID(id)
	switch {
	case id == 0 && frontend == "" && len(backends) != 0:
		// When service ID is 0 and frontend is not specified, the intended use
		// of the API is to update backend state(s) for service(s) selecting those
		// backend(s).
		if len(backendStates) == 0 {
			Fatalf("Cannot update empty backend states")
		}
		spec = &models.ServiceSpec{ID: 0}
		skipFrontendCheck = true
		spec.UpdateServices = true
		fmt.Printf("Updating backend states \n")

	case err == nil && (svc.Status == nil || svc.Status.Realized == nil):
		Fatalf("Cannot update service %d: empty state", id)

	case err == nil:
		spec = svc.Status.Realized
		fmt.Printf("Updating existing service with id '%v'\n", id)

	default:
		spec = &models.ServiceSpec{ID: id}
		fmt.Printf("Creating new service with id '%v'\n", id)
	}

	// This can happen when we create a new service or when the service returned
	// to us has no flags set
	if spec.Flags == nil {
		spec.Flags = &models.ServiceSpecFlags{}
	}

	if boolToInt(k8sExternalIPs)+boolToInt(k8sNodePort)+boolToInt(k8sHostPort)+boolToInt(k8sLoadBalancer)+boolToInt(localRedirect) > 1 {
		Fatalf("Can only set one of --k8s-external, --k8s-node-port, --k8s-load-balancer, --k8s-host-port, --local-redirect for a service")
	} else if k8sExternalIPs {
		spec.Flags = &models.ServiceSpecFlags{Type: models.ServiceSpecFlagsTypeExternalIPs}
	} else if k8sNodePort {
		spec.Flags = &models.ServiceSpecFlags{Type: models.ServiceSpecFlagsTypeNodePort}
	} else if k8sLoadBalancer {
		spec.Flags = &models.ServiceSpecFlags{Type: models.ServiceSpecFlagsTypeLoadBalancer}
	} else if k8sHostPort {
		spec.Flags = &models.ServiceSpecFlags{Type: models.ServiceSpecFlagsTypeHostPort}
	} else if localRedirect {
		spec.Flags = &models.ServiceSpecFlags{Type: models.ServiceSpecFlagsTypeLocalRedirect}
	} else {
		spec.Flags = &models.ServiceSpecFlags{Type: models.ServiceSpecFlagsTypeClusterIP}
	}

	if strings.ToLower(k8sExtTrafficPolicy) == "local" {
		spec.Flags.TrafficPolicy = models.ServiceSpecFlagsTrafficPolicyLocal
		spec.Flags.ExtTrafficPolicy = models.ServiceSpecFlagsExtTrafficPolicyLocal
	} else {
		spec.Flags.TrafficPolicy = models.ServiceSpecFlagsTrafficPolicyCluster
		spec.Flags.ExtTrafficPolicy = models.ServiceSpecFlagsExtTrafficPolicyCluster
	}

	if strings.ToLower(k8sIntTrafficPolicy) == "local" {
		spec.Flags.IntTrafficPolicy = models.ServiceSpecFlagsIntTrafficPolicyLocal
	} else {
		spec.Flags.IntTrafficPolicy = models.ServiceSpecFlagsIntTrafficPolicyCluster
	}

	spec.FrontendAddress = fa

	if len(backends) == 0 {
		fmt.Printf("Reading backend list from stdin...\n")

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			backends = append(backends, scanner.Text())
		}
	}

	if len(backendWeights) > 0 {
		resp, err := client.ConfigGet()
		if err != nil {
			Fatalf("Unable to retrieve cilium configuration: %s", err)
		}
		if resp.Status == nil {
			Fatalf("Unable to retrieve cilium configuration: empty response")
		}
		if len(backendWeights) != len(backends) {
			Fatalf("Mismatch between number of backend weights and number of backends")
		}
	}

	spec.BackendAddresses = nil
	backendState, _ := loadbalancer.BackendStateActive.String()

	switch {
	case len(backendStates) == 0:
	case len(backendStates) == 1:
		backendState = backendStates[0]
		if !loadbalancer.IsValidBackendState(backendState) {
			Fatalf("Invalid backend state (%v)", backendState)
		}
	case len(backendStates) == len(backends):
	default:
		Fatalf("Invalid number of backend states (%v) for backends (%v)", backendStates, backends)
	}

	for i, backend := range backends {
		ip, port, proto, err := parseAddress(protocol, backend)
		if err != nil {
			Fatalf("Cannot parse backend address %q: %s", backend, err)
		}
		// Backend ID will be set by the daemon
		be := loadbalancer.NewLegacyBackend(0, loadbalancer.L4Type(strings.ToUpper(proto)), cmtypes.MustAddrClusterFromIP(ip), uint16(port))

		if !skipFrontendCheck && fa.Port == 0 && port != 0 {
			Fatalf("L4 backend found (%s:%d) with L3 frontend", ip, port)
		}

		ba := be.GetBackendModel()

		if i < len(backendStates) {
			if !loadbalancer.IsValidBackendState(backendStates[i]) {
				Fatalf("Invalid backend state (%v) for backend (%v)", backendStates[i], backends[i])
			}
			ba.State = backendStates[i]
		} else {
			ba.State = backendState
		}

		if i < len(backendWeights) {
			w := uint16(backendWeights[i])
			ba.Weight = &w
		}

		spec.BackendAddresses = append(spec.BackendAddresses, ba)
	}

	if created, err := client.PutServiceID(id, spec); err != nil {
		Fatalf("Cannot add/update service: %s %+v", err, spec)
	} else if created {
		fmt.Printf("Added service with %d backends\n", len(spec.BackendAddresses))
	} else {
		fmt.Printf("Updated service with %d backends\n", len(spec.BackendAddresses))
	}
}
