# ko-operator

A Kubernetes operator to build and deploy go apps using [ko](https://github.com/google/ko).

`ko-operator` is a Kubernetes operator using the [ko-builder](https://github.com/feloy/ko-builder) image to build and deploy go apps on Kubernetes with a minimum of configuration.

`ko-operator` is built with [kubebuilder](https://book.kubebuilder.io/).

## On Google Cloud GKE

### Prepare the cluster

- define the `PROJECT` variable with the name of your Google Cloud project:

  ```sh
  $ PROJECT=my-project
  ```

- Deploy the operator:

  ```sh
  $ kubectl apply -f https://raw.githubusercontent.com/feloy/ko-operator/master/dist/ko-operator.yaml
  ```

- Create a Google Cloud service account with access to the registry of the project and get a JSON key for this service account, for example:

  ```sh
  $ gcloud iam service-accounts create ko-builder-sa \
     --description "SA for KO builder" \
     --display-name "SA for KO builder"
  Created service account [ko-builder-sa].

  $ gcloud iam service-accounts enable \
     ko-builder-sa@$PROJECT.iam.gserviceaccount.com
  Enabled service account [ko-builder-sa@$PROJECT.iam.gserviceaccount.com].

  $ gcloud iam service-accounts keys create key.json \
     --iam-account \
     ko-builder-sa@$PROJECT.iam.gserviceaccount.com
  created key [...] of type [json] as [key.json] for [ko-builder-sa@ko-demo.iam.gserviceaccount.com]

  $ gcloud projects add-iam-policy-binding $PROJECT \
    --member \
    "serviceAccount:ko-builder-sa@$PROJECT.iam.gserviceaccount.com" --role "roles/storage.admin"
  Updated IAM policy for project [$PROJECT].
  ```

### Prepare namespaces

For each namespace you want to deploy apps using the ko-opertor:

- Create the namespace:

  ```sh
  $ kubectl create namespace my-ns
  namespace/my-ns created
  ```

- Create a Kubernetes secret named `gcloud` with the key.json contents:

  ```sh
  $ kubectl create secret generic gcloud \
     -n my-ns \
     --from-file=key.json
  secret/gcloud created
  ```

- Create a Kubernetes service account named `ko-builder` with credentials to create the resources specified in your repository, for example:

  ```sh
  # Create a service account in the namespace
  $ kubectl create sa ko-builder -n my-ns
  serviceaccount/ko-builder created

  # Create the role with credentials necessary to work with resources
  # This resource is non-namespaced, create it only once per cluster
  $ kubectl create clusterrole ko-builder-role \
   --verb=list,get,create,patch \
   --resource=deployments.apps,services,namespaces,serviceaccounts

  # Bind the role to the service account on this namespace
  $ kubectl create clusterrolebinding ko-builder-rolebinding-my-ns \
     --clusterrole=ko-builder-role \
     --serviceaccount=my-ns:ko-builder
  clusterrolebinding.rbac.authorization.k8s.io/ko-builder-rolebinding-my-ns created
  ```

### For each program you want to build and deploy

- Create a `KoBuilder` custom resource template. Adapt the fields with your own values:

  ```yaml
  # sample.yaml
  apiVersion: ko.feloy.dev/v1alpha1
  kind: KoBuilder
  metadata:
    name: kobuilder-sample
    namespace: my-ns
  spec:
    # the registry on which to push built images
    registry: eu.gcr.io/PROJECT
    # a Google Cloud service account with access to registry
    serviceAccount: ko-builder-sa@PROJECT.iam.gserviceaccount.com
    # repository containing sources to build
    repository: github.com/feloy/kopond
    # branch / tag / commit to checkout
    checkout: "2.1.0"
    # path containing manifests of resources to deploy
    configPath: /config
  ```

- Apply the template:

  ```sh
  $ kubectl apply -f sample.yaml
  kobuilder.ko.feloy.dev/kobuilder-sample created
  ```

> You can verify that the app is deploying or has been deployed by using the [kubectl tree plugin](https://github.com/ahmetb/kubectl-tree) and see that the resources created are "owned" by the `KoBuilder` resource you created.

- The app is deploying (you see a configmap containing configuration and a `KoBuilder` job deploying the app):

  ```sh
  $ kubectl tree kobuilders.ko.feloy.dev kobuilder-sample -n my-ns
  NAMESPACE  NAME                                 READY  REASON  AGE
  my-ns      KoBuilder/kobuilder-sample           -              11s
  my-ns      ├─ConfigMap/kobuilder-sample-config  -              11s
  my-ns      └─Job/kobuilder-sample-job           -              11s
  my-ns        └─Pod/kobuilder-sample-job-mcfkm   True           11s
  ```

- The app has been deployed (the Job finished and has been deleted, the configmap is still here for information and the resources from the repository has been deployed):

  ```sh
  $ kubectl tree kobuilders.ko.feloy.dev kobuilder-sample -n my-ns
  NAMESPACE  NAME                                        READY  REASON  AGE 
  my-ns      KoBuilder/kobuilder-sample                  -              2m3s
  my-ns      ├─ConfigMap/kobuilder-sample-config         -              2m3s
  my-ns      ├─Deployment/echo-controller                -              92s 
  my-ns      │ └─ReplicaSet/echo-controller-777bc46cf8   -              92s 
  my-ns      │   └─Pod/echo-controller-777bc46cf8-fjhhm  True           92s 
  my-ns      └─Service/echo-service                      -              92s 
  ```

- You can also use `kubectl get kobuilders`:

  ```sh
  $ kubectl get kobuilders -w -n my-ns
  NAME               REPOSITORY                CHECKOUT   STATE
  kobuilder-sample   github.com/feloy/kopond   2.1.0      Deploying
  kobuilder-sample   github.com/feloy/kopond   2.1.0      Deployed
  ```

- You can later update your deployment with a new release of your app by patching the `KoBuilder` resource:

  ```sh
  $ kubectl patch kobuilders.ko.feloy.dev \
     -n my-ns kobuilder-sample \
     -p '{"spec":{"checkout":"2.2.0"}}' \
     --type=merge
  kobuilder.ko.feloy.dev/kobuilder-sample patched
  ```

- Thanks to these owner references, the created objects will be deleted when you delete the `KoBuilder` resource:

  ```sh
  $ kubectl delete kobuilders kobuilder-sample -n my-ns
  kobuilder.ko.feloy.dev "kobuilder-sample" deleted
  # After several seconds
  $ kubectl get deployments -n my-ns
  No resources found in demoop namespace.
  $ kubectl get svc -n my-ns
  No resources found in demoop namespace.
  ```
