package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
)

func newResource() *sdkresource.Resource {
	var (
		newResourcesOnce sync.Once
		resource         *sdkresource.Resource
	)

	newResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
			sdkresource.WithAttributes(
				attribute.String("environment", "demo"),
			),
		)

		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)

	})
	return resource
}
