## TODO ##

* Github actions for CI / CD
* CD to dev environment
* Metrics with Prometheus + Grafana
* Extract golang microservice bootstrap code

## DONE ##

* Logging
* CLI flags
* Config file
* Error handling
* SQL repository bootstrap
* SQL migrations
* HTTP API bootstrap
* HTTP API docs
* Docker
* Makefile
* Unit testing

## DEV WORKFLOW ##

* Locally, run unit tests for specific package
* Push to every branch -> run Github actions for CI and CD to dev environment : lint, tests, docker
* Push to master branch -> deploy to prod env
