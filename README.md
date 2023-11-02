## Examples

Create a github release
`dagger call -m github.com/jpadams/github-release create --repo github.com/foo/bar --tag v0.0.0 --title 'my release' --token $GH_TOKEN`

Create a github release with artifacts
`dagger call -m github.com/jpadams/github-release create --repo github.com/foo/bar --tag v0.0.0 --title 'my release' --token $GH_TOKEN --assets ./build`

Create a draft release with notes
`dagger call -m github.com/jpadams/github-release create --repo github.com/foo/bar --tag v0.0.0 --title 'my release' --token $GH_TOKEN --assets ./build --notes 'my release notes' --draft`

