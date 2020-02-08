# ko-operator

A Kubernetes operator to build and deploy go apps

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
  $ git clone https://github.com/feloy/ko-operator.git
  # Build and containerize the operator, then push the image
  $ make docker-build docker-push \
     IMG=eu.gcr.io/$PROJECT/ko-operator
  # Install the CRD
  $ make install
  customresourcedefinition.apiextensions.k8s.io/  kobuilders.ko.feloy.dev created
  # Deploy the operator on the cluster
  $ make deploy IMG=eu.gcr.io/$PROJECT/ko-operator
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

- Create a namespace on which you want to deploy your app:

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

- Create a Kubernetes service account with credentials to create the resources specified in your repository, for example:

  ```sh
  $ kubectl create sa ko-builder -n my-ns
  serviceaccount/ko-builder created

  $ kubectl create clusterrole ko-builder-role \
   --verb=list,get,create,patch \
   --resource=deployments.apps,services,namespaces,serviceaccounts

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
    serviceAccount:   ko-builder-sa@PROJECT.iam.gserviceaccount.com
    # repository containing sources to build
    repository: github.com/feloy/kopond
    # branch / tag / commit to checkout
    checkout: "1.0.0"
    # path containing manifests of resources to deploy
    configPath: /config
  ```

- Apply the template:

  ```sh
  $ kubectl apply -f sample.yaml
  kobuilder.ko.feloy.dev/kobuilder-sample created
  ```

- You can verify that the app has been deployed:

  > you can use [kubectl tree plugin](https://github.com/ahmetb/kubectl-tree) and see that the resources created are "owned" by the `KoBuilder` resource you created.

  ```sh
  $ kubectl tree kobuilders.ko.feloy.dev kobuilder-sample -n my-ns
  NAMESPACE  NAME                                           READY  REASON        AGE
  my-ns      KoBuilder/kobuilder-sample                     -                    45s
  my-ns      ├─ConfigMap/kobuilder-sample-configxzljd       -                    45s
  my-ns      └─Job/kobuilder-sample-job-k8djj               -                    45s
  my-ns        └─Pod/kobuilder-sample-job-k8djj-7lnr7       False  PodCompleted  45s
  my-ns          ├─Deployment/echo-controller               -                    12s
  my-ns          │ └─ReplicaSet/echo-controller-cd9fd5c75   -                    12s
  my-ns          │   └─Pod/echo-controller-cd9fd5c75-p7gwl  True                 12s
  my-ns          └─Service/echo-service                     -                    12s
  ```

- You can later update your deployment with a new release of your app by patching the `KoBuilder` resource:

  ```sh
  $ kubectl patch kobuilders.ko.feloy.dev \
     -n my-ns kobuilder-sample \
     -p '{"spec":{"checkout":"2.0.0"}}' \
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
