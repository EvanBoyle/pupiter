module github.com/evanboyle/pupiter

go 1.14

require (
	github.com/dustinkirkland/golang-petname v0.0.0-20191129215211-8e5a1ed0cff0
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/sdk/v2 v2.6.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1
)

replace github.com/pulumi/pulumi/sdk/v2 => ../../pulumi/pulumi/sdk
