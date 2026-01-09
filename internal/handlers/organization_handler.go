// internal/handlers/organization_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/mappers"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/services"
)

type OrganizationHandler struct {
	*BaseHandler[models.Organization, dto.CreateOrganizationRequest, dto.UpdateOrganizationRequest, dto.OrganizationResponse]
	orgService services.OrganizationService
	mapper     *mappers.UnifiedMapper
}

func NewOrganizationHandler(
	orgService services.OrganizationService,
	mapper *mappers.UnifiedMapper,
) *OrganizationHandler {
	base := NewBaseHandler[models.Organization, dto.CreateOrganizationRequest, dto.UpdateOrganizationRequest, dto.OrganizationResponse](
		orgService,
		mapper,
		"organization",
	)

	base.SetLimits(20, 100)

	return &OrganizationHandler{
		BaseHandler: base,
		orgService:  orgService,
		mapper:      mapper,
	}
}

// VerifyOrganization verifica una organización (admin only)
func (h *OrganizationHandler) VerifyOrganization(c *gin.Context) {
	orgID := c.Param("id")
	userCtx := extractUserContext(c)

	if !userCtx.IsAdmin() {
		common.ErrorResponse(c, common.NewBusinessError("admin_required", "Solo administradores pueden verificar organizaciones"))
		return
	}

	org, err := h.orgService.VerifyOrganization(c.Request.Context(), orgID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	response := h.mapper.OrganizationToResponse(org, userCtx)
	common.SuccessResponse(c, http.StatusOK, "Organización verificada exitosamente", response)
}

// GetMembers obtiene miembros de la organización
func (h *OrganizationHandler) GetMembers(c *gin.Context) {
	orgID := c.Param("id")
	opts := extractQueryOptions(c)
	userCtx := extractUserContext(c)

	users, pagination, err := h.orgService.GetMembers(c.Request.Context(), orgID, *opts, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	response := h.mapper.UsersToListResponse(users, pagination, userCtx)
	common.SuccessWithPagination(c, "Miembros de la organización", response.Users, pagination)
}

// GetActiveOrganizations obtiene organizaciones activas
func (h *OrganizationHandler) GetActiveOrganizations(c *gin.Context) {
	opts := extractQueryOptions(c)

	orgs, pagination, err := h.orgService.GetActiveOrganizations(c.Request.Context(), *opts)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// No hay userCtx porque es público
	response := h.mapper.OrganizationsToListResponse(orgs, pagination, nil)
	common.SuccessWithPagination(c, "Organizaciones activas", response.Organizations, pagination)
}
