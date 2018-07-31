# Create SSM Params

### Create Parameters For Development
Start with the KMS key

``` bash
$ aws kms create-key --profile default --description "Pulpfree test key"

$ aws kms create-alias --alias-name alias/testPulpfree --target-key-id <key-id here>
# target-key-id is KeyId from above return

# $ aws ssm put-parameter --name /test/gales-dips2/S3Bucket \
#   --value ca-gales --type String
#   --key-id alias/testPulpfree --overwrite

$ aws ssm put-parameter --name /test/gales-dips2/DBName \
  --value gales-sales --type String --overwrite

# --key-id alias/testPulpfree --overwrite

# fetch params by path
$ aws ssm get-parameters-by-path --path /test/gales-dips2

```