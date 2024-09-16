# Introduction

This project is designed to provide a comprehensive solution for sorting and generating data using Go and Docker. The project includes several Makefile targets to streamline the development and deployment process. Below is a detailed explanation of each Makefile target and its usage.

# How to run project?

## To run on your machine

```bash
make gen in=10000
make sort in=10000
```

## To run with Docker

```bash
make gen in=10000
make limited in=10000
```

# Full makefile documentation:

This Makefile helps with running Go sorting and generation commands, managing Docker containers, and cleaning generated files. The following targets are defined:

## Variables

- `mem`: Memory limit for the Docker container (default: `512m`).
- `in`: Input file identifier (default: `1`).

## Targets

### `sort`

Runs the Go sorting program using the input file with the specified identifier.

```
make sort in=<filename without extension>
```

Example:

```
make sort in=1000
```

### `docker-build`

Builds a Docker image named "asd-1" for the project.

```
make docker-build
```

### `docker-run`

Runs the Docker container with the specified memory limit and input file.

```
make docker-run mem=<amount of memory> in=<filename without extension>
```

Example:

```
make docker-run mem=256m in=1000
```

### `pull`

Pulls the sorted output file from the Docker container to the host machine.

```
make pull
```

### `docker-clean`

Removes the Docker container named "asd-1".

```
make docker-clean
```

### `limited`

Runs a complete cycle of building the Docker image, running the container, pulling the output file, and cleaning the container.

```
make limited in=<filename without extension>
```

> Note: This target internally runs `docker-build`, `docker-run`, `pull`, and `docker-clean`.

Example

```
make limited in=1000
```

### `gen`

Runs a Go program to generate input files using the specified identifier. It creates file with specified number of elements under files/in/ in the root of the project

```
make gen in=<filename without extesnion>
```

Example

```
make gen in=1000
```
