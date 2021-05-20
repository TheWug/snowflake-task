# snowflake-task

## Written for:
* Pulumi 3.2.1
* go 1.16.4
* aws-cli 2.2.5

## How to use:
1. `git clone`
1. make sure awscli and pulumi have API keys set up
1. `pulumi config set public-key-file [some RSA public key]`
1. `pulumi config set deploy-script deploy_webserver.sh`
1. `pulumi up`
1. wait for it to do its thing
1. click the link. The link will be posted before nginx is available. If you catch it too early and it shows the default page, or if it doesn't load at all, wait a few seconds and refresh it.
