# k8s-split

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/leominov/k8s-split)](https://github.com/leominov/k8s-split/releases/latest)
[![Build Status](https://travis-ci.com/leominov/k8s-split.svg?branch=master)](https://travis-ci.com/leominov/k8s-split)
[![codecov](https://codecov.io/gh/leominov/k8s-split/branch/master/graph/badge.svg)](https://codecov.io/gh/leominov/k8s-split)

Split multi-document or `kind: List` Kubernetes specification file into separate files by `name` and `kind`. It is possible to save splitted files in directories based on longest name prefix or on the value of `app.kubernetes.io/part-of` tag.

## Usage

### File

```
$ k8s-split -f test_data/correct_multi.yaml -o ./
Found single.Pod
Saved to single.Pod.yaml
Found single.CronJob
Saved to single.CronJob.yaml
```

or

```
$ k8s-split -f test_data/correct_multi.yaml -o ./ --prefix
Found single.Pod
Saved to single/single.Pod.yaml
Found single.CronJob
Saved to single/single.CronJob.yaml
```

or

```
$ k8s-split -f test_data/correct_multi_prefix.yaml -o ./ --tag
Found application.Pod
Saved to bar/application.Pod.yaml
Found application.Service
Saved to bar/application.Service.yaml
Found application-backup.CronJob
Saved to foo/application-backup.CronJob.yaml
```

or

```
$ k8s-split -f test_data/correct_list.yaml -o ./
Found dco-manager-core-credentials.Secret
Saved to dco-manager-core-credentials.Secret.yaml
Found default-token-kzrjn.Secret
Saved to default-token-kzrjn.Secret.yaml
```

### Stdin

```
$ cat test_data/correct_multi.yaml | k8s-split -f -
Found single.Pod
Saved to single.Pod.yaml
Found single.CronJob
Saved to single.CronJob.yaml
```

or

```
$ kubectl get secrets -o yaml | k8s-split -f -
Found dco-manager-core-credentials.Secret
Saved to dco-manager-core-credentials.Secret.yaml
Found default-token-kzrjn.Secret
Saved to default-token-kzrjn.Secret.yaml
```

or

```
$ kustomize build test_data | k8s-split -f - --prefix
Found single.CronJob
Saved to single/single.CronJob.yaml
Found single.Pod
Saved to single/single.Pod.yaml
```

## Download

* https://github.com/leominov/k8s-split/releases/latest
