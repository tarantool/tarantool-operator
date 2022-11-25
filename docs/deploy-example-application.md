# Deploy example application

This guide provides the instructions how to install the example cartridge application.
To deploy your own application you'd need to build your own docker image
via [Cartridge CLI](https://www.tarantool.io/en/doc/latest/book/cartridge/cartridge_cli/commands/pack/docker/).

### Install or update application using Helm

To install or upgrade the Cartridge App using Helm, run following commands from the shell:

- Add helm repository
  ```shell
  helm repo add tarantool https://tarantool.github.io/helm-charts/
  ```
  
- Configure application params
  ```shell
  helm show values tarantool/cartridge > values.yaml
  ```
  Edit `values.yaml` file to configure your deployment, the most useful variables described [here](#useful-helm-variables)

- Install/update app:
  - Replace the `my-namespace` with any name of namespace what you want.
  - Add the `--create-namespace` flag if your namespace is not created yet.
  - Execute following command:
    ```shell
    helm upgrade --install tarantool-app tarantool/cartridge --values ./values.yaml -n my-namespace [--create-namespace]
    ```

### Useful helm variables

| JSON Path                                                      | Type          | Description                                   |
|----------------------------------------------------------------|---------------|-----------------------------------------------|
| dockerconfigjson                                               | array         | Docker registry(ies) credentials              |
| storageClass                                                   | string        | An StorageClass name for requested disks      |
| tarantool.image.repository                                     | string        | Name of application docker image              |
| tarantool.image.tag                                            | string        | Tag of application docker image               |
| tarantool.bucketCount                                          | number        | Count of vshard buckets                       |
| tarantool.memtxMemory                                          | quantity      | Size of reserved memtx memory                 |
| tarantool.auth.user                                            | string        | Tarantool super-admin username                |
| tarantool.auth.password                                        | string        | Tarantool super-admin password/cluster cookie |
| tarantool.roles                                                | array         | Roles definition                              |
| tarantool.roles.*.replicasets                                  | number        | Count of replicasets in role                  |
| tarantool.roles.*.replicas                                     | number        | Count of replicas in each replicaset o role   |
| tarantool.roles.*.vshard.roles                                 | array<string> | List of vshard roles                          |
| tarantool.roles.*.persistence.spec.resources.requests.storage  | quantity      | Size of disk for each instance of role        |
