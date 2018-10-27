$(aws ecr get-login --no-include-email --region eu-west-1)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker build -t dev-trinity-match-server-registry .
docker tag dev-trinity-match-server-registry:latest 524454272832.dkr.ecr.eu-west-1.amazonaws.com/dev-trinity-match-server-registry:latest
docker push 524454272832.dkr.ecr.eu-west-1.amazonaws.com/dev-trinity-match-server-registry:latest

aws ecs update-service --cluster dev-trinity-match-server-cluster --service dev-trinity-match-server-service --force-new-deployment
