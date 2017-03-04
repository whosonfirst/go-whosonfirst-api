// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package cloudtrail

const (

	// ErrCodeARNInvalidException for service response error code
	// "ARNInvalidException".
	//
	// This exception is thrown when an operation is called with an invalid trail
	// ARN. The format of a trail ARN is:
	//
	// arn:aws:cloudtrail:us-east-1:123456789012:trail/MyTrail
	ErrCodeARNInvalidException = "ARNInvalidException"

	// ErrCodeCloudWatchLogsDeliveryUnavailableException for service response error code
	// "CloudWatchLogsDeliveryUnavailableException".
	//
	// Cannot set a CloudWatch Logs delivery for this region.
	ErrCodeCloudWatchLogsDeliveryUnavailableException = "CloudWatchLogsDeliveryUnavailableException"

	// ErrCodeInsufficientEncryptionPolicyException for service response error code
	// "InsufficientEncryptionPolicyException".
	//
	// This exception is thrown when the policy on the S3 bucket or KMS key is not
	// sufficient.
	ErrCodeInsufficientEncryptionPolicyException = "InsufficientEncryptionPolicyException"

	// ErrCodeInsufficientS3BucketPolicyException for service response error code
	// "InsufficientS3BucketPolicyException".
	//
	// This exception is thrown when the policy on the S3 bucket is not sufficient.
	ErrCodeInsufficientS3BucketPolicyException = "InsufficientS3BucketPolicyException"

	// ErrCodeInsufficientSnsTopicPolicyException for service response error code
	// "InsufficientSnsTopicPolicyException".
	//
	// This exception is thrown when the policy on the SNS topic is not sufficient.
	ErrCodeInsufficientSnsTopicPolicyException = "InsufficientSnsTopicPolicyException"

	// ErrCodeInvalidCloudWatchLogsLogGroupArnException for service response error code
	// "InvalidCloudWatchLogsLogGroupArnException".
	//
	// This exception is thrown when the provided CloudWatch log group is not valid.
	ErrCodeInvalidCloudWatchLogsLogGroupArnException = "InvalidCloudWatchLogsLogGroupArnException"

	// ErrCodeInvalidCloudWatchLogsRoleArnException for service response error code
	// "InvalidCloudWatchLogsRoleArnException".
	//
	// This exception is thrown when the provided role is not valid.
	ErrCodeInvalidCloudWatchLogsRoleArnException = "InvalidCloudWatchLogsRoleArnException"

	// ErrCodeInvalidEventSelectorsException for service response error code
	// "InvalidEventSelectorsException".
	//
	// This exception is thrown when the PutEventSelectors operation is called with
	// an invalid number of event selectors, data resources, or an invalid value
	// for a parameter:
	//
	//    * Specify a valid number of event selectors (1 to 5) for a trail.
	//
	//    * Specify a valid number of data resources (1 to 50) for an event selector.
	//
	//    * Specify a valid value for a parameter. For example, specifying the ReadWriteType
	//    parameter with a value of read-only is invalid.
	ErrCodeInvalidEventSelectorsException = "InvalidEventSelectorsException"

	// ErrCodeInvalidHomeRegionException for service response error code
	// "InvalidHomeRegionException".
	//
	// This exception is thrown when an operation is called on a trail from a region
	// other than the region in which the trail was created.
	ErrCodeInvalidHomeRegionException = "InvalidHomeRegionException"

	// ErrCodeInvalidKmsKeyIdException for service response error code
	// "InvalidKmsKeyIdException".
	//
	// This exception is thrown when the KMS key ARN is invalid.
	ErrCodeInvalidKmsKeyIdException = "InvalidKmsKeyIdException"

	// ErrCodeInvalidLookupAttributesException for service response error code
	// "InvalidLookupAttributesException".
	//
	// Occurs when an invalid lookup attribute is specified.
	ErrCodeInvalidLookupAttributesException = "InvalidLookupAttributesException"

	// ErrCodeInvalidMaxResultsException for service response error code
	// "InvalidMaxResultsException".
	//
	// This exception is thrown if the limit specified is invalid.
	ErrCodeInvalidMaxResultsException = "InvalidMaxResultsException"

	// ErrCodeInvalidNextTokenException for service response error code
	// "InvalidNextTokenException".
	//
	// Invalid token or token that was previously used in a request with different
	// parameters. This exception is thrown if the token is invalid.
	ErrCodeInvalidNextTokenException = "InvalidNextTokenException"

	// ErrCodeInvalidParameterCombinationException for service response error code
	// "InvalidParameterCombinationException".
	//
	// This exception is thrown when the combination of parameters provided is not
	// valid.
	ErrCodeInvalidParameterCombinationException = "InvalidParameterCombinationException"

	// ErrCodeInvalidS3BucketNameException for service response error code
	// "InvalidS3BucketNameException".
	//
	// This exception is thrown when the provided S3 bucket name is not valid.
	ErrCodeInvalidS3BucketNameException = "InvalidS3BucketNameException"

	// ErrCodeInvalidS3PrefixException for service response error code
	// "InvalidS3PrefixException".
	//
	// This exception is thrown when the provided S3 prefix is not valid.
	ErrCodeInvalidS3PrefixException = "InvalidS3PrefixException"

	// ErrCodeInvalidSnsTopicNameException for service response error code
	// "InvalidSnsTopicNameException".
	//
	// This exception is thrown when the provided SNS topic name is not valid.
	ErrCodeInvalidSnsTopicNameException = "InvalidSnsTopicNameException"

	// ErrCodeInvalidTagParameterException for service response error code
	// "InvalidTagParameterException".
	//
	// This exception is thrown when the key or value specified for the tag does
	// not match the regular expression ^([\\p{L}\\p{Z}\\p{N}_.:/=+\\-@]*)$.
	ErrCodeInvalidTagParameterException = "InvalidTagParameterException"

	// ErrCodeInvalidTimeRangeException for service response error code
	// "InvalidTimeRangeException".
	//
	// Occurs if the timestamp values are invalid. Either the start time occurs
	// after the end time or the time range is outside the range of possible values.
	ErrCodeInvalidTimeRangeException = "InvalidTimeRangeException"

	// ErrCodeInvalidTokenException for service response error code
	// "InvalidTokenException".
	//
	// Reserved for future use.
	ErrCodeInvalidTokenException = "InvalidTokenException"

	// ErrCodeInvalidTrailNameException for service response error code
	// "InvalidTrailNameException".
	//
	// This exception is thrown when the provided trail name is not valid. Trail
	// names must meet the following requirements:
	//
	//    * Contain only ASCII letters (a-z, A-Z), numbers (0-9), periods (.), underscores
	//    (_), or dashes (-)
	//
	//    * Start with a letter or number, and end with a letter or number
	//
	//    * Be between 3 and 128 characters
	//
	//    * Have no adjacent periods, underscores or dashes. Names like my-_namespace
	//    and my--namespace are invalid.
	//
	//    * Not be in IP address format (for example, 192.168.5.4)
	ErrCodeInvalidTrailNameException = "InvalidTrailNameException"

	// ErrCodeKmsException for service response error code
	// "KmsException".
	//
	// This exception is thrown when there is an issue with the specified KMS key
	// and the trail can’t be updated.
	ErrCodeKmsException = "KmsException"

	// ErrCodeKmsKeyDisabledException for service response error code
	// "KmsKeyDisabledException".
	//
	// This exception is deprecated.
	ErrCodeKmsKeyDisabledException = "KmsKeyDisabledException"

	// ErrCodeKmsKeyNotFoundException for service response error code
	// "KmsKeyNotFoundException".
	//
	// This exception is thrown when the KMS key does not exist, or when the S3
	// bucket and the KMS key are not in the same region.
	ErrCodeKmsKeyNotFoundException = "KmsKeyNotFoundException"

	// ErrCodeMaximumNumberOfTrailsExceededException for service response error code
	// "MaximumNumberOfTrailsExceededException".
	//
	// This exception is thrown when the maximum number of trails is reached.
	ErrCodeMaximumNumberOfTrailsExceededException = "MaximumNumberOfTrailsExceededException"

	// ErrCodeOperationNotPermittedException for service response error code
	// "OperationNotPermittedException".
	//
	// This exception is thrown when the requested operation is not permitted.
	ErrCodeOperationNotPermittedException = "OperationNotPermittedException"

	// ErrCodeResourceNotFoundException for service response error code
	// "ResourceNotFoundException".
	//
	// This exception is thrown when the specified resource is not found.
	ErrCodeResourceNotFoundException = "ResourceNotFoundException"

	// ErrCodeResourceTypeNotSupportedException for service response error code
	// "ResourceTypeNotSupportedException".
	//
	// This exception is thrown when the specified resource type is not supported
	// by CloudTrail.
	ErrCodeResourceTypeNotSupportedException = "ResourceTypeNotSupportedException"

	// ErrCodeS3BucketDoesNotExistException for service response error code
	// "S3BucketDoesNotExistException".
	//
	// This exception is thrown when the specified S3 bucket does not exist.
	ErrCodeS3BucketDoesNotExistException = "S3BucketDoesNotExistException"

	// ErrCodeTagsLimitExceededException for service response error code
	// "TagsLimitExceededException".
	//
	// The number of tags per trail has exceeded the permitted amount. Currently,
	// the limit is 50.
	ErrCodeTagsLimitExceededException = "TagsLimitExceededException"

	// ErrCodeTrailAlreadyExistsException for service response error code
	// "TrailAlreadyExistsException".
	//
	// This exception is thrown when the specified trail already exists.
	ErrCodeTrailAlreadyExistsException = "TrailAlreadyExistsException"

	// ErrCodeTrailNotFoundException for service response error code
	// "TrailNotFoundException".
	//
	// This exception is thrown when the trail with the given name is not found.
	ErrCodeTrailNotFoundException = "TrailNotFoundException"

	// ErrCodeTrailNotProvidedException for service response error code
	// "TrailNotProvidedException".
	//
	// This exception is deprecated.
	ErrCodeTrailNotProvidedException = "TrailNotProvidedException"

	// ErrCodeUnsupportedOperationException for service response error code
	// "UnsupportedOperationException".
	//
	// This exception is thrown when the requested operation is not supported.
	ErrCodeUnsupportedOperationException = "UnsupportedOperationException"
)
