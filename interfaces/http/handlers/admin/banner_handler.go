package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	responses "shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"
	"shop/interfaces/http/response"
	"shop/pkg/errors"
	"shop/pkg/old"
	"shop/pkg/sessions"
	"shop/pkg/util"
	"slices"
)

// IndexBanner lists every banner (sidebar: لیست بنر ها).
func (a *AdminHandler) IndexBanner(c *gin.Context) {
	banners, err := a.bannerSrv.Index(c)
	if err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/banners/create")
		return
	}

	response.Render(c, http.StatusOK, "admin_index_banner", gin.H{
		"TITLE":      "لیست بنرها",
		"BANNERS":    responses.ToBanners(banners),
		"MEDIA_PATH": util.GetBannerStoragePath(),
	})
	return
}

func (a *AdminHandler) CreateBanner(c *gin.Context) {
	response.Render(c, 200, "create_banner", gin.H{"TITLE": "ایجاد بنر"})
}

func (a *AdminHandler) StoreBanner(c *gin.Context) {
	var req requests.CreateBannerRequest
	_ = c.Request.ParseForm()
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/banners/create")
		return
	}

	//check slider type validation
	if !entities.IsValidBannerType(req.Type) {
		errors.Init()
		errors.Add("type", domain_err.SliderTypeIsNotValid)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/banners/create")
		return
	}

	imageFile, _ := c.FormFile("image")

	//check required validation
	if imageFile == nil {
		errors.Init()
		errors.Add("image", domain_err.IsRequired)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/banners/create")
		return
	}

	extension := filepath.Ext(imageFile.Filename)

	// file extension validation
	ok := slices.Contains(util.AllowImageExtensions(), extension)
	if !ok {
		errors.Init()
		errors.Add("image", domain_err.MustBeImage)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/banners/create")
		return

	}

	//generate file name
	imageGenerateFileName := util.GenerateFilename(imageFile.Filename)

	//check upload and store file to storage bucket
	if os.Getenv("STORAGE_STATUS") == "active" {
		go func() {
			{
				imageFile, _ := imageFile.Open()
				err := a.dep.Storage.UploadFile(imageFile, os.Getenv("STORAGE_BANNER_PATH")+imageGenerateFileName)
				if err != nil {
					log.Println("-- failed to upload file to s3 : ", err)
				}
			}
		}()
	}

	//store images on disk
	saveUploadedImageErr := c.SaveUploadedFile(imageFile, viper.GetString("Upload.Banners")+imageGenerateFileName)
	if saveUploadedImageErr != nil {
		_ = os.Remove(viper.GetString("Upload.Products") + imageGenerateFileName)

		errors.Init()
		errors.Add("images", domain_err.StoreImageOnDiskFailed)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/banners/create")
		return
	}

	req.BannerImage = imageGenerateFileName

	bErr := a.bannerSrv.Create(c, req)
	if bErr != nil {
		//remove images from disk
		_ = os.Remove(viper.GetString("Upload.Categories") + imageGenerateFileName)

		//delete from bucket
		go func() {
			if os.Getenv("STORAGE_STATUS") == "active1" {
				go func() {
					{
						err := a.dep.Storage.DeleteFile(os.Getenv("STORAGE_BANNER_PATH") + imageGenerateFileName)
						if err != nil {
							log.Println("-- failed to delete file from s3 : ", err)
						}
					}
				}()
			}
		}()

		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/banners/create")
		return
	}

	sessions.Set(c, "message", domain_err.SuccessfullyCreated)
	c.Redirect(http.StatusFound, "/admins/banners/create")
	return
}
