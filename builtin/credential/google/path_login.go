package google

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	oauth "google.golang.org/api/oauth2/v2"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"code": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Google authentication code",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},
	}
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	code := data.Get("code").(string)

	var verifyResp *verifyCredentialsResp
	if verifyResponse, resp, err := b.verifyCredentials(req, code); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		verifyResp = verifyResponse
	}

	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	ttl, _, err := b.SanitizeTTL(config.TTL.String(), config.MaxTTL.String())
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("[ERR]:%s", err)), nil
	}

	return &logical.Response{
		Auth: &logical.Auth{
			InternalData: map[string]interface{}{
				"code": code,
			},
			Policies: verifyResp.Policies,
			Metadata: map[string]string{
				"username": 	verifyResp.User,
				"domain":      	verifyResp.Domain,
			},
			DisplayName: 		verifyResp.Name,
			LeaseOptions: logical.LeaseOptions{
				TTL:       ttl,
				Renewable: true,
			},
		},
	}, nil
}

//TODO: nathang this is probably wrong... we can't renew login with the redeemed code... understand this flow and fix
func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	token := req.Auth.InternalData["code"].(string)

	var verifyResp *verifyCredentialsResp
	if verifyResponse, resp, err := b.verifyCredentials(req, token); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		verifyResp = verifyResponse
	}
	sort.Strings(req.Auth.Policies)
	if !reflect.DeepEqual(verifyResp.Policies, req.Auth.Policies) {
		return logical.ErrorResponse("policies do not match"), nil
	}

	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	return framework.LeaseExtend(config.TTL, config.MaxTTL, b.System())(req, d)
}

func (b *backend) verifyCredentials(req *logical.Request, code string) (*verifyCredentialsResp, *logical.Response, error) {

	config, err := b.Config(req.Storage)

	if err != nil {
		return nil, nil, err
	}

	if config.ApplicationId == "" {
		return nil, logical.ErrorResponse(configErrorMsg), nil
	}

	if config.ApplicationSecret == "" {
		return nil, logical.ErrorResponse(configErrorMsg), nil
	}

	googleConfig := &oauth2.Config{
		ClientID:     config.ApplicationId,
		ClientSecret: config.ApplicationSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{ "email" },
	}

	tok, err := googleConfig.Exchange(oauth2.NoContext, code)
	if (err != nil) {
		return nil, nil, err
	}

	httpClient := googleConfig.Client(context.Background(), tok)
	service, err := oauth.New(httpClient)
	if (err != nil) {
		return nil, nil, err
	}

	me := oauth.NewUserinfoV2MeService(service)
	info, err := me.Get().Do()
	if (err != nil) {
		return nil, nil, err
	}

	user := info.Email
	domain := info.Hd


	if domain != config.Domain {
		return nil, logical.ErrorResponse(fmt.Sprintf("user is of domain %s, not part of required domain %s", domain, config.Domain)), nil
	}

	teamNames := []string{ "default" }
	// Get the teams that this user is part of to determine the policies
	//var teamNames []string
	//
	//teamOpt := &github.ListOptions{
	//	PerPage: 100,
	//}
	//
	//var allTeams []github.Team
	//for {
	//	teams, resp, err := client.Organizations.ListUserTeams(teamOpt)
	//	if err != nil {
	//		return nil, nil, err
	//	}
	//	allTeams = append(allTeams, teams...)
	//	if resp.NextPage == 0 {
	//		break
	//	}
	//	teamOpt.Page = resp.NextPage
	//}
	//
	//for _, t := range allTeams {
	//	// We only care about teams that are part of the organization we use
	//	if *t.Organization.ID != *org.ID {
	//		continue
	//	}
	//
	//	// Append the names so we can get the policies
	//	teamNames = append(teamNames, *t.Name)
	//	if *t.Name != *t.Slug {
	//		teamNames = append(teamNames, *t.Slug)
	//	}
	//}

	policiesList, err := b.Map.Policies(req.Storage, teamNames...)
	if err != nil {
		return nil, nil, err
	}
	return &verifyCredentialsResp{
		User:     user,
		Domain:      domain,
		Policies: policiesList,
	}, nil, nil
}

type verifyCredentialsResp struct {
	User    string
	Domain  string
	Name	string
	Policies []string
}
