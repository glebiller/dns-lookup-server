# Interview challenge

## Implementation

I used [goswagger](https://goswagger.io/) to autogenerate an HTTP server that comply with the given the swagger specification.
This allows me to primarily focus on implementation & CI/CD tasks.

Endpoints are registered in the `restapi.configure_dns_lookup.go` file, and implementations are located in the `dnslookup` package.
I also modified the swagger definition to add an error 500 to the history endpoint. 
Most of the errors returned are server-side and a 5xx error fit better than a 4xx. 

In order to persist the history of calls, I decided to use InfluxDB as a fast time-series database.
As per the requirements it will be a perfect fit for logging the successful queries and retrieving the history quickly.

It's also a bit more uncommon and fun than a traditional SQL database for this project, it was allowing me to not bother with managing table and index. 
however in a more complete product with different features planned it would probably a better choice to pick a standard relational database ;)

With more time, I would also have improved logging using the Log service from the generated server.
Error handling could be improved by having more explicit messages.

## Buildings

The release process is using SemVer for versioning, and use Release Please to automatically create new releases.
Release Please also take care of updating the CHANGELOG, as long as all commits are following the Conventional Commits pattern.
The Commitsar tool is used to validate that PRs are following this convention.

Build workflows using Github Actions. 
All PRs are validated using Review Dog to provide fast feedback as comments inside the PRs.
Once merged, the Release workflow trigger a GoReleaser build that create both binaries for multiple platform
along with multi-arch docker images.

## Deployment

#### Docker-compose

The docker-compose allow a simple local dev environment to be run locally.
Build and run with `docker-compose up -d --build`.

#### Kubernetes

Kubernetes manifests are located in `deploy` folder. It's using Kustomize as primary tool. The resources located in this folder
would serve as a "base". In case we need to deploy to multiple environments (staging, production for instance), the use of
Kustomize Overlays would allow us to only change the attributes that need to be updated (i.e. replicas, env variables).

Additionally, as part of the CI/CD the K8s manifests should be tested for good practices and security. 
I picked "yamllint" and "kube-score" for code validation (run with make). 

As I time-boxed myself, I did not have the time to lunch & verify that the manifests are working.
It's also missing a readiness probe (currently I use liveness with the health endpoint), 
but for that app I cannot find any use case where the application would be unavailable temporarily.
