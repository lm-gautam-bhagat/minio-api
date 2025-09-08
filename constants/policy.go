package constants

var READ_ONLY_POLICY = "read-only-policy"

var read_only_policy_document = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"s3:ListBucket",
					"s3:GetObject"
				],
				"Resource": [
					"arn:aws:s3:::my-data-bucket/*"
				]
			}
		]
	}`

var PolicyMap = map[string]string{
	READ_ONLY_POLICY: read_only_policy_document,
}
