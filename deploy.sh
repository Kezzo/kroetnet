$(aws ecr get-login --no-include-email --region eu-west-1)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker build -t dev-trinity-container-registry .
docker tag dev-trinity-container-registry:latest 524454272832.dkr.ecr.eu-west-1.amazonaws.com/dev-trinity-container-registry:latest
docker push 524454272832.dkr.ecr.eu-west-1.amazonaws.com/dev-trinity-container-registry:latest
aws ecs update-service --cluster dev-trinity-cluster --service dev-trinity-container-service --force-new-deployment
