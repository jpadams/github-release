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
		// +optional
		assets *Directory,
		// +optional
		notes string,
		// +optional
		draft bool,
		// +optional
		latest bool,
		// +optional
		prerelease bool,
	) (string, error) {
		createCmd := []string{"release", "create", tag, "--title", title}

		if notes != "" {
			createCmd = append(createCmd, "--notes", notes)
		}

		if draft {
			createCmd = append(createCmd, "--draft")
		}

		if latest {
			createCmd = append(createCmd, "--latest")
		}

		if prerelease {
			createCmd = append(createCmd, "--prerelease")
		}

		releaser := ghImage().
		WithSecretVariable("GH_TOKEN", token).
		WithEnvVariable("GH_REPO", repo).
		WithExec(createCmd)

		if assets != nil {
			entries, err := assets.Entries(ctx)
			if err != nil {
				return "", err
			}

			uploadCmd := append([]string{"release", "upload", tag}, entries...)

			releaser = releaser.
			WithMountedDirectory("/assets", assets).
			WithWorkdir("/assets").
			WithExec(uploadCmd)
		}

		return releaser.Stdout(ctx)
}

func ghImage() *Container {
	return dag.Container().
	From(RUNTIME).
	WithFile("/bin/gh", dag.Gh().Get(GhGetOpts{ Version: GH_CLI_VERSION})).
	WithEntrypoint([]string{"/bin/gh"})
}
