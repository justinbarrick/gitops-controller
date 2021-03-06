gitops-controller is a proof-of-concept GitOps controller that can sync from Git to
Kubernetes or from Kubernetes to Git.

*It is not yet ready for production use, but feedback is desired.*

# Overview

Users can choose on a per-resource basis whether or not a resource should be kept in
sync with its manifest in Git, whether the manifest in Git should be kept in sync with
the resource in Kubernetes, or whether it is ignored.

The gitops-controller clones a repository containing Kubernetes manifests and can
synchronize the repository in either direction. If a resource is configured to be
kept in sync with the Git repository, it will be updated any time it goes out of sync
with the Git repository. If a resource is configured to keep the Git repository in
sync, the Git repository will be updated any time it is out of sync with the
repository in Kubernetes.

# Rationale

Current GitOps tools typically only work in one way - resources are synchronized from
Git to Kubernetes. This works very well in most cases for ensuring that the cluster
state is enforced by Git.

However, there are some use-cases that require or are made easier by the ability to
synchronize from a cluster to the Git repository:

* Updating image versions in manifests when images are updated. Flux already
  implements this feature. The main limitation of the feature as implemented is that
  it is very specific to image updates and update logic is coupled to Flux.
* Composing GitOps operators with existing tools - for example, if Flux could write
  generic resource state back to a repository, it could be composed with Keel to
  provide image update functionality.
* If cluster state is written back to the Git repository, operations teams can easily
  make emergency changes in an auditable and repeatable way without needing to go
  through a full change cycle.
* Kubernetes RBAC can again be used to manage access to resources, some organizations
  may prefer this to managing access on a monorepo or simply to enable new workflows.
* Stateful resources that are created by automation in the cluster that should
  automatically be persisted to Git to ensure they are not tied to the life of the
  cluster. This can be useful for backing up PersistentVolumeClaims, VolumeSnapshots,
  and VolumeSnapshotContents.

This project was spawned by the last use-case, to enable fully automated backup
workflows for persistent volume claims in Kubernetes, see [backup-controller](https://github.com/justinbarrick/backup-controller)
for more details.

# Configuration

In the initial stages, the controller is configured with a rules configuration file
that tells the controller how to handle particular resources, but may eventually be
able to apply heuristics to changes to determine how they should be handled
automatically.

Configuration format:

* `branch`: the git branch to checkout.
* `gitUrl`: the URL of the Git repository to checkout.
* `gitPath`: a subdirectory in the Git repository to work in.
* `rules`: a list of `rule` objects to use when determining how
           changes should be handled.

`rule` objects:

* `apiGroups`: a list of API groups to match the rule on. If empty, the rule
               matches all API groups.
* `resources`: a list of resource types to match the rule on. If empty, the rule
               matches any resources.
* `labels`: a string label selector to match the rule on. If empty, the rule matches
            any labels.
* `filters`: a list of JSON path strings indicating which changes to include. For example, if filters is
             `["/metadata/annotations"]` then only changes to annotations will be matched.
* `syncTo`: the direction to synchronize matching resources - `kubernetes` to sync
            resources from Git to Kubernetes, `git` to sync resources from Kubernetes
            to Git.

# Running locally

To try out the gitops-controller, create a configuration file:

```
rules:
- apiGroups:
  - snapshot.storage.k8s.io
  resources: 
  - volumesnapshots
  - volumesnapshotcontents
  syncTo: kubernetes
```

Save it as `config.yaml` and then build and run the controller:

```
export GO111MODULE=on
go build
./gitops-controller git@github.com:justinbarrick/gitops-controller-test.git
```

Any `VolumeSnapshots` and `VolumeSnapshotContent` resources in the repository will
be created in your Kubernetes cluster.

# Deployment

To deploy the gitops-controller, create an SSH secret:

```
ssh-keygen -f gitops
kubectl create secret generic git-ssh -n gitops-controller --from-file=identity=$(pwd)/gitops -o yaml --dry-run >> deploy.yaml
```

Add the `gitops.pub` key to your repository's deploy keys with read and write
access.

Update the ConfigMap in deploy.yaml as required and then deploy it:

```
kubectl apply -f deploy.yaml
```

## Custom SSH known hosts

The image comes preloaded with known hosts for gitlab and github, however, you can add a new known_hosts file by mounting
a ConfigMap to `/ssh/known_hosts`.
