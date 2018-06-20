package models

type UserData struct {
	Id                 int    `json:"id"`
	IdmGuid            string `json:"idm_guid"`
	Username           string `json:"username"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	OrgId              int    `json:"org_id"`
	CrmUid             string `json:"crm_uid"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	AllowImpersonation bool   `json:"allowImpersonation"`
}

type NoUser struct {
	NoFound string `json:"not_found"`
}
