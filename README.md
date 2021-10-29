[![Build Status](https://cloud.drone.io/api/badges/13rentgen/nexus-cli/status.svg)](https://cloud.drone.io/13rentgen/nexus-cli) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

<div align="center">
<img src="logo.png" width="60%"/>
</div>

Nexus CLI for Docker Registry

## Usage

<div align="center">
<img src="example.png"/>
</div>

## Getting started

To run `nexus-cli` download the [latest release](https://github.com/13rentgen/nexus-cli/releases/latest) distribution.

## Building from Source
Build the binary using `make`:
```
make build
```

## Configure
For configure `nexus-cli` use `nexus-cli configure`.

⚠️ Important: For the `Enter Nexus Host` question type your Nexus host with *web port*. Typically, Nexus run on port `8081`.  

## Available Commands

```
$ nexus-cli configure
```

```
$ nexus-cli image ls
```

```
$ nexus-cli image tags -name mlabouardy/nginx
```

```
$ nexus-cli image info -name mlabouardy/nginx -tag 1.2.0
```

```
$ nexus-cli image delete -name mlabouardy/nginx -tag 1.2.0
```

```
$ nexus-cli image delete -name mlabouardy/nginx -keep 4
```

```
$ nexus-cli image size -name mlabouardy/nginx
```
## Tutorials

* [Cleanup old Docker images from Nexus Repository](http://www.blog.labouardy.com/cleanup-old-docker-images-from-nexus-repository/)

## License

* [MIT License](LICENSE)
