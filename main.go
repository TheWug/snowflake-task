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
			UserData:            pulumi.String(string(script)),
			KeyName:             keypair.KeyName,
		})
		if err != nil {
			return err
		}

		// Export the name of the bucket
		ctx.Export("web-server-name", instance.ID())
		return nil
	})
}
