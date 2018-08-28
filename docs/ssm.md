# Create SSM Params

### Create Parameters For Development
``` bash
# Start with the KMS key if it does not already exists
$ aws kms create-key --profile default --description "Pulpfree test key"

$ aws kms create-alias --alias-name alias/testPulpfree --target-key-id <key-id here>
# target-key-id is KeyId from above return

# encrypted parameter
# $ aws ssm put-parameter --name /test/gales-dips2/S3Bucket \
#   --value ca-gales --type String
#   --key-id alias/testPulpfree --overwrite

$ aws ssm put-parameter --name /test/gales-dips2/MongoDBName \
  --value gales-sales --type String --overwrite

$ aws ssm put-parameter --name /test/gales-dips2/CognitoClientID \
  --value us-east-1_gsB59wfzW --type String --overwrite

$ aws ssm put-parameter --name /test/gales-dips2/CognitoPoolID \
  --value 2084ukslsc831pt202t2dudt7c --type String --overwrite

$ aws ssm put-parameter --name /test/gales-dips2/CognitoRegion \
  --value us-east-1 --type String --overwrite

# delete a parameter
$ aws ssm delete-parameter --name /test/gales-dips2/DBName

# fetch params by path
$ aws ssm get-parameters-by-path --path /test/gales-dips2

```

### Production Parameters
``` bash
# Start with the KMS key if it does not already exists
$ aws kms create-key --profile default --description "Gales Production key"

$ aws kms create-alias --alias-name alias/GalesProd --target-key-id <key-id here>

# add tag
$ aws kms tag-resource --key-id <key-id-here> --tags TagKey=BillTo,TagValue=Gales

# create parameters
$ aws ssm put-parameter --name /prod/gales-dips2/MongoDBName \
  --value gales-sales --type String

$ aws ssm put-parameter --name /prod/gales-dips2/MongoDBHost \
  --value gales.mongo --type String

$ aws ssm put-parameter --name /prod/gales-dips2/MongoDBPassword \
  --value "<secret-password>" --type SecureString \
  --key-id alias/GalesProd --overwrite

$ aws ssm put-parameter --name /prod/gales-dips2/MongoDBUser \
  --value "<secret-username>" --type SecureString \
  --key-id alias/GalesProd --overwrite

$ aws ssm put-parameter --name /prod/gales-dips2/CognitoClientID \
  --value us-east-1_gsB59wfzW --type String --overwrite

$ aws ssm put-parameter --name /prod/gales-dips2/CognitoPoolID \
  --value 2084ukslsc831pt202t2dudt7c --type String --overwrite

$ aws ssm put-parameter --name /prod/gales-dips2/CognitoRegion \
  --value us-east-1 --type String --overwrite

# get parameter single
$ aws ssm get-parameters --names /prod/gales-dips2/MongoDBPassword \
  --with-decryption

# get parameters by path
$ aws ssm get-parameters-by-path --path /prod/gales-dips2
$ aws ssm get-parameters-by-path --path /prod/gales-dips2 --with-decryption

```