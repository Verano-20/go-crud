package controller

import (
	"net/http"
	"strconv"

	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/Verano-20/go-crud/internal/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SimpleController struct {
	SimpleRepository *repository.SimpleRepository
}

func NewSimpleController(db *gorm.DB) *SimpleController {
	return &SimpleController{SimpleRepository: repository.NewSimpleRepository(db)}
}

// Create godoc
// @Summary Create a new Simple
// @Description Create a new Simple with the provided details. The name field is required and must be a non-empty string.
// @Tags Simple
// @Accept json
// @Produce json
// @Param simple body model.SimpleForm true "Simple details" example({"name": "My Simple"})
// @Success 201 {object} response.ApiResponse "Simple created successfully" example({"message": "Simple created successfully", "data": {"id": 1, "name": "My Simple"}})
// @Failure 400 {object} response.ErrorResponse "Invalid request format" example({"error": "Invalid request format"})
// @Failure 500 {object} response.ErrorResponse "Internal server error during resource creation" example({"error": "Failed to create Simple"})
// @Router /simple [post]
func (c *SimpleController) Create(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	var simpleForm model.SimpleForm
	if err := ctx.ShouldBindJSON(&simpleForm); err != nil {
		log.Warn("Invalid create request format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request format"})
		return
	}

	log.Debug("Creating Simple...", zap.Object("simple", &simpleForm))

	simple, err := c.SimpleRepository.Create(simpleForm.ToModel())
	if err != nil {
		log.Error("Failed to create Simple",
			zap.Object("simple", &simpleForm),
			zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to create Simple"})
		return
	}

	log.Debug("Simple created successfully", zap.Object("simple", simple))

	ctx.JSON(http.StatusCreated, response.ApiResponse{Message: "Simple created successfully", Data: simple.ToDTO()})
}

// GetAll godoc
// @Summary Get all Simples
// @Description Get all Simples. Returns an array of Simple objects. Returns an empty array if none exist.
// @Tags Simple
// @Produce json
// @Success 200 {object} response.ApiResponse "Simples retrieved successfully" example({"message": "Simples retrieved successfully", "data": [{"id": 1, "name": "Simple 1"}, {"id": 2, "name": "Simple 2"}]})
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving Simples" example({"error": "Failed to retrieve Simples"})
// @Router /simple [get]
func (c *SimpleController) GetAll(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	log.Debug("Retrieving all Simples")

	var simples model.Simples
	simples, err := c.SimpleRepository.GetAll()
	if err != nil {
		log.Error("Failed to retrieve Simples", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to retrieve Simples"})
		return
	}

	log.Debug("Simples retrieved successfully", zap.Int("count", len(simples)))

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Simples retrieved successfully", Data: simples.ToDTOs()})
}

// GetByID godoc
// @Summary Get Simple by ID
// @Description Find a Simple by its unique ID
// @Tags Simple
// @Param id path int true "Simple ID"
// @Produce json
// @Param id path int true "Simple ID" minimum(1) example(1)
// @Success 200 {object} response.ApiResponse "Simple retrieved successfully" example({"message": "Simple retrieved successfully", "data": {"id": 1, "name": "My Simple"}})
// @Failure 400 {object} response.ErrorResponse "Invalid ID format or value" example({"error": "Invalid ID"})
// @Failure 404 {object} response.ErrorResponse "Simple not found" example({"error": "Simple not found"})
// @Router /simple/{id} [get]
func (c *SimpleController) GetByID(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Warn("Invalid ID format for get by id", zap.String("id_param", idParam), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid ID"})
		return
	}

	log.Debug("Retrieving Simple by ID", zap.Uint64("id", id))

	simple, err := c.SimpleRepository.GetByID(uint(id))
	if err != nil {
		log.Warn("Simple not found", zap.Uint64("id", id), zap.Error(err))
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Simple not found"})
		return
	}

	log.Debug("Simple retrieved successfully", zap.Object("simple", simple))

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Simple retrieved successfully", Data: simple.ToDTO()})
}

// Update godoc
// @Summary Update an existing Simple
// @Description Update a Simple identified by its ID with new data. The ID must exist and the request body must contain valid data.
// @Tags Simple
// @Accept json
// @Produce json
// @Param id path int true "Simple ID to update" minimum(1) example(1)
// @Param simple body model.SimpleForm true "Updated Simple details" example({"name": "Updated Simple"})
// @Success 200 {object} response.ApiResponse "Simple updated successfully" example({"message": "Simple updated successfully", "data": {"id": 1, "name": "Updated Simple"}})
// @Failure 400 {object} response.ErrorResponse "Invalid ID or request body format" example({"error": "Invalid request format"})
// @Failure 404 {object} response.ErrorResponse "Simple not found" example({"error": "Simple not found"})
// @Failure 500 {object} response.ErrorResponse "Internal server error during update operation" example({"error": "Failed to update Simple"})
// @Router /simple/{id} [put]
func (c *SimpleController) Update(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Warn("Invalid ID format for update", zap.String("id_param", idParam), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid ID"})
		return
	}

	var simpleForm model.SimpleForm
	if err := ctx.ShouldBindJSON(&simpleForm); err != nil {
		log.Warn("Invalid update request format", zap.Uint64("id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request format"})
		return
	}

	existingSimple, err := c.SimpleRepository.GetByID(uint(id))
	if err != nil {
		log.Warn("Simple not found for update", zap.Uint64("id", id), zap.Error(err))
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Simple not found"})
		return
	}

	log.Debug("Updating Simple",
		zap.Object("existing", existingSimple),
		zap.Object("update", &simpleForm))

	existingSimple.Name = simpleForm.Name
	simple, err := c.SimpleRepository.Update(existingSimple)
	if err != nil {
		log.Error("Failed to update Simple",
			zap.Object("existing", existingSimple),
			zap.Object("update", &simpleForm),
			zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to update Simple"})
		return
	}

	log.Debug("Simple updated successfully", zap.Object("updated", simple))

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Simple updated successfully", Data: simple.ToDTO()})
}

// Delete godoc
// @Summary Delete a Simple
// @Description Permanently delete a Simple identified by its ID. This operation cannot be undone.
// @Tags Simple
// @Produce json
// @Param id path int true "Simple ID to delete" minimum(1) example(1)
// @Success 200 {object} response.ApiResponse "Simple deleted successfully" example({"message": "Simple deleted successfully", "data": null})
// @Failure 400 {object} response.ErrorResponse "Invalid ID format or value" example({"error": "Invalid ID"})
// @Failure 404 {object} response.ErrorResponse "Simple not found" example({"error": "Simple not found"})
// @Failure 500 {object} response.ErrorResponse "Internal server error during deletion" example({"error": "Failed to delete Simple"})
// @Router /simple/{id} [delete]
func (c *SimpleController) Delete(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Warn("Invalid ID format for delete", zap.String("id_param", idParam), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid ID"})
		return
	}

	existingSimple, err := c.SimpleRepository.GetByID(uint(id))
	if err != nil {
		log.Warn("Simple not found for deletion", zap.Uint64("id", id), zap.Error(err))
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Simple not found"})
		return
	}

	log.Debug("Deleting Simple", zap.Object("simple", existingSimple))

	err = c.SimpleRepository.Delete(uint(id))
	if err != nil {
		log.Error("Failed to delete Simple",
			zap.Object("simple", existingSimple),
			zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to delete Simple"})
		return
	}

	log.Debug("Simple deleted successfully", zap.Object("simple", existingSimple))

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Simple deleted successfully", Data: nil})
}
