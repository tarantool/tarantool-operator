<a href="http://tarantool.org">
   <img src="https://static.tarantool.io/pub/221123-0838-43389b9/tarantool/images/current-logo.svg" align="right">
</a>

# Tarantool Community Kubernetes Operator

[![Tests][gh-test-actions-badge]][gh-actions-url]
[![Lint][gh-lint-actions-badge]][gh-actions-url]

This is a [Kubernetes Operator](https://coreos.com/operators/) which deploys [Tarantool Cartridge](https://github.com/tarantool/cartridge)-based
cluster on Kubernetes.

If you are a Tarantool Enterprise customer, or need Enterprise features such as rolling update, scaling down and may others
you can use the [Tarantool Operator Enterprise](https://www.tarantool.io/ru/kubernetesoperator).

## IMPORTANT NOTICE

Begins from v1.0.0-rc1 Tarantool Community Kubernetes Operator was completely rewrote.

API version was bumped and any backward compatibility was dropped.

There is only one approved method to migrate from version 0.0.0 to versions >=1.0.0-rc1, 
please follow [migration guide](./docs/migrate-from-0.0.*-to-1.0.0.md).

## Table of contents

* [Getting started](#getting-started)
* [Documentation](#documentation)
* [Contribute](#contribute)

## Getting started

- [Install the Operator](./docs/installation.md)
- [Deploy example application](./docs/deploy-example-application.md)

## Documentation

The documentation is work in progress...

At the moment you can use official [helm-chart](https://github.com/tarantool/helm-charts/tree/master/charts/tarantool-operator) 
and receive useful information from comments in default [values.yaml](https://github.com/tarantool/helm-charts/blob/master/charts/tarantool-operator/values.yaml) file 

## Contribute

Please follow the [development guide](./docs/development-guide.md)

[gh-lint-actions-badge]: https://github.com/tarantool/tarantool-operator/actions/workflows/lint.yml/badge.svg
[gh-test-actions-badge]: https://github.com/tarantool/tarantool-operator/actions/workflows/test.yml/badge.svg
[gh-actions-url]: https://github.com/tarantool/tarantool-operator/actions
