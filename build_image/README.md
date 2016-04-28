# Build image

Docker image for building statically linked go applications

It calls `glide install` and `go build` all gitlab repositories have to have
`gitlab_public` deploy key allowed

## Build
```
make build
```

## Push to registry
```
make push
```

## Usage
```
PROJECT_NAME=github.com/sejvlond/tarsier
PROJECT_SOURCE=$GOPATH/src/$PROJECT_NAME

docker run --rm -v $PROJECT_SOURCE:/src sejvlond/tarsier_build $PROJECT_NAME
