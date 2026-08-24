// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package watchhandlers

import (
	"context"
	"log/slog"

	"github.com/cilium/cilium/pkg/logging/logfields"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gateway_inf_ext "sigs.k8s.io/gateway-api-inference-extension/api/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// EnqueueRequestForOwningInferencePool returns an event handler that, when passed a InferencePool, returns reconcile.Requests
// for all Cilium-relevant Gateways associated with that InferencePool.
func EnqueueRequestForOwningInferencePool(c client.Client, logger *slog.Logger, controllerName string) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, a client.Object) []reconcile.Request {
		// get all the httproutes that the inferencepool refers to
		inferencePool, ok := a.(*gateway_inf_ext.InferencePool)
		httpRouteList := &gatewayv1.HTTPRoute{}

		if err := c.List(ctx, httpRouteList); err != nil {
			scopedLog.WarnContext(ctx, "Unable to list httproutes", logfields.Error, err)
		}

		httprouteInfExt:= make(map[&gatewayv1.HTTPRoute]struct{})

		// for each httproute, check if the rules.backendref is the inferencepool
		for hr := range httpRouteList {
			for _, backend := range hr.Spec.Rules {
				if backend.kind == "InferencePool" && backend.name == inferencePool.Metadata.Name {
					// if it is a match, then put the httproute in a list
					httprouteInfExt[hr] = struct{}{}
				}else{
					continue
				}
			}
		}
		// get all the gateways associated with the httproutes
	})
}
