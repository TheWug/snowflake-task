package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		vpc_args := &ec2.VpcArgs{
			CidrBlock:	pulumi.String("172.20.0.0/16"),
		}
		test_vpc, err := ec2.NewVpc(ctx, "test-vpc", vpc_args)
		if err != nil {
			return err
		}

		subnet_args := &ec2.SubnetArgs{
			VpcId:		test_vpc.ID(),
			CidrBlock:	pulumi.String("172.20.20.0/24"),
		}
		test_subnet, err := ec2.NewSubnet(ctx, "test-subnet", subnet_args)
		if err != nil {
			return err
		}

		iface_args := &ec2.NetworkInterfaceArgs{
			SubnetId:	test_subnet.ID(),
			PrivateIps:	pulumi.StringArray{
				pulumi.String("172.20.20.20"),
			}
		}
		iface, err := ec2.NewNetworkInterface(ctx, "test-interface", iface_args)
		if err != nil {
			return err
		}

		args := &ec2.InstanceArgs{
			Ami: pulumi.String("ami-0d382e80be7ffdae5"), // ubuntu 20.04 x86-64
			InstanceType: pulumi.String(ec2.InstanceType_T2_Micro),
			NetworkInterfaces: ec2.InstanceNetworkInterfaceArray{
				&ec2.InstanceNetworkInterfaceArgs{
					NetworkInterfaceId: iface.ID(),
					DeviceIndex: pulumi.Int(0),
				},
			},
		}
		ec2_instance_type := "t2.micro"

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
