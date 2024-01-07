run: build
	cd cmd/ctfwriteups;  ./ctfwriteups 2>&1  | jq -Rrc '. as $$line | fromjson? // $$line'

build:
	cd cmd/ctfwriteups; go build

test:
	ENV=test go test ./...


debug:
	kubectl get pods
	kubectl logs ctfwriteup-77b669955f-bhjdv
	kubectl desribe ctfwriteup-77b669955f-bhjdv

# Prod
build-docker-prod:
	docker build -t ctfwriteups .

run-docker-prod:
	docker run -it --rm -p '8080:8080' ctfwriteup

push-docker-prod:
	doctl auth switch --context ctfwriteups
	doctl registry login
	docker tag ctfwriteups registry.digitalocean.com/ctfwriteups/ctfwriteups:latest
	docker push registry.digitalocean.com/ctfwriteups/ctfwriteups:latest

deploy: build-docker-prod push-docker-prod
	doctl auth switch --context ctfwriteups
	doctl auth init
	kubectl config use-context do-nyc1-k8s-ctfwriteups-org
	cd infra; kubectl apply -f ctfwriteups-prod.yml
	cd infra; kubectl rollout restart deployment/ctfwriteups-web-prod -n prod
