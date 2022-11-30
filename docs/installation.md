# Install Tarantool Kubernetes Operator CE

## Table of Contents

- [Prerequisites](#prerequisites)
- [Install or Update the Operator using Helm](#install-or-update-the-operator-using-helm)

## Prerequisites

- A kubernetes cluster in version ≥1.16, with [dynamic volume provisioning](https://kubernetes.io/docs/concepts/storage/dynamic-provisioning/)
- [kubectl](https://kubernetes.io/ru/docs/tasks/tools/install-kubectl/)
- [helm](https://helm.sh/docs/intro/install/)

## Install or update the Operator using Helm

To install or upgrade the CRD's and the Tarantool Operator using Helm, run following commands from the shell:

- Add helm repository

  ```shell
  helm repo add tarantool https://tarantool.github.io/helm-charts/
  ```

- Configure the Operator deployment (optional)

  ```shell
  helm show values tarantool/tarantool-operator > operator-values.yaml
  ```

  Edit `operator-values.yaml` file if needed, all parameters described in comments.
- Install/Update operator (choose on of)
  - To `default` namespace:
  
    - Add the `--values ./operator-values.yaml` option if you are using custom configuration of deployment.
    - Execute following command:
      ```shell
      helm upgrade --install tarantool-operator-ce tarantool/tarantool-operator [--values ./operator-values.yaml]
      ```
      
  - To custom namespace:
    
    - Replace the `my-namespace` with any name of namespace what you want.
    - Add the `--create-namespace` flag if your namespace is not created yet.
    - Add the `--values ./operator-values.yaml` option if you are using custom configuration of deployment.
    - Execute following command:
      ```shell
      helm upgrade --install tarantool-operator-ce tarantool/tarantool-operator \
	  -n my-namespace [--create-namespace] [--values ./operator-values.yaml]
      ```
