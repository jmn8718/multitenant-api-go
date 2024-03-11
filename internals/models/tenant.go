package models

type Tenant struct {
	BaseModel
	Name      string `json:"name"`
	ApiKey    string `json:"apiKey"`
	Image     string `json:"image"`
	Status    string `json:"status" gorm:"default:ACTIVE"`
	IsSandbox bool   `json:"isSandbox" gorm:"<-:create,not null"`
}

type TenantCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type TenantCreateResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MyTenantsData struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	IsSandbox bool   `json:"isSandbox"`
}

type MyTenantsDataResponse struct {
	Count int64           `json:"count"`
	Items []MyTenantsData `json:"items"`
}

type TenantApiKeyResponse struct {
	ApiKey string `json:"apiKey"`
}
