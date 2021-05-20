package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func pS2A(str pulumi.String) pulumi.StringArray {
	return pulumi.StringArray {
		str,
	}
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		args := &ec2.InstanceArgs{
			Ami: pulumi.String("ami-0d382e80be7ffdae5"), // ubuntu 20.04 x86-64
			InstanceType: pulumi.String(ec2.InstanceType_T2_Micro),
		}

		// Create an AWS resource (S3 Bucket)
		instance, err := ec2.NewInstance(ctx, "web-server-basic-hello", args)
		if err != nil {
			return err
		}

		// Export the name of the bucket
		ctx.Export("web-server-name", instance.ID())
		return nil
	})
}
