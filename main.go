package main

import (
	"context"
	"fmt"
)

const (
	GH_CLI_VERSION = "2.38.0"
	RUNTIME = "alpine:3.18.4"
)

type GithubRelease struct {}

func (g *GithubRelease) CreateWithAssets(
		ctx context.Context,
		repo string,
		tag string,
		title string,
		token *Secret,
		assets *Directory,
		notes Optional[string],
		draft Optional[bool],
		latest Optional[bool],
		prerelease Optional[bool]
	) (string, error) {
		createCmd := []string{"release", "create", tag}

		notes_, isset := notes.Get()
		if isset {
			createCmd = append(createCmd, fmt.Sprintf("'%s'", notes_))
		}

		if draft.GetOr(false) {
			createCmd = append(createCmd, "--draft")
		}

		if latest.GetOr(false) {
			createCmd = append(createCmd, "--latest")
		}

		if prerelease.GetOr(false) {
			createCmd = append(createCmd, "--prerelease")
		}
		return ghImage().
		WithSecretVariable("GH_TOKEN", token).
		WithEnvVariable("GH_REPO", repo).
		WithMountedDirectory("/assets", assets).
		WithExec(createCmd).
		WithExec([]string{"release", "upload", tag, "/assets/"}).
		Stdout(ctx)
}

func ghImage() *Container {
	return dag.Container().
	From(RUNTIME).
	WithFile("/bin/gh", dag.Gh().Get(GH_CLI_VERSION)).
	WithEntrypoint([]string{"/bin/gh"})
}
