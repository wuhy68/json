package godropbox

import (
	"net/http"

	"encoding/json"

	"fmt"

	"github.com/joaosoft/go-error/service"
	"github.com/joaosoft/go-manager/service"
)

type user struct {
	client gomanager.IGateway
	config *goDropboxConfig
}

type getUserResponse struct {
	AccountID string `json:"account_id"`
	Name      struct {
		GivenName       string `json:"given_name"`
		Surname         string `json:"surname"`
		FamiliarName    string `json:"familiar_name"`
		DisplayName     string `json:"display_name"`
		AbbreviatedName string `json:"abbreviated_name"`
	} `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Disabled      bool   `json:"disabled"`
	Locale        string `json:"locale"`
	ReferralLink  string `json:"referral_link"`
	IsPaired      bool   `json:"is_paired"`
	AccountType   struct {
		Tag string `json:".tag"`
	} `json:"account_type"`
	RootInfo struct {
		Tag             string `json:".tag"`
		RootNamespaceID string `json:"root_namespace_id"`
		HomeNamespaceID string `json:"home_namespace_id"`
	} `json:"root_info"`
	Country string `json:"country"`
	Team    struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		SharingPolicies struct {
			SharedFolderMemberPolicy struct {
				Tag string `json:".tag"`
			} `json:"shared_folder_member_policy"`
			SharedFolderJoinPolicy struct {
				Tag string `json:".tag"`
			} `json:"shared_folder_join_policy"`
			SharedLinkCreatePolicy struct {
				Tag string `json:".tag"`
			} `json:"shared_link_create_policy"`
		} `json:"sharing_policies"`
		OfficeAddinPolicy struct {
			Tag string `json:".tag"`
		} `json:"office_addin_policy"`
	} `json:"team"`
	TeamMemberID string `json:"team_member_id"`
}

// Get ...
func (u *user) Get() (*getUserResponse, *goerror.ErrorData) {
	dropboxResponse := &getUserResponse{}
	headers := gomanager.Headers{
		"Authorization": {fmt.Sprintf("%s %s", u.config.Authorization.Access, u.config.Authorization.Token)},
	}

	if status, response, err := u.client.Request(http.MethodPost, u.config.Hosts.Api, "/users/get_current_account", headers, nil); err != nil {
		newErr := goerror.NewError(err)
		log.Error("error getting user account").ToErrorData(newErr)
		return nil, newErr
	} else if status != http.StatusOK {
		var err error
		log.Errorf("response status %d instead of %d", status, http.StatusOK).ToError(&err)
		return nil, goerror.NewError(err)
	} else if response == nil {
		var err error
		log.Error("error getting user account").ToError(&err)
		return nil, goerror.NewError(err)
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			newErr := goerror.NewError(err)
			log.Error("error converting dropbox user data").ToErrorData(newErr)
			return nil, newErr
		}
		return dropboxResponse, nil
	}

	return nil, nil
}
