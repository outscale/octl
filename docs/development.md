# Development processes

## Release

1. In [CHANGELOG.md](../CHANGELOG.md), add a new vX.Y.Z version
2. Update [actions/setup_octl/action.yaml](../actions/setup_octl/action.yaml)
 - Set `version.default` to the new release
3. Open a pull request with the above changes and merge it once approved.
4. Tag and push the container release:
```shell
export VERSION=vX.Y.Z
git tag -a $VERSION -m "🔖 octl $VERSION" && git push origin $VERSION
```
5. Publish the Github release
6. Merge the [generated homebrew-tap PR](https://github.com/outscale/homebrew-tap/pulls)
