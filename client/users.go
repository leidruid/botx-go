package client

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

type searchUserResult struct {
	UserHUID        uuid.UUID  `json:"user_huid"`
	ADLogin         string     `json:"ad_login"`
	ADDomain        string     `json:"ad_domain"`
	Name            string     `json:"name"`
	Company         string     `json:"company"`
	CompanyPosition string     `json:"company_position"`
	Department      string     `json:"department"`
	Emails          []string   `json:"emails"`
	OtherID         string     `json:"other_id"`
	UserKind        string     `json:"user_kind"`
	Active          *bool      `json:"active"`
	Description     string     `json:"description"`
	IPPhone         string     `json:"ip_phone"`
	Manager         string     `json:"manager"`
	Office          string     `json:"office"`
	OtherIPPhone    string     `json:"other_ip_phone"`
	OtherPhone      string     `json:"other_phone"`
	PublicName      string     `json:"public_name"`
	CTSID           *uuid.UUID `json:"cts_id"`
	RTSID           *uuid.UUID `json:"rts_id"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type searchUserResponse struct {
	Status string           `json:"status"`
	Result searchUserResult `json:"result"`
}

type searchUsersResponse struct {
	Status string             `json:"status"`
	Result []searchUserResult `json:"result"`
}

// SearchUserByHUID returns user profile by HUID.
func (c *Client) SearchUserByHUID(ctx context.Context, botID uuid.UUID, huid uuid.UUID) (models.UserFromSearch, error) {
	path := "/api/v3/botx/users/by_huid"
	return c.searchUserGET(ctx, botID, path, url.Values{"user_huid": []string{huid.String()}})
}

// SearchUserByEmail returns user profile by email.
func (c *Client) SearchUserByEmail(ctx context.Context, botID uuid.UUID, email string) (models.UserFromSearch, error) {
	path := "/api/v3/botx/users/by_email"
	return c.searchUserGET(ctx, botID, path, url.Values{"email": []string{email}})
}

// SearchUserByEmails returns user profiles by emails.
func (c *Client) SearchUserByEmails(ctx context.Context, botID uuid.UUID, emails []string) ([]models.UserFromSearch, error) {
	path := "/api/v3/botx/users/by_email"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return nil, err
	}
	payload := map[string]any{"emails": emails}
	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return nil, err
	}
	var resp searchUsersResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != "ok" {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	out := make([]models.UserFromSearch, 0, len(resp.Result))
	for _, r := range resp.Result {
		out = append(out, mapSearchUser(r))
	}
	return out, nil
}

// SearchUserByLogin returns user profile by AD login/domain.
func (c *Client) SearchUserByLogin(ctx context.Context, botID uuid.UUID, adLogin string, adDomain string) (models.UserFromSearch, error) {
	path := "/api/v3/botx/users/by_login"
	params := url.Values{"ad_login": []string{adLogin}, "ad_domain": []string{adDomain}}
	return c.searchUserGET(ctx, botID, path, params)
}

// SearchUserByOtherID returns user profile by other id.
func (c *Client) SearchUserByOtherID(ctx context.Context, botID uuid.UUID, otherID string) (models.UserFromSearch, error) {
	path := "/api/v3/botx/users/by_other_id"
	return c.searchUserGET(ctx, botID, path, url.Values{"other_id": []string{otherID}})
}

func (c *Client) searchUserGET(ctx context.Context, botID uuid.UUID, path string, params url.Values) (models.UserFromSearch, error) {
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.UserFromSearch{}, err
	}
	if enc := params.Encode(); enc != "" {
		urlStr = urlStr + "?" + enc
	}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return models.UserFromSearch{}, err
	}
	var resp searchUserResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.UserFromSearch{}, err
	}
	if resp.Status != "ok" {
		return models.UserFromSearch{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return mapSearchUser(resp.Result), nil
}

func mapSearchUser(r searchUserResult) models.UserFromSearch {
	return models.UserFromSearch{
		HUID:            r.UserHUID,
		ADLogin:         r.ADLogin,
		ADDomain:        r.ADDomain,
		Username:        r.Name,
		Company:         r.Company,
		CompanyPosition: r.CompanyPosition,
		Department:      r.Department,
		Emails:          r.Emails,
		OtherID:         r.OtherID,
		UserKind:        models.UserKind(r.UserKind),
		Active:          r.Active,
		Description:     r.Description,
		IPPhone:         r.IPPhone,
		Manager:         r.Manager,
		Office:          r.Office,
		OtherIPPhone:    r.OtherIPPhone,
		OtherPhone:      r.OtherPhone,
		PublicName:      r.PublicName,
		CTSID:           r.CTSID,
		RTSID:           r.RTSID,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

// UpdateUserProfile updates user profile fields.
func (c *Client) UpdateUserProfile(
	ctx context.Context,
	botID uuid.UUID,
	userHUID uuid.UUID,
	avatar optional.Optional[models.OutgoingAttachment],
	name optional.Optional[string],
	publicName optional.Optional[string],
	company optional.Optional[string],
	companyPosition optional.Optional[string],
	description optional.Optional[string],
	department optional.Optional[string],
	office optional.Optional[string],
	manager optional.Optional[string],
) error {
	path := "/api/v3/botx/users/update_profile"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{"user_huid": userHUID}
	if name.Set {
		payload["name"] = name.Value
	}
	if publicName.Set {
		payload["public_name"] = publicName.Value
	}
	if company.Set {
		payload["company"] = company.Value
	}
	if companyPosition.Set {
		payload["company_position"] = companyPosition.Value
	}
	if description.Set {
		payload["description"] = description.Value
	}
	if department.Set {
		payload["department"] = department.Value
	}
	if office.Set {
		payload["office"] = office.Value
	}
	if manager.Set {
		payload["manager"] = manager.Value
	}
	if avatar.Set {
		apiAtt, err := models.AttachmentToAPI(avatar.Value)
		if err != nil {
			return err
		}
		payload["avatar"] = apiAtt.Data
	}

	req, err := BuildJSONRequest(http.MethodPut, urlStr, payload)
	if err != nil {
		return err
	}
	var resp struct {
		Status string `json:"status"`
		Result bool   `json:"result"`
	}
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return err
	}
	if resp.Status != "ok" {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}

// UsersAsCSV downloads users list as CSV and parses it.
func (c *Client) UsersAsCSV(ctx context.Context, botID uuid.UUID, ctsUser bool, unregistered bool, botx bool) ([]models.UserFromCSV, error) {
	path := "/api/v3/botx/users/users_as_csv"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return nil, err
	}
	params := url.Values{
		"cts_user":     []string{fmt.Sprintf("%t", ctsUser)},
		"unregistered": []string{fmt.Sprintf("%t", unregistered)},
		"botx":         []string{fmt.Sprintf("%t", botx)},
	}
	urlStr = urlStr + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.doAuthorizedJSON(ctx, botID, req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseUsersCSV(resp.Body)
}

func parseUsersCSV(r io.Reader) ([]models.UserFromCSV, error) {
	csvr := csv.NewReader(r)
	csvr.TrimLeadingSpace = true
	headers, err := csvr.Read()
	if err != nil {
		return nil, err
	}
	index := map[string]int{}
	for i, h := range headers {
		index[strings.TrimSpace(h)] = i
	}

	var users []models.UserFromCSV
	for {
		rec, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		get := func(key string) string {
			if i, ok := index[key]; ok && i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}

		huid, _ := uuid.Parse(get("HUID"))
		managerHUID := (*uuid.UUID)(nil)
		if v := get("Manager HUID"); v != "" {
			if id, err := uuid.Parse(v); err == nil {
				managerHUID = &id
			}
		}
		active := strings.EqualFold(get("Active"), "true")

		users = append(users, models.UserFromCSV{
			HUID:            huid,
			ADLogin:         get("AD Login"),
			ADDomain:        get("Domain"),
			Email:           get("AD E-mail"),
			Username:        get("Name"),
			SyncSource:      models.SyncSourceType(get("Sync source")),
			Active:          active,
			UserKind:        models.UserKind(get("Kind")),
			Company:         get("Company"),
			Department:      get("Department"),
			Position:        get("Position"),
			Avatar:          get("Avatar"),
			AvatarPreview:   get("Avatar preview"),
			Office:          get("Office"),
			Manager:         get("Manager"),
			ManagerHUID:     managerHUID,
			Description:     get("Description"),
			Phone:           get("Phone"),
			OtherPhone:      get("Other phone"),
			IPPhone:         get("IP phone"),
			OtherIPPhone:    get("Other IP phone"),
			PersonnelNumber: get("Personnel number"),
		})
	}
	return users, nil
}
