$(aws ecr get-login --no-include-email --region eu-west-1)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker build -t dev-trinity-container-registry .
docker tag dev-trinity-container-registry:latest 524454272832.dkr.ecr.eu-west-1.amazonaws.com/dev-trinity-container-registry:latest
docker push 524454272832.dkr.ecr.eu-west-1.amazonaws.com/dev-trinity-container-registry:latest

ssh ec2-user@34.248.140.154 -i ~/.ssh/dev-trinity-key-pair.pem 'docker kill $(docker ps --format "{{.ID}}" --filter "name=dev-trinit");uname -a;exit'

aws ecs update-service --cluster dev-trinity-cluster --service dev-trinity-container-service --force-new-deployment
