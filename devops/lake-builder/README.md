# lake-builder

golang builder image for lake, including go1.19.0 gcc10.3.1 g++10.3.1, libgit2 v1.3 based on alpine linux 3.16.

https://hub.docker.com/r/mericodev/lake-builder

## release

```shell
export VERSION=0.0.2
docker build -t mericodev/lake-builder:$VERSION .
docker push mericodev/lake-builder:$VERSION
```