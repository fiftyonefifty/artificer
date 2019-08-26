# Artificer

![Artificer Logo](/assets/images/cover.png)

Artificer is an OAuth2 token creator. Nothing more. Nothing less.

## Developing

This application is backed by the Echo framework. To get started with development, you'll need to restore the build locally.

```bash
go mod download
```

Once all dependencies are installed, you can start the server easily.

```bash
go run ./cmd/artificer
```

If all goes well, Artificer will be available at `localhost:8000`.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Versioning

This project uses semantic versioning ([SemVer 2.0.0](https://semver.org/)). Incrementing versions is managed by [bumpversion](https://github.com/peritus/bumpversion).

To ensure that the repo is properly versioned, you will need to install `bumpversion`.

```bash
pip install bumpversion
```

Once installed, bump the version before pushing your code or created a pull request.

```bash
# Examples

# Bumping the major version to indicate a backwards incompatible change
bumpversion major

# Bumping the minor version
bumpversion minor

# Bumping the subminor due to a hotfix
bumpversion patch
```

*Note: Bumpversion is configured to automatically create a commit when executed.*
