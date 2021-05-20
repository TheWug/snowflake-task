package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"io/ioutil"
)

func pS2A(str pulumi.String) pulumi.StringArray {
	return pulumi.StringArray {
		str,
	}
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		/*
		required configuration options:
		deploy-script (you can set it to deploy_webserver.sh)
		public-key-file (you can set it to your favorite RSA public key)
		*/

		conf := config.New(ctx, "")
		script, err := ioutil.ReadFile(conf.Require("deploy-script"))
		if err != nil {
			return err
		}

		pubkey, err := ioutil.ReadFile(conf.Require("public-key-file"))
		if err != nil {
			return err
		}

		// create a new security group, which permits inbound on tcp ports 22 and 80, and any outbound
		sec_group, err := ec2.NewSecurityGroup(ctx, "public-http-server", &ec2.SecurityGroupArgs{
			Name:        pulumi.String("basic-http-ssh"),
			Description: pulumi.String("Allow incoming HTTP and SSH traffic, and any outgoing traffic."),
			Ingress:     ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					CidrBlocks:  pS2A(pulumi.String("0.0.0.0/0")),
					FromPort:    pulumi.Int(80),
					ToPort:      pulumi.Int(80),
					Protocol:    pulumi.String("tcp"),
					Description: pulumi.String("Allow inbound simple HTTP"),
				},
				ec2.SecurityGroupIngressArgs{
					CidrBlocks:  pS2A(pulumi.String("0.0.0.0/0")),
					FromPort:    pulumi.Int(22),
					ToPort:      pulumi.Int(22),
					Protocol:    pulumi.String("tcp"),
					Description: pulumi.String("Allow inbound SSH"),
				},
			},
			Egress:      ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					CidrBlocks:  pS2A(pulumi.String("0.0.0.0/0")),
					FromPort:    pulumi.Int((0)),
					ToPort:      pulumi.Int(0),
					Protocol:    pulumi.String("-1"),
					Description: pulumi.String("Allow all outbound traffic"),
				},
			},
		})
		if err != nil {
			return err
		}

		// create a keypair from your provided public key
		keypair, err := ec2.NewKeyPair(ctx, "instance-sshkey", &ec2.KeyPairArgs{
			PublicKey:     pulumi.String(string(pubkey)),
			KeyNamePrefix: pulumi.String("instance-sshkey"),
		})
		if err != nil {
			return err
		}

		// create an ec2 instance with the provided public key and feed it the deployment script
		instance, err := ec2.NewInstance(ctx, "web-server-basic-demo", &ec2.InstanceArgs{
			Ami:                 pulumi.String("ami-0d382e80be7ffdae5"), // ubuntu 20.04 x86-64
			InstanceType:        pulumi.String(ec2.InstanceType_T2_Micro),
			VpcSecurityGroupIds: pulumi.StringArray{
				sec_group.ID(),
			},
			UserData:            pulumi.String(string(script)),
			KeyName:             keypair.KeyName,
		})
		if err != nil {
			return err
		}

		ctx.Export("web-server-id", instance.ID())
		ctx.Export("public-key-id", keypair.ID())
		ctx.Export("public-IP",     instance.PublicIp)
		ctx.Export("public-DNS",    instance.PublicDns)
		return nil
	})
}
