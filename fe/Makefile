run:
	npm start

deploy:
	ng build --configuration production --deploy-url https://assets.ctfwriteups.org/
	s3cmd put --recursive --acl-public ./dist/ctfwriteups/* s3://assets-ctfwriteups-org
	rm -rf ../be/cmd/ctfwriteups/dist
	cp -r dist/ctfwriteups ../be/cmd/ctfwriteups/dist