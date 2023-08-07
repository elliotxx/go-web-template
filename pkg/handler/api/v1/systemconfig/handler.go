package systemconfig

import (
	"context"
	"strconv"

	"github.com/elliotxx/errors"
	"github.com/elliotxx/go-web-template/pkg/domain/entity"
	"github.com/elliotxx/go-web-template/pkg/domain/repository"
	"github.com/elliotxx/go-web-template/pkg/errcode"
	"github.com/elliotxx/go-web-template/pkg/util/kdump"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	repo repository.SystemConfigRepository
}

func NewHandler(repo repository.SystemConfigRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// @Summary      Create system config
// @Description  Create a new system config instance
// @Accept       json
// @Produce      json
// @Param        system  config    body                 CreateSystemConfigRequest  true  "Created system config"
// @Success      200     {object}  entity.SystemConfig  "Success"
// @Failure      400     {object}  errors.DetailError   "Bad Request"
// @Failure      401     {object}  errors.DetailError   "Unauthorized"
// @Failure      429     {object}  errors.DetailError   "Too Many Requests"
// @Failure      404     {object}  errors.DetailError   "Not Found"
// @Failure      500     {object}  errors.DetailError   "Internal Server Error"
// @Router       /api/v1/systemconfig [post]
func (h *Handler) CreateSystemConfig(c *gin.Context, log logrus.FieldLogger) (any, error) {
	// Parse payload from requestPayload
	var requestPayload CreateSystemConfigRequest
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		return nil, errcode.ErrDeserializedParams.Causewf(err, "failed to decode json")
	}
	log.Infof("Request payload: %v", kdump.FormatN(requestPayload))

	// Convert request payload to domain model
	var systemConfig entity.SystemConfig
	if err := copier.Copy(&systemConfig, &requestPayload); err != nil {
		return nil, errors.Wrap(err, "failed to convert request payload to domain model")
	}

	// Create systemConfig with repository
	err := h.repo.Create(context.TODO(), &systemConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to creating systemConfig with repository")
	}

	// Return created systemConfig
	return systemConfig, nil
}

// @Summary      Delete system config
// @Description  Delete specified system config by ID
// @Produce      json
// @Param        id   path      int                  true  "SystemConfig ID"
// @Success      200  {object}  entity.SystemConfig  "Success"
// @Failure      400  {object}  errors.DetailError   "Bad Request"
// @Failure      401  {object}  errors.DetailError   "Unauthorized"
// @Failure      429  {object}  errors.DetailError   "Too Many Requests"
// @Failure      404  {object}  errors.DetailError   "Not Found"
// @Failure      500  {object}  errors.DetailError   "Internal Server Error"
// @Router       /api/v1/systemconfig/{id} [delete]
func (h *Handler) DeleteSystemConfig(c *gin.Context, log logrus.FieldLogger) (any, error) {
	// Parse payload from requestPayload
	paramID := c.Param("id")
	log.Infof("Request params id: %s", paramID)

	// Get systemConfig with repository
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return nil, err
	}

	// Delete systemConfig with repository
	err = h.repo.Delete(context.TODO(), uint(id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to deleting systemConfig with repository")
	}

	// Return deleted systemConfig
	return nil, nil
}

// @Summary      Update system config
// @Description  Update the specified system config
// @Accept       json
// @Produce      json
// @Param        system  config    body                 UpdateSystemConfigRequest  true  "Updated system config"
// @Success      200     {object}  entity.SystemConfig  "Success"
// @Failure      400     {object}  errors.DetailError   "Bad Request"
// @Failure      401     {object}  errors.DetailError   "Unauthorized"
// @Failure      429     {object}  errors.DetailError   "Too Many Requests"
// @Failure      404     {object}  errors.DetailError   "Not Found"
// @Failure      500     {object}  errors.DetailError   "Internal Server Error"
// @Router       /api/v1/systemconfig [put]
func (h *Handler) UpdateSystemConfig(c *gin.Context, log logrus.FieldLogger) (any, error) {
	// Parse payload from requestPayload
	var requestPayload UpdateSystemConfigRequest
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		return nil, errcode.ErrDeserializedParams.Causewf(err, "failed to decode json")
	}
	log.Infof("Request payload: %v", kdump.FormatN(requestPayload))

	// Convert request payload to domain model
	var requestEntity entity.SystemConfig
	if err := copier.Copy(&requestEntity, &requestPayload); err != nil {
		return nil, errors.Wrap(err, "failed to convert request payload to domain model")
	}

	// Get the existed systemConfig by id
	updatedEntity, err := h.repo.Get(context.TODO(), requestEntity.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.NotFound.Causewf(err, "failed to update system config")
		}
		return nil, errcode.InvalidParams.Cause(err)
	}

	// Overwrite non-zero values in request entity to existed entity
	copier.CopyWithOption(updatedEntity, requestEntity, copier.Option{IgnoreEmpty: true})

	// Update systemConfig with repository
	err = h.repo.Update(context.TODO(), updatedEntity)
	if err != nil {
		return nil, errors.Wrap(err, "failed to updating systemConfig with repository")
	}

	// Return updated systemConfig
	return updatedEntity, nil
}

// @Summary      Get system config
// @Description  Get system config information by system config ID
// @Produce      json
// @Param        id   path      int                  true  "SystemConfig ID"
// @Success      200  {object}  entity.SystemConfig  "Success"
// @Failure      400  {object}  errors.DetailError   "Bad Request"
// @Failure      401  {object}  errors.DetailError   "Unauthorized"
// @Failure      429  {object}  errors.DetailError   "Too Many Requests"
// @Failure      404  {object}  errors.DetailError   "Not Found"
// @Failure      500  {object}  errors.DetailError   "Internal Server Error"
// @Router       /api/v1/systemconfig/{id} [get]
func (h *Handler) GetSystemConfig(c *gin.Context, log logrus.FieldLogger) (any, error) {
	// Parse payload from request
	paramID := c.Param("id")
	log.Infof("Request params id: %s", paramID)

	// Get systemConfig with repository
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return nil, err
	}
	existedEntity, err := h.repo.Get(context.TODO(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.NotFound.Causewf(err, "failed to get system config")
		}
		return nil, errors.Wrap(err, "failed to get systemConfig with repository")
	}

	// Return systemConfig
	return existedEntity, nil
}

// @Summary      Find system configs
// @Description  Find system configs with query
// @Accept       json
// @Produce      json
// @Param        query  body      QuerySystemConfigRequest  true  "query body"
// @Success      200    {object}  entity.SystemConfig       "Success"
// @Failure      400    {object}  errors.DetailError        "Bad Request"
// @Failure      401    {object}  errors.DetailError        "Unauthorized"
// @Failure      429    {object}  errors.DetailError        "Too Many Requests"
// @Failure      404    {object}  errors.DetailError        "Not Found"
// @Failure      500    {object}  errors.DetailError        "Internal Server Error"
// @Router       /api/v1/systemconfigs [get]
func (h *Handler) FindSystemConfigs(c *gin.Context, log logrus.FieldLogger) (any, error) {
	// Parse payload from request
	var requestPayload QuerySystemConfigRequest
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		return nil, errcode.ErrDeserializedParams.Causewf(err, "failed to decode json")
	}
	log.Infof("Request payload: %v", kdump.FormatN(requestPayload))

	// Calculate the limit and offset based on the pagination request
	limit := requestPayload.PerPage
	offset := (requestPayload.Page - 1) * requestPayload.PerPage

	// Find systemConfigs with repository
	dataEntities, err := h.repo.Find(context.TODO(), repository.Query{
		Offset:  offset,
		Limit:   limit,
		Keyword: requestPayload.Keyword,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all systemConfig with repository")
	}

	// Return all systemConfig
	return dataEntities, nil
}

// @Summary      Count system configs
// @Description  Count the total number of system configs
// @Produce      json
// @Success      200  {object}  entity.SystemConfig  "Success"
// @Failure      400  {object}  errors.DetailError   "Bad Request"
// @Failure      401  {object}  errors.DetailError   "Unauthorized"
// @Failure      429  {object}  errors.DetailError   "Too Many Requests"
// @Failure      404  {object}  errors.DetailError   "Not Found"
// @Failure      500  {object}  errors.DetailError   "Internal Server Error"
// @Router       /api/v1/systemconfig/count [get]
func (h *Handler) CountSystemConfigs(c *gin.Context, log logrus.FieldLogger) (any, error) {
	// Count systemConfigs with repository
	total, err := h.repo.Count(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "failed to count systemConfig with repository")
	}

	// Return total of all systemConfig
	return CountSystemConfigResponse{
		Total: total,
	}, nil
}
