package utils

import (
	"context"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

type UploadResult struct {
	URL      string
	PublicID string
}

func NewCloudinaryService() *CloudinaryService {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)

	if err != nil {
		panic("Cloudinary init failed: " + err.Error())
	}

	return &CloudinaryService{
		cld: cld,
	}
}

func (c *CloudinaryService) UploadFile(
	ctx context.Context,
	filePath string,
	folder string,
	mimeType string,
) (UploadResult, error) {

	resourceType := "image"

	if strings.HasPrefix(mimeType, "video") {
		resourceType = "video"
	}

	res, err := c.cld.Upload.Upload(
		ctx,
		filePath,
		uploader.UploadParams{
			Folder:       folder,
			ResourceType: resourceType,
		},
	)

	if err != nil {
		return UploadResult{}, err
	}

	return UploadResult{
		URL:      res.SecureURL,
		PublicID: res.PublicID,
	}, nil
}

func (c *CloudinaryService) Delete(publicID string) error {

	_, err := c.cld.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)

	return err
}
