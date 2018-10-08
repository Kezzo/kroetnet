ssh ec2-user@34.254.60.219 -i ~/.ssh/dev-trinity-key-pair.pem 'docker restart $(docker ps --format "{{.ID}}" --filter "name=dev-trinit");uname -a;exit'
