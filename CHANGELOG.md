# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [unpublished]
- Add ability to specify key in failover password secret
- Improve leader election logic

## [1.0.0-rc1]

### Added
- A completely new version of operator does not have any compatibility with versions less than 1.0.0-rc1 
- Migration guide from 1.0.0-rc1
- API `tarantoo.io/v1beta1`
- Bump operator-sdk version (and other dependencies)
- HELM charts is now placed in standalone repository https://github.com/tarantool/helm-charts
- Update CRDs api version from `apiextensions.k8s.io/v1beta1` to `apiextensions.k8s.io/v1`
- Example cartridge is now placed in standalone repository
- Custom license
- Configurable image pull secrets
- Ability to configure cluster cookie
- Ability to configure application pods resources
- Ability to configure operator pods resources
- Release automation
- Ability to pass custom environment variables to application
- Update StatefulSet fields on the fly
- Configurable vshard groups
- Configurable vshard weights
- Configurable failover
- A lot of configurable params in new helm charts
- Publish kubernetes events with useful info 
- Topology leader switch
- Release automatically
- Support kubectl version 1.25+

### Removed
- API `tarantoo.io/v1alpha1` is not serve anymore

### Fixed
- Operator was not able to manage multiple cartridge clusters in multiple namespaces
- Remove all deprecated cartridge calls
- Reconciling every 5 second
- After restart minikube, tarantool cluster cannot be configured

## [0.0.9] - 2021-03-30

### Added
- Integration test for cluster_controller written with envtest and ginkgo
- Description of failover setting in the Cartridge Kubernetes guide
- Section to troubleshooting about CrashLoopBackOff
- Lua memory reserve for tarantool containers
- Guide to troubleshooting about replicas recreating

### Changed
- Requested verbs for a RBAC role Tarantool: remove all * verbs and resources

### Fixed
- Not working update of replicaset roles
- Not working update of container env vars
- Problem with a non-existent leader of cluster
- Flaky role_controller unit test

## [0.0.8] - 2020-12-16

### Added
- Support custom cluster domain name via variable `ClusterDomainName` in cartrige chart `values.yaml`
- New chart for deploying ready-to-use crud based application
- Ability to change TARANTOOL_WORKDIR in the Cartridge helm chart and the **default value is set to** `/var/lib/tarantool`