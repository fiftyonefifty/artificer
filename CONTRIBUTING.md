# Contributing

## Contribution Process

### 1. Create your feature branch

```bash
git checkout -b my-new-feature
```

### 2. Commit your changes

If you have multiple commits, squash merge them into a single commit. This helps keep the history light. [Here's a nice guide](https://www.devroom.io/2011/07/05/git-squash-your-latests-commits-into-one/). This project uses [bumpversion](https://github.com/peritus/bumpversion) to manage the version number. Please install it to manage the version number.

```bash
git commit -am 'Added some feature'
```

### 3. Bump the version

Before pushing you change, please ensure that you bump the version. To do so install [bump2version](https://github.com/c4urself/bump2version). We use Semantic Versioning, so ensure that you bump accordingly.

```bash
# Example Major Change
bump2version major

# Example Minor Change
bump2version minor

# Example Patch (Hotfix) Change
bump2version patch
```

### 4. Push to the feature branch

```bash
git push origin my-new-feature
```

### 5. Create a Pull Request

## Contribution Guidelines

1. Coding style must match what exists in the project.
2. Must not contain any merge conflicts.
3. Include tests for the modifications in your commit.
