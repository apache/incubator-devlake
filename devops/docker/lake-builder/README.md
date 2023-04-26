# lake-builder

Golang builder image for DevLake. Used by GitHub workflows.

https://hub.docker.com/r/mericodev/lake-builder

## Manual Release

```shell
export VERSION=0.0.11
docker build -t mericodev/lake-builder:$VERSION .
docker push mericodev/lake-builder:$VERSION
```

## Tagged Release
1. Create a tag matching the pattern `builder-v#.#.#`, e.g. `builder-v0.0.11`. Determine the previous tag first so you version
it correctly. Example command: `git tag builder-v0.0.11`
2. Push the tag to origin. Example command: `git push origin --tag builder-v0.0.11`
3. Done! This will trigger a GitHub workflow that will push this image to the Docker registry.