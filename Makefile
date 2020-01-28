include .env

default: build

deploy: build awsPackage awsDeploy

clean:
	@rm -rf dist
	@mkdir -p dist

build: clean
	@for dir in `ls handler`; do \
		GOOS=linux go build -o dist/$$dir github.com/pulpfree/gales-fuelsale-export/handler/$$dir; \
	done
	cp ./config/defaults.yaml dist/

validate:
	sam validate

run: build
	sam local start-api -n env.json

awsPackage:
	aws cloudformation package \
   --template-file template.yml \
   --output-template-file packaged-template.yml \
   --s3-bucket $(AWS_BUCKET_NAME) \
   --s3-prefix $(AWS_BUCKET_PREFIX) \
   --profile $(AWS_PROFILE)

awsDeploy:
	aws cloudformation deploy \
   --template-file packaged-template.yml \
   --stack-name $(AWS_STACK_NAME) \
   --capabilities CAPABILITY_IAM \
   --profile $(AWS_PROFILE)

describe:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(AWS_STACK_NAME)