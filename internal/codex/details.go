package codex

import (
	"context"
	"github.com/pathbird/pbauthor/internal/graphql"
	"github.com/pathbird/pbauthor/internal/graphql/transport"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
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
query pbauthor_CodexBuildStatus($id: ID!, $offset: Int, $limit: Int) {
	node(id: $id) { ... on CodexMetadata {
		id
		name
		kernelSpec {
			id
			buildStatus
			events
			buildLog(offset: $offset, limit: $limit)
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
	// If the kernel is built right away, it indicates that we're using a previous kernel build
	// so we can skip waiting for the build to complete (and in particular we don't want to
	// write the buildlogs to stdout again).
	spec, err := queryKernelSpec(ctx, client, codexId, 0, 0)
	if err != nil {
		return nil, err
	}
	if spec.BuildStatus != "pending" {
		log.Info("kernel image already exists (using cached image)")
		return spec, nil
	}

	var offset int64
	for {
		log.WithField("codex_id", codexId).Debug("querying kernel build status")
		spec, err := queryKernelSpec(ctx, client, codexId, offset, 100)
		if err != nil {
			return nil, err
		}
		offset += int64(len(spec.BuildLog))
		for _, logentry := range spec.BuildLog {
			_, _ = os.Stdout.WriteString(logentry)
		}
		if spec.BuildStatus != "pending" {
			return spec, nil
		}

		// sleep for a few seconds, then try again
		select {
		case <-time.After(2 * time.Second):
			// pass
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func queryKernelSpec(
	ctx context.Context,
	client *graphql.Client,
	codexId string,
	offset int64,
	length int64,
) (*KernelSpec, error) {
	req := transport.NewRequest(waitForBuildStatusQuery)
	req.Var("id", codexId)
	req.Var("offset", offset)
	req.Var("limit", length)
	var res struct {
		Node struct {
			ID         string
			Name       string
			KernelSpec KernelSpec
		}
	}
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}
	if res.Node.ID == "" {
		return nil, errors.Errorf("codex (id: %s) could not be found", codexId)
	}
	if res.Node.KernelSpec.ID == "" {
		return nil, errors.Errorf("query didn't return KernelSpec data (for codex: %s)", codexId)
	}
	return &res.Node.KernelSpec, nil
}
