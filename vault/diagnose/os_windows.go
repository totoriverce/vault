// +build windows

package diagnose

import (
	"context"
)

func OSChecks(ctx context.Context) {
	ctx, span := StartSpan(ctx, "Operating System")
	defer span.End()
	diskUsage(ctx)
}
