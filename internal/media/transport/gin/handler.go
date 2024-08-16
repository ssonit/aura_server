package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/media/utils"
)

var (
	cloudinaryCloudName = common.EnvConfig("CLOUDINARY_CLOUD_NAME", "")
	cloudinaryAPIKey    = common.EnvConfig("CLOUDINARY_API_KEY", "")
	cloudinaryAPISecret = common.EnvConfig("CLOUDINARY_API_SECRET", "")
)

type handler struct {
	service utils.MediaService
}

func NewHandler(service utils.MediaService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/upload-image", h.UploadImage())
	group.GET("/get-all-images", h.GetAllImages())
}

func (h *handler) GetAllImages() func(*gin.Context) {
	return func(c *gin.Context) {
		images, err := h.service.GetAllImages(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, images)
	}
}

func generatePublicID() string {
	publicId := fmt.Sprintf("cld-%s", uuid.New().String())
	return publicId
}

func (h *handler) UploadImage() func(*gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
			return
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open the file"})
			return
		}
		defer f.Close()

		fmt.Println(generatePublicID())

		// cld, err := cloudinary.NewFromParams(
		//    cloudinaryCloudName,
		//    cloudinaryAPIKey,
		//    cloudinaryAPISecret,
		// )

		// if err != nil {
		//     log.Fatalf("Failed to initialize Cloudinary, %v", err)
		// }

		id, err := h.service.UploadImage(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(id))
	}
}
