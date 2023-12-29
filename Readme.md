# vCenter filesd Exporter
This is a Go application used to discover Virtual Machine (VM) IPs in vCenter and return a JSON for Prometheus file service discovery.

# Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

# Prerequisites
- Go (version 1.16 or later)
- Docker (optional)
# Installation
1. Clone the repository:
```bash
git clone <repository_url>
```
2. Navigate to the project directory:
```
cd <project_directory>
```
3. Install the dependencies:
```
go mod download
go run .
```

# Usage
The application uses environment variables for configuration. The following environment variables are required:

- `VC_HOST`: The vCenter host.
- `VC_USER`: The vCenter username.
- `VC_PASSWORD`: The vCenter password.
- `PORT`: The port on which the application will run.
The application also accepts the following flags:

- `-tls`: Enable TLS. Default is `false`.

## Controllers
The application has a Prometheus controller defined in controller/prom.go. This controller handles the /prom endpoint and is responsible for returning the JSON for Prometheus file service discovery.

## vCenter Helper
The vCenter helper defined in vcenter-helper/vcenter-helper.go is responsible for discovering VM IPs in vCenter.

## Docker
A Dockerfile is provided if you wish to build a Docker image of the application.

# API reference

`localhost:<port>/<path>`

- `port`: This is set as an environment variable
- `path`: This is the path to vCenter folder that you want to scrape. Keep in mind that this discovers the folder recursively.