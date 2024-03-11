package repositories

import (
	"multitenant-api-go/internals/constants"
	"multitenant-api-go/internals/errors"
	"multitenant-api-go/internals/models"
	"multitenant-api-go/internals/utils"
	"time"

	"gorm.io/gorm"
)

type TenantsRepository struct {
	db *gorm.DB
}

func NewTenantsRepository(db *gorm.DB) TenantsRepository {
	return TenantsRepository{db: db}
}

func (tenantsRepository TenantsRepository) CreateTenantWithUser(body models.TenantCreateRequest, userId string) (*models.TenantCreateResponse, error) {
	var count int64
	if err := tenantsRepository.db.Model(&models.UserTenants{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.NewHttpError(400, "user_reached_tenant_limit")
	}

	apiKey, err := utils.GenerateApiKey(12, "sb")
	if err != nil {
		return nil, errors.NewHttpErrorWithReason(500, "internal_error", err.Error(), false)
	}
	dbTenant := &models.Tenant{Name: body.Name, ApiKey: apiKey, IsSandbox: true}

	dbErr := tenantsRepository.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dbTenant).Error; err != nil {
			return errors.NewHttpErrorWithReason(500, "internal_error", err.Error(), false)
		}
		dbUserTenant := models.UserTenants{UserId: userId, TenantId: dbTenant.ID, Role: constants.ADMIN, AcceptedInvitation: true}
		if err := tx.Create(&dbUserTenant).Error; err != nil {
			return errors.NewHttpErrorWithReason(500, "internal_error", err.Error(), false)
		}
		return nil
	})

	if dbErr != nil {
		return nil, dbErr
	}
	return &models.TenantCreateResponse{Id: dbTenant.ID, Name: dbTenant.Name}, nil
}

func (tenantsRepository TenantsRepository) GetMyTenants(userId string) (*models.MyTenantsDataResponse, error) {
	var count int64
	var items []models.MyTenantsData

	if err := tenantsRepository.db.Model(&models.UserTenants{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return nil, err
	}

	if err := tenantsRepository.db.Table("user_tenants as u").Where("u.user_id = ?", userId).Joins("left join tenants as t on t.id = u.tenant_id").Select("t.id as id, t.name as name, u.role as role, t.is_sandbox as is_sandbox").Find(&items).Error; err != nil {
		return nil, err
	}
	return &models.MyTenantsDataResponse{Count: count, Items: items}, nil
}

func (tenantsRepository TenantsRepository) GetMyTenantUsers(userId string, tenantId string) (*models.TenantUsersResponse, error) {
	var users []models.TenantUserData
	var userTenant models.UserTenants

	if err := tenantsRepository.db.Where("user_id = ? AND tenant_id = ?", userId, tenantId).Select("role").First(&userTenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewHttpError(404, "not_found")
		} else {
			return nil, errors.NewHttpError(500, "internal_error")
		}
	}

	if err := tenantsRepository.db.Table("user_tenants as t").Where("t.tenant_id = ?", tenantId).Joins("left join users as u on t.user_id = u.id").Select("u.id as id, t.role as role, t.invited_at as invited_at,t.accepted_invitation as accepted_invitation, u.image as image, u.name as name, u.email as email").Find(&users).Error; err != nil {
		return nil, err
	}
	var items = []models.TenantUserDataResponse{}
	for _, user := range users {
		items = append(items, user.ToTenantUserDataResponse())
	}
	return &models.TenantUsersResponse{Count: len(items), Items: items}, nil
}

func (tenantsRepository TenantsRepository) AddUserToMyTenant(userId string, tenantId string, email string) error {
	dbErr := tenantsRepository.db.Transaction(func(tx *gorm.DB) error {
		var tenant models.Tenant
		if err := tx.Where("id = ?", tenantId).First(&tenant).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.NewHttpError(404, "tenant_not_found")
			}
			return errors.NewHttpError(500, "db_user_ex")
		}
		var dbUser models.User
		if err := tx.Where("email = ?", email).First(&dbUser).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.NewHttpError(500, "db_user_ex")
			}
		}
		if dbUser.ID == "" {
			newDbUser := models.User{Email: email, Password: "", Name: email, EmailVerified: false}
			result := tx.Create(&newDbUser)
			if result.Error != nil {
				return errors.NewHttpError(500, "db_user")
			}
			dbUser = newDbUser
		}
		var dbUserTenant models.UserTenants
		result := tx.Where("user_id = ? AND tenant_id = ?", dbUser.ID, tenantId).First(&dbUserTenant)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.NewHttpError(500, "db_user_ex")
		}
		if result.RowsAffected == 0 {
			dbUserTenant := models.UserTenants{UserId: dbUser.ID, TenantId: tenantId, Role: "USER", InvitedAt: time.Now()}
			if err := tx.Create(&dbUserTenant).Error; err != nil {
				return errors.NewHttpError(500, "db_tenant_user")
			}
		}
		return nil
	})

	return dbErr
}

func (tenantsRepository TenantsRepository) HasUserAccessToTenant(userId string, tenantId string) (*models.UserTenants, error) {
	var tenant models.Tenant
	if err := tenantsRepository.db.Where("id = ?", tenantId).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewHttpError(404, "not_found")
		} else {
			return nil, errors.NewHttpError(500, "internal_error")
		}
	}
	var userTenant models.UserTenants
	if err := tenantsRepository.db.Where("user_id = ? AND tenant_id = ?", userId, tenantId).First(&userTenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewHttpError(404, "not_found")
		} else {
			return nil, errors.NewHttpError(500, "internal_error")
		}
	}
	return &userTenant, nil
}

func (tenantsRepository TenantsRepository) GetTenantByTenantId(tenantId string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := tenantsRepository.db.Where("id = ?", tenantId).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewHttpError(404, "not_found")
		} else {
			return nil, errors.NewHttpError(500, "internal_error")
		}
	}
	return &tenant, nil
}

func (tenantsRepository TenantsRepository) GetMyTenantApiKey(tenantId string) (*models.TenantApiKeyResponse, error) {
	var tenant models.Tenant
	if err := tenantsRepository.db.Where("id = ?", tenantId).Select("api_key").First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewHttpError(404, "not_found")
		} else {
			return nil, errors.NewHttpErrorWithReason(500, "internal_error", err.Error(), false)
		}
	}
	return &models.TenantApiKeyResponse{ApiKey: tenant.ApiKey}, nil
}

func (tenantsRepository TenantsRepository) UpdateMyTenantApiKey(tenantId string) (*models.TenantApiKeyResponse, error) {
	tenant, err := tenantsRepository.GetTenantByTenantId(tenantId)
	if err != nil {
		return nil, err
	}
	var apiKeyPrefix = "sb"
	if !tenant.IsSandbox {
		apiKeyPrefix = "lv"
	}
	apiKey, err := utils.GenerateApiKey(12, apiKeyPrefix)
	if err != nil {
		return nil, errors.NewHttpErrorWithReason(500, "internal_error", err.Error(), false)
	}
	if err := tenantsRepository.db.Model(&models.Tenant{}).Where("id = ?", tenantId).Update("api_key", apiKey).Error; err != nil {
		return nil, err
	}
	return &models.TenantApiKeyResponse{ApiKey: apiKey}, nil
}
