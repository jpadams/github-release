package main

import (
	"context"
)

const (
	GH_CLI_VERSION = "2.38.0"
	RUNTIME = "alpine:3.18.4"
)

type GithubRelease struct {}

func (g *GithubRelease) Create(
		ctx context.Context,
		repo string,
		tag string,
		title string,
		token *Secret,
		assets Optional[*Directory],
		notes Optional[string],
		draft Optional[bool],
		latest Optional[bool],
		prerelease Optional[bool],
	) (string, error) {
		createCmd := []string{"release", "create", tag, "--title", title}

		notes_, isset := notes.Get()
		if isset {
			createCmd = append(createCmd, "--notes", notes_)
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

		releaser := ghImage().
		WithSecretVariable("GH_TOKEN", token).
		WithEnvVariable("GH_REPO", repo).
		WithExec(createCmd)

		assets_, isset := assets.Get()
		if isset {
			releaser = releaser.
			WithMountedDirectory("/assets", assets).
			WithExec([]string{"release", "upload", tag, "/assets/*"})
		}

		return releaser.Stdout(ctx)
}

func ghImage() *Container {
	return dag.Container().
	From(RUNTIME).
	WithFile("/bin/gh", dag.Gh().Get(GhGetOpts{ Version: GH_CLI_VERSION})).
	WithEntrypoint([]string{"/bin/gh"})
}
