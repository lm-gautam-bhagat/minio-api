package constants

type PolicyType string

const (
	READ_ONLY_POLICY PolicyType = "read-only-policy"
)

var read_only_policy_document = `{
		"Version": "2012-10-17",
		"Statement": [
			{
			"Effect": "Allow",
			"Action": [
				"s3:ListAllMyBuckets",
				"s3:ListBucket",
				"s3:GetObject"
			],
			"Resource": [
				"arn:aws:s3:::*",
				"arn:aws:s3:::*/*"
			]
			}
		]
		}
`

var PolicyMap = map[PolicyType]string{
	READ_ONLY_POLICY: read_only_policy_document,
}

func (p PolicyType) IsValid() bool {
	switch p {
	case READ_ONLY_POLICY:
		return true
	default:
		return false
	}
}
