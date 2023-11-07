# simple-controller
This repository implements a simple controller based on controller-runtime for watching AppService resources as defined with a CustomResourceDefinition (CRD).

**Note:** When your clone is done, you should "go mod tidy" and then "go mod vendor".

This example mainly demonstrates how to use the controller-runtime to create an operator from zero, including the following:
1. How do initialize CRD-related files
2. How to use code-generator to generate code
3. How to implement controller business logic
4. How to debug/deploy locally on a cluster

We use `./hack/update-codegen.sh` to generate the deepcopy and register file.(`update-codegen.sh` comes from [kosmos-io/kosmos](https://github.com/kosmos-io/kosmos/tree/main)&[k8s.io/code-generator](https://github.com/kubernetes/code-generator))

The update-codegen script will automatically generate the following files:

+ pkg/apis/v1/zz_generated.deepcopy.go
+ pkg/apis/v1/zz_generated.register.go

> You should not edit these files, but run the script to generate your own files while writing your own controller.

