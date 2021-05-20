#!/bin/bash

# write a log in the web folder so we can view it from a browser
mkdir -p /var/www/html
exec 1>> /var/www/html/deploy_webserver.log 2>&1
set -x

# I said I'd do it with nginx and im sticking to it
sudo apt update
sudo apt install -y nginx

# write new index page
sudo tee /var/www/html/snowflake.html <<< '
<html>
  <head>
    <title>Hello Terry</title>
  </head>
  <body>
    <h1>This is a snowflake demo.</h1>
    <p><a href="https://github.com/TheWug/snowflake-task">Find it on github here.</a></p>
    <h3>A brief overview of how it works</h3>
    <p>1. using pulumi config, set deploy-script = <tt>deploy_webserver.sh</tt> and public-key-file = your preferred RSA public key (e.g. <tt>~/.ssh/id_rsa.pub</tt>)</p>
    <p>2. <tt>pulumi up</tt> will provision an instance with SSH and HTTP inbound and open outbound, deploy your public key to it, and run the deploy script.</p>
    <p>3. the provided deploy script will wait until the network is available, then install nginx, supply a custom landing page, update the config to point to it, and reload nginx. </p>
    <p><a href="deploy_webserver.log">Here is a link to the log file.</a></p>
  </body>
</html>'

# update nginx config to use snowflake.html as the default index page
sudo sed 's/\tindex .*/\tindex snowflake.html index.html/' -i /etc/nginx/sites-enabled/default

# instruct nginx to re-read config file
sudo systemctl reload nginx
