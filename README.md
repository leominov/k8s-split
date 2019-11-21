# k8s-split

[![Build Status](https://travis-ci.com/leominov/k8s-split.svg?branch=master)](https://travis-ci.com/leominov/k8s-split)
[![codecov](https://codecov.io/gh/leominov/k8s-split/branch/master/graph/badge.svg)](https://codecov.io/gh/leominov/k8s-split)

Split multi document Kubernetes specification file into separate files by `name` and `kind`.

## Usage

### File

```
$ k8s-split -f test_data/correct_multi.yaml -o ./
Found single.Pod
Saved to single.Pod.yaml
Found single.CronJob
Saved to single.CronJob.yaml
```

### Stdin

```
$ cat test_data/correct_multi.yaml | k8s-split -f - -o ./
Found single.Pod
Saved to single.Pod.yaml
Found single.CronJob
Saved to single.CronJob.yaml
```
