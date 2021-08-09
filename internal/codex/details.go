package codex

import (
	"context"
	"github.com/pathbird/pbauthor/internal/graphql"
	"github.com/pathbird/pbauthor/internal/graphql/transport"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	time "time"
)

type Details struct {
	Name       string
	KernelSpec KernelSpec `json:"kernelSpec"`
}

type KernelSpec struct {
	ID          string   `json:"id"`
	BuildStatus string   `json:"buildStatus"`
	Events      []string `json:"events"`
	BuildLog    []string `json:"buildLog"`
}

const waitForBuildStatusQuery = `
query pbauthor_CodexBuildStatus($id: ID!) {
	node(id: $id) { ... on CodexMetadata {
		id
		name
		kernelSpec {
			id
			buildStatus
			events
			buildLog
		}
	}}
}
`

// WaitForKernelBuildCompleted polls the API server until the kernel build is completed.
func WaitForKernelBuildCompleted(
	ctx context.Context,
	client *graphql.Client,
	codexId string,
) (*KernelSpec, error) {
	for {
		log.WithField("codex_id", codexId).Debug("querying kernel build status")
		status, err := queryKernelStatus(ctx, client, codexId)
		if err != nil {
			return nil, err
		}
		if status.BuildStatus != "pending" {
			return status, nil
		}

		// sleep for a few seconds, then try again
		select {
		case <-time.After(5 * time.Second):
			// pass
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func queryKernelStatus(
	ctx context.Context,
	client *graphql.Client,
	codexId string,
) (*KernelSpec, error) {
	req := transport.NewRequest(waitForBuildStatusQuery)
	req.Var("id", codexId)
	var res struct {
		Node struct {
			ID         string
			Name       string
			KernelSpec KernelSpec
		}
	}
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	if res.Node.ID == "" {
		return nil, errors.Errorf("codex (id: %s) could not be found", codexId)
	}
	if res.Node.KernelSpec.ID == "" {
		return nil, errors.Errorf("query didn't return KernelSpec data (for codex: %s)", codexId)
	}
	return &res.Node.KernelSpec, nil
}
