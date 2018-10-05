ssh ec2-user@34.248.140.154 -i ~/.ssh/dev-trinity-key-pair.pem 'docker restart $(docker ps --format "{{.ID}}" --filter "name=dev-trinit");uname -a;exit'
