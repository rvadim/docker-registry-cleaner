
# Docker Registry Cleaner

## Description  
Simple app to delete images from your private docker registry

### Build
```
go get
go build
```

Linux
```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .
```

### Docker Build


Docker
```
docker build -t docker-registry-cleaner .

```

or 

The Dockerfile.build in this project uses a multi stage docker build which requires docker >= 17.05 

```
docker build -t docker-registry-cleaner -f Dockerfile.build .
```

### Example  
CLI  
```
docker-registry-cleaner --path your-registry/repo/path --keep 3 --imageversion "^1.0.*" --dry-run
```

Docker
```
docker run -it -e PATH=your-registry/repo/path -e IMAGE_VERSION=".*-TEST" -e DRYRUN=true this-image:latest
```

## Documentation
Simple app that hits the docker registry api to delete images.  This is usefully if you want to only keep x number of the latest versions of an image.

### options
Available command line optoins  

```
   --path value                      Regular docker path '<registry>/repo/image' [$PATH]
   --username value, -u value        Registry username (optional) [$USERNAME]
   --password value, -p value        Registry password (optional) [$PASSWORD]
   --imageversion value, --iv value  Image Version to delete, this can be a regex ".*-SNAPSHOT.*" (default: ".*-SNAPSHOT.*") [$IMAGE_VERSION]
   --keep value, -k value            The number of images you want to keep, usefully if you are deleting images by regex (default: 3) [$KEEP]
   --dryrun, -d                      Do not actually delete anything [$DRYRUN]
   --help, -h                        show help
   --version, -v                     print the version
```
