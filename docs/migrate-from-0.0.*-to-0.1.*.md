# Migration guide from version 0.0.* to version ≥0.1.0

## Terms

- `legacy operator` - Tarantool operator in version <0.1.0
- `new operator` - Tarantool operator in version ≥0.1.0

## Important

1. **Read whole guide before start**
2. Be careful with your actions.
3. There is no simple way to rollback
4. Make a whole backup and be ready to restore from it.
5. All cartridge application controlled by legacy operator in your kubernetes cluster MUST be migrated at the same time.
6. It can be helpful if you can change connection address in all application which communicating with your cartridge app.
7. Make sure you have enough resources in your kubernetes cluster 

## Migration process

1. Find out a namespace of your application
2. Find out a most stable instance of your application in kubernetes pods and remember pod name

   ```shell
   kubectl -n tarantool-app get pod
   NAME          READY   STATUS    RESTARTS   AGE
   routers-0-0   1/1     Running   0          44h
   storage-0-0   1/1     Running   0          44h
   storage-0-1   1/1     Running   0          44h
   ```
   
   It seems good to choose vshard-router or any other instance without data.
   In this migration guide we will use `routers-0-0` pod.
3. Configure new helm chart for your application using pod name from previous step

   ```yaml
   tarantool:
     foreignLeader: "routers-0-0"  # The name of pod from step 2 - most important field
     bucketCount: 30000 # doest meter in migration flow
     auth:
       password: "secret-cluster-cookie" # you should use your actual cluster cookie here 
     image:
       repository: tarantool/tarantool-operator-examples-kv
       tag: 0.0.4
       pullPolicy: IfNotPresent
   roles: [...] # Your application roles here 
   ```
4. Add helm repository
   
   ```shell
   helm repo add tarantool https://tarantool.github.io/helm-charts/
   ```   

5. Make sure you are ready to update, at next step you going to lose a way to rollback

6. Install new operator using official helm-chart of new operator
   You can follow [installation guide](./installation.md) at this step

   ```shell
   helm upgrade --install tarantool-operator-ce tarantool/tarantool-operator 
   ```
   
   Custom resource definitions will not be updated by helm, you MUST update it manually:

   ```shell
   kubectl apply -f https://raw.githubusercontent.com/tarantool/helm-charts/master/charts/tarantool-operator/crds/apiextensions.k8s.io_v1_customresourcedefinition_cartridgeconfigs.tarantool.io.yaml
   kubectl apply -f https://raw.githubusercontent.com/tarantool/helm-charts/master/charts/tarantool-operator/crds/apiextensions.k8s.io_v1_customresourcedefinition_clusters.tarantool.io.yaml
   kubectl apply -f https://raw.githubusercontent.com/tarantool/helm-charts/master/charts/tarantool-operator/crds/apiextensions.k8s.io_v1_customresourcedefinition_roles.tarantool.io.yaml
   ```
   
7. Manually delete legacy operator (**DO NOT USE HELM**)

   ```shell
   kubectl delete ns tarantool-operator
   ```

8. Wait for new operator to be ready
9. Install new cartridge helm chart using values file from step 3

   - **Make sure that new helm release have DIFFERENT NAME from old app**
     Replace `you-app` with name which you wish in following command.
   - **Make sure that new helm release have SAME NAMESPACE as old app**
     Replace `your-namespace` with name of namespace where your app was deployed.
   - new cartridge app will be installed into you kubernetes cluster  
   - replicasets of new app will join existing cartridge app immediately
   - vshard will start buckets re-balancing immediately 
   - it can produce lots of internal network traffic 

   ```shell
   helm upgrade --install -n your-namespace you-app tarantool/cartridge --values ./operator-values.yaml
   ```

10. Wait for all instances of deployed cartridge app to be ready
11. At this step you need to **connect to Cartridge UI on any new instance** using port-forwarding or any other method 
12. At this step in Cartridge UI you can see topology similar to following picture:
    <img src="./assets/migration.png"> 
    
    As you can see there are new and old instances joins one cartridge cluster

13. Set zero weight on all old replicasets with vshard-storage role
14. Wait until all buckets goes to new replicasets
15. Disable and expel all old instance one by one
16. Manually delete all of old resource (DO NOT USE HELM)
    
    **Be careful with names, in this guide we are using names from legacy docs**

    ```shell
    kubectl -n tarantool-app delete clusters.tarantool.io -l tarantool-cluster
    kubectl -n tarantool-app delete roles.tarantool.io -l tarantool.io/cluster-id=tarantool-cluster
    kubectl -n tarantool-app delete sts -l tarantool.io/cluster-id=tarantool-cluster
    kubectl -n tarantool-app delete pod -l tarantool.io/cluster-id=tarantool-cluster
    kubectl -n tarantool-app delete svc tarantool-cluster
    kubectl -n tarantool-app delete svc routers
    kubectl -n tarantool-app delete svc storage
    kubectl -n tarantool-app delete secret sh.helm.release.v1.cartridge-app.v1
    ```

17. Repeat step 1-3 and 8-15 for each cartridge app in your kubernetes cluster    
18. Delete unused CRD

    ```shell
    kubectl delete crd replicasettemplates.tarantool.io
    ```

19. Done, now you are using new operator!
