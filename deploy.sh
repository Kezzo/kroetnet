$(aws ecr get-login --no-include-email --region eu-west-1)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker build -t dev-trinity-match-server-registry .
docker tag dev-trinity-match-server-registry:latest 524454272832.dkr.ecr.eu-west-1.amazonaws.com/dev-trinity-match-server-registry:latest
docker push 524454272832.dkr.ecr.eu-west-1.amazonaws.com/dev-trinity-match-server-registry:latest

ssh ec2-user@34.254.60.219 -i ~/.ssh/dev-trinity-key-pair.pem 'docker kill $(docker ps --format "{{.ID}}" --filter "name=dev-trinit");uname -a;exit'

aws ecs update-service --cluster dev-trinity-cluster --service dev-trinity-container-service --force-new-deployment
