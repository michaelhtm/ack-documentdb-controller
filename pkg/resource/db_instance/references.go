// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package db_instance

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kmsapitypes "github.com/aws-controllers-k8s/kms-controller/apis/v1alpha1"
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"

	svcapitypes "github.com/aws-controllers-k8s/documentdb-controller/apis/v1alpha1"
)

// +kubebuilder:rbac:groups=kms.services.k8s.aws,resources=keys,verbs=get;list
// +kubebuilder:rbac:groups=kms.services.k8s.aws,resources=keys/status,verbs=get;list

// ClearResolvedReferences removes any reference values that were made
// concrete in the spec. It returns a copy of the input AWSResource which
// contains the original *Ref values, but none of their respective concrete
// values.
func (rm *resourceManager) ClearResolvedReferences(res acktypes.AWSResource) acktypes.AWSResource {
	ko := rm.concreteResource(res).ko.DeepCopy()

	if ko.Spec.PerformanceInsightsKMSKeyRef != nil {
		ko.Spec.PerformanceInsightsKMSKeyID = nil
	}

	return &resource{ko}
}

// ResolveReferences finds if there are any Reference field(s) present
// inside AWSResource passed in the parameter and attempts to resolve those
// reference field(s) into their respective target field(s). It returns a
// copy of the input AWSResource with resolved reference(s), a boolean which
// is set to true if the resource contains any references (regardless of if
// they are resolved successfully) and an error if the passed AWSResource's
// reference field(s) could not be resolved.
func (rm *resourceManager) ResolveReferences(
	ctx context.Context,
	apiReader client.Reader,
	res acktypes.AWSResource,
) (acktypes.AWSResource, bool, error) {
	ko := rm.concreteResource(res).ko

	resourceHasReferences := false
	err := validateReferenceFields(ko)
	if fieldHasReferences, err := rm.resolveReferenceForPerformanceInsightsKMSKeyID(ctx, apiReader, ko); err != nil {
		return &resource{ko}, (resourceHasReferences || fieldHasReferences), err
	} else {
		resourceHasReferences = resourceHasReferences || fieldHasReferences
	}

	return &resource{ko}, resourceHasReferences, err
}

// validateReferenceFields validates the reference field and corresponding
// identifier field.
func validateReferenceFields(ko *svcapitypes.DBInstance) error {

	if ko.Spec.PerformanceInsightsKMSKeyRef != nil && ko.Spec.PerformanceInsightsKMSKeyID != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("PerformanceInsightsKMSKeyID", "PerformanceInsightsKMSKeyRef")
	}
	return nil
}

// resolveReferenceForPerformanceInsightsKMSKeyID reads the resource referenced
// from PerformanceInsightsKMSKeyRef field and sets the PerformanceInsightsKMSKeyID
// from referenced resource. Returns a boolean indicating whether a reference
// contains references, or an error
func (rm *resourceManager) resolveReferenceForPerformanceInsightsKMSKeyID(
	ctx context.Context,
	apiReader client.Reader,
	ko *svcapitypes.DBInstance,
) (hasReferences bool, err error) {
	if ko.Spec.PerformanceInsightsKMSKeyRef != nil && ko.Spec.PerformanceInsightsKMSKeyRef.From != nil {
		hasReferences = true
		arr := ko.Spec.PerformanceInsightsKMSKeyRef.From
		if arr.Name == nil || *arr.Name == "" {
			return hasReferences, fmt.Errorf("provided resource reference is nil or empty: PerformanceInsightsKMSKeyRef")
		}
		namespace := ko.ObjectMeta.GetNamespace()
		if arr.Namespace != nil && *arr.Namespace != "" {
			namespace = *arr.Namespace
		}
		obj := &kmsapitypes.Key{}
		if err := getReferencedResourceState_Key(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
			return hasReferences, err
		}
		ko.Spec.PerformanceInsightsKMSKeyID = (*string)(obj.Status.ACKResourceMetadata.ARN)
	}

	return hasReferences, nil
}

// getReferencedResourceState_Key looks up whether a referenced resource
// exists and is in a ACK.ResourceSynced=True state. If the referenced resource does exist and is
// in a Synced state, returns nil, otherwise returns `ackerr.ResourceReferenceTerminalFor` or
// `ResourceReferenceNotSyncedFor` depending on if the resource is in a Terminal state.
func getReferencedResourceState_Key(
	ctx context.Context,
	apiReader client.Reader,
	obj *kmsapitypes.Key,
	name string, // the Kubernetes name of the referenced resource
	namespace string, // the Kubernetes namespace of the referenced resource
) error {
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	err := apiReader.Get(ctx, namespacedName, obj)
	if err != nil {
		return err
	}
	var refResourceSynced, refResourceTerminal bool
	for _, cond := range obj.Status.Conditions {
		if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
			cond.Status == corev1.ConditionTrue {
			refResourceSynced = true
		}
		if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
			cond.Status == corev1.ConditionTrue {
			return ackerr.ResourceReferenceTerminalFor(
				"Key",
				namespace, name)
		}
	}
	if refResourceTerminal {
		return ackerr.ResourceReferenceTerminalFor(
			"Key",
			namespace, name)
	}
	if !refResourceSynced {
		return ackerr.ResourceReferenceNotSyncedFor(
			"Key",
			namespace, name)
	}
	if obj.Status.ACKResourceMetadata == nil || obj.Status.ACKResourceMetadata.ARN == nil {
		return ackerr.ResourceReferenceMissingTargetFieldFor(
			"Key",
			namespace, name,
			"Status.ACKResourceMetadata.ARN")
	}
	return nil
}
