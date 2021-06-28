package aws

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/awsutil"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/errwrap"
)

const secretAccessKeyType = "access_keys"

func secretAccessKeys(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: secretAccessKeyType,
		Fields: map[string]*framework.FieldSchema{
			"access_key": {
				Type:        framework.TypeString,
				Description: "Access Key",
			},

			"secret_key": {
				Type:        framework.TypeString,
				Description: "Secret Key",
			},
			"security_token": {
				Type:        framework.TypeString,
				Description: "Security Token",
			},
		},

		Renew:  b.secretAccessKeysRenew,
		Revoke: b.secretAccessKeysRevoke,
	}
}

func genUsername(displayName, policyName, userType string) (ret string, warning string) {
	var midString string

	switch userType {
	case "iam_user":
		// IAM users are capped at 64 chars; this leaves, after the beginning and
		// end added below, 42 chars to play with.
		midString = fmt.Sprintf("%s-%s-",
			normalizeDisplayName(displayName),
			normalizeDisplayName(policyName))
		if len(midString) > 42 {
			midString = midString[0:42]
			warning = "the calling token display name/IAM policy name were truncated to fit into IAM username length limits"
		}
	case "sts":
		// Capped at 32 chars, which leaves only a couple of characters to play
		// with, so don't insert display name or policy name at all
	}

	ret = fmt.Sprintf("vault-%s%d-%d", midString, time.Now().Unix(), rand.Int31n(10000))
	return
}

func (b *backend) getFederationToken(ctx context.Context, s logical.Storage,
	displayName, policyName, policy string, policyARNs []string,
	iamGroups []string, lifeTimeInSeconds int64) (*logical.Response, error) {

	groupPolicies, groupPolicyARNs, err := b.getGroupPolicies(ctx, s, iamGroups)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if groupPolicies != nil {
		groupPolicies = append(groupPolicies, policy)
		policy, err = combinePolicyDocuments(groupPolicies...)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}
	if len(groupPolicyARNs) > 0 {
		policyARNs = append(policyARNs, groupPolicyARNs...)
	}

	stsClient, err := b.clientSTS(ctx, s)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	username, usernameWarning := genUsername(displayName, policyName, "sts")

	getTokenInput := &sts.GetFederationTokenInput{
		Name:            aws.String(username),
		DurationSeconds: &lifeTimeInSeconds,
	}
	if len(policy) > 0 {
		getTokenInput.Policy = aws.String(policy)
	}
	if len(policyARNs) > 0 {
		getTokenInput.PolicyArns = convertPolicyARNs(policyARNs)
	}

	// If neither a policy document nor policy ARNs are specified, then GetFederationToken will
	// return credentials equivalent to that of the Vault server itself. We probably don't want
	// that by default; the behavior can be explicitly opted in to by associating the Vault role
	// with a policy ARN or document that allows the appropriate permissions.
	if policy == "" && len(policyARNs) == 0 {
		return logical.ErrorResponse("must specify at least one of policy_arns or policy_document with %s credential_type", federationTokenCred), nil
	}

	tokenResp, err := stsClient.GetFederationToken(getTokenInput)
	if err != nil {
		return logical.ErrorResponse("Error generating STS keys: %s", err), awsutil.CheckAWSError(err)
	}

	resp := b.Secret(secretAccessKeyType).Response(map[string]interface{}{
		"access_key":     *tokenResp.Credentials.AccessKeyId,
		"secret_key":     *tokenResp.Credentials.SecretAccessKey,
		"security_token": *tokenResp.Credentials.SessionToken,
	}, map[string]interface{}{
		"username": username,
		"policy":   policy,
		"is_sts":   true,
	})

	// Set the secret TTL to appropriately match the expiration of the token
	resp.Secret.TTL = tokenResp.Credentials.Expiration.Sub(time.Now())

	// STS are purposefully short-lived and aren't renewable
	resp.Secret.Renewable = false

	if usernameWarning != "" {
		resp.AddWarning(usernameWarning)
	}

	return resp, nil
}

func (b *backend) assumeRole(ctx context.Context, s logical.Storage,
	displayName, roleName, roleArn, policy string, policyARNs []string,
	iamGroups []string, lifeTimeInSeconds int64, roleSessionName string) (*logical.Response, error) {

	// grab any IAM group policies associated with the vault role, both inline
	// and managed
	groupPolicies, groupPolicyARNs, err := b.getGroupPolicies(ctx, s, iamGroups)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if len(groupPolicies) > 0 {
		groupPolicies = append(groupPolicies, policy)
		policy, err = combinePolicyDocuments(groupPolicies...)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}
	if len(groupPolicyARNs) > 0 {
		policyARNs = append(policyARNs, groupPolicyARNs...)
	}

	stsClient, err := b.clientSTS(ctx, s)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	roleSessionNameWarning := ""
	if roleSessionName == "" {
		roleSessionName, roleSessionNameWarning = genUsername(displayName, roleName, "iam_user")
	} else {
		roleSessionName = normalizeDisplayName(roleSessionName)
		if len(roleSessionName) > 64 {
			roleSessionName = roleSessionName[0:64]
			roleSessionNameWarning = "the role session name was truncated to 64 characters to fit within IAM session name length limits"
		}
	}

	assumeRoleInput := &sts.AssumeRoleInput{
		RoleSessionName: aws.String(roleSessionName),
		RoleArn:         aws.String(roleArn),
		DurationSeconds: &lifeTimeInSeconds,
	}
	if policy != "" {
		assumeRoleInput.SetPolicy(policy)
	}
	if len(policyARNs) > 0 {
		assumeRoleInput.SetPolicyArns(convertPolicyARNs(policyARNs))
	}
	tokenResp, err := stsClient.AssumeRole(assumeRoleInput)
	if err != nil {
		return logical.ErrorResponse("Error assuming role: %s", err), awsutil.CheckAWSError(err)
	}

	resp := b.Secret(secretAccessKeyType).Response(map[string]interface{}{
		"access_key":     *tokenResp.Credentials.AccessKeyId,
		"secret_key":     *tokenResp.Credentials.SecretAccessKey,
		"security_token": *tokenResp.Credentials.SessionToken,
		"arn":            *tokenResp.AssumedRoleUser.Arn,
	}, map[string]interface{}{
		"username": roleSessionName,
		"policy":   roleArn,
		"is_sts":   true,
	})

	// Set the secret TTL to appropriately match the expiration of the token
	resp.Secret.TTL = tokenResp.Credentials.Expiration.Sub(time.Now())

	// STS are purposefully short-lived and aren't renewable
	resp.Secret.Renewable = false

	if roleSessionNameWarning != "" {
		resp.AddWarning(roleSessionNameWarning)
	}

	return resp, nil
}

func (b *backend) secretAccessKeysCreate(
	ctx context.Context,
	s logical.Storage,
	displayName, policyName string,
	role *awsRoleEntry) (*logical.Response, error) {
	iamClient, err := b.clientIAM(ctx, s)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	username, usernameWarning := genUsername(displayName, policyName, "iam_user")

	// Write to the WAL that this user will be created. We do this before
	// the user is created because if switch the order then the WAL put
	// can fail, which would put us in an awkward position: we have a user
	// we need to rollback but can't put the WAL entry to do the rollback.
	walID, err := framework.PutWAL(ctx, s, "user", &walUser{
		UserName: username,
	})
	if err != nil {
		return nil, fmt.Errorf("error writing WAL entry: %w", err)
	}

	userPath := role.UserPath
	if userPath == "" {
		userPath = "/"
	}

	createUserRequest := &iam.CreateUserInput{
		UserName: aws.String(username),
		Path:     aws.String(userPath),
	}
	if role.PermissionsBoundaryARN != "" {
		createUserRequest.PermissionsBoundary = aws.String(role.PermissionsBoundaryARN)
	}

	// Create the user
	_, err = iamClient.CreateUser(createUserRequest)
	if err != nil {
		if walErr := framework.DeleteWAL(ctx, s, walID); walErr != nil {
			iamErr := fmt.Errorf("error creating IAM user: %w", err)
			return nil, errwrap.Wrap(fmt.Errorf("failed to delete WAL entry: %w", walErr), iamErr)
		}
		return logical.ErrorResponse("Error creating IAM user: %s", err), awsutil.CheckAWSError(err)
	}

	for _, arn := range role.PolicyArns {
		// Attach existing policy against user
		_, err = iamClient.AttachUserPolicy(&iam.AttachUserPolicyInput{
			UserName:  aws.String(username),
			PolicyArn: aws.String(arn),
		})
		if err != nil {
			return logical.ErrorResponse("Error attaching user policy: %s", err), awsutil.CheckAWSError(err)
		}

	}
	if role.PolicyDocument != "" {
		// Add new inline user policy against user
		_, err = iamClient.PutUserPolicy(&iam.PutUserPolicyInput{
			UserName:       aws.String(username),
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(role.PolicyDocument),
		})
		if err != nil {
			return logical.ErrorResponse("Error putting user policy: %s", err), awsutil.CheckAWSError(err)
		}
	}

	for _, group := range role.IAMGroups {
		// Add user to IAM groups
		_, err = iamClient.AddUserToGroup(&iam.AddUserToGroupInput{
			UserName:  aws.String(username),
			GroupName: aws.String(group),
		})
		if err != nil {
			return logical.ErrorResponse("Error adding user to group: %s", err), awsutil.CheckAWSError(err)
		}
	}

	var tags []*iam.Tag
	for key, value := range role.IAMTags {
		// This assignment needs to be done in order to create unique addresses for
		// these variables. Without doing so, all the tags will be copies of the last
		// tag listed in the role.
		k, v := key, value
		tags = append(tags, &iam.Tag{Key: &k, Value: &v})
	}

	if len(tags) > 0 {
		_, err = iamClient.TagUser(&iam.TagUserInput{
			Tags:     tags,
			UserName: &username,
		})

		if err != nil {
			return logical.ErrorResponse("Error adding tags to user: %s", err), awsutil.CheckAWSError(err)
		}
	}

	// Create the keys
	keyResp, err := iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return logical.ErrorResponse("Error creating access keys: %s", err), awsutil.CheckAWSError(err)
	}

	// Remove the WAL entry, we succeeded! If we fail, we don't return
	// the secret because it'll get rolled back anyways, so we have to return
	// an error here.
	if err := framework.DeleteWAL(ctx, s, walID); err != nil {
		return nil, fmt.Errorf("failed to commit WAL entry: %w", err)
	}

	// Return the info!
	resp := b.Secret(secretAccessKeyType).Response(map[string]interface{}{
		"access_key":     *keyResp.AccessKey.AccessKeyId,
		"secret_key":     *keyResp.AccessKey.SecretAccessKey,
		"security_token": nil,
	}, map[string]interface{}{
		"username": username,
		"policy":   role,
		"is_sts":   false,
	})

	lease, err := b.Lease(ctx, s)
	if err != nil || lease == nil {
		lease = &configLease{}
	}

	resp.Secret.TTL = lease.Lease
	resp.Secret.MaxTTL = lease.LeaseMax

	if usernameWarning != "" {
		resp.AddWarning(usernameWarning)
	}

	return resp, nil
}

func (b *backend) secretAccessKeysRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// STS already has a lifetime, and we don't support renewing it
	isSTSRaw, ok := req.Secret.InternalData["is_sts"]
	if ok {
		isSTS, ok := isSTSRaw.(bool)
		if ok {
			if isSTS {
				return nil, nil
			}
		}
	}

	lease, err := b.Lease(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = lease.Lease
	resp.Secret.MaxTTL = lease.LeaseMax
	return resp, nil
}

func (b *backend) secretAccessKeysRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// STS cleans up after itself so we can skip this if is_sts internal data
	// element set to true. If is_sts is not set, assumes old version
	// and defaults to the IAM approach.
	isSTSRaw, ok := req.Secret.InternalData["is_sts"]
	if ok {
		isSTS, ok := isSTSRaw.(bool)
		if ok {
			if isSTS {
				return nil, nil
			}
		} else {
			return nil, fmt.Errorf("secret has is_sts but value could not be understood")
		}
	}

	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}

	// Use the user rollback mechanism to delete this user
	err := b.pathUserRollback(ctx, req, "user", map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func normalizeDisplayName(displayName string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9+=,.@_-]")
	return re.ReplaceAllString(displayName, "_")
}

func convertPolicyARNs(policyARNs []string) []*sts.PolicyDescriptorType {
	size := len(policyARNs)
	retval := make([]*sts.PolicyDescriptorType, size, size)
	for i, arn := range policyARNs {
		retval[i] = &sts.PolicyDescriptorType{
			Arn: aws.String(arn),
		}
	}
	return retval
}
