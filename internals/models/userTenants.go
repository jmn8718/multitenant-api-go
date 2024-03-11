package models

import (
	"multitenant-api-go/internals/utils"
	"time"
)

type UserTenants struct {
	UserId             string    `gorm:"unique_index:user_tenant,priority:2"`
	TenantId           string    `gorm:"unique_index:user_tenant,priority:1"`
	Role               string    `json:"role" gorm:"required,default:USER"`
	Image              string    `json:"image"`
	InvitedAt          time.Time `json:"invitedAt"`
	AcceptedInvitation bool      `json:"acceptedInvitation" gorm:"default:false"`
}

type TenantUserData struct {
	Id                 string    `json:"id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Role               string    `json:"role"`
	Image              string    `json:"image"`
	InvitedAt          time.Time `json:"invitedAt"`
	AcceptedInvitation bool      `json:"acceptedInvitation"`
}

type TenantUserDataResponse struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	Image              string `json:"image"`
	InvitedAt          string `json:"invitedAt"`
	AcceptedInvitation bool   `json:"acceptedInvitation"`
}

func (tenantUser TenantUserData) ToTenantUserDataResponse() TenantUserDataResponse {
	return TenantUserDataResponse{
		Id:                 tenantUser.Id,
		Name:               tenantUser.Name,
		Email:              tenantUser.Email,
		Role:               tenantUser.Role,
		Image:              tenantUser.Image,
		InvitedAt:          utils.ToValidDateString(tenantUser.InvitedAt),
		AcceptedInvitation: tenantUser.AcceptedInvitation,
	}
}

type TenantUsersResponse struct {
	Count int                      `json:"count"`
	Items []TenantUserDataResponse `json:"items"`
}

type TenantAddUserRequest struct {
	Email string `json:"email" binding:"required"`
}
