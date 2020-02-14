# k8s-split

[![Build Status](https://travis-ci.com/leominov/k8s-split.svg?branch=master)](https://travis-ci.com/leominov/k8s-split)
[![codecov](https://codecov.io/gh/leominov/k8s-split/branch/master/graph/badge.svg)](https://codecov.io/gh/leominov/k8s-split)

Split multi-document or `kind: List` Kubernetes specification file into separate files by `name` and `kind`.

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
