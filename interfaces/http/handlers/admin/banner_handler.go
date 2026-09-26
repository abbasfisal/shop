package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	responses "shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/infrastructure/messages"
	"shop/interfaces/http/requests/admin"
	"shop/interfaces/http/response"
	"shop/pkg/errors"
	"shop/pkg/old"
	"shop/pkg/sessions"
	"shop/pkg/util"
)

// IndexBanner lists every promotion banner (sidebar: لیست بنر ها).
func (a *AdminHandler) IndexBanner(c *gin.Context) {
	banners, err := a.bannerSrv.Index(c)
	if err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/banners")
		return
	}

	response.Render(c, http.StatusOK, "admin_index_banner", gin.H{
		"TITLE":      "لیست بنرهای پروموشن",
		"BANNERS":    responses.ToBanners(banners),
		"MEDIA_PATH": util.GetBannerStoragePath(),
	})
	return
}

// CreateBanner renders the empty promotion form.
func (a *AdminHandler) CreateBanner(c *gin.Context) {
	response.Render(c, http.StatusOK, "create_banner", gin.H{
		"TITLE":  "ایجاد بنر پروموشن",
		"LAYOUT": "two",
		"STATUS": "1",
	})
	return
}

// EditBanner renders the promotion form pre-filled with the banner.
func (a *AdminHandler) EditBanner(c *gin.Context) {
	banner, err := a.bannerSrv.Show(c, uint(util.StringToUint(c.Param("id"))))
	if err != nil {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/banners")
		return
	}

	response.Render(c, http.StatusOK, "create_banner", gin.H{
		"TITLE":      "ویرایش بنر پروموشن",
		"BANNER":     responses.ToBanner(banner),
		"IS_EDIT":    true,
		"MEDIA_PATH": util.GetBannerStoragePath(),
	})
	return
}

// bannerFormData is the common validation + old-input flash for store/update.
// It returns false when the request must not continue.
func (a *AdminHandler) bannerFormData(c *gin.Context) (requests.CreateBannerRequest, bool) {
	_ = c.Request.ParseMultipartForm(32 << 20)

	var req requests.CreateBannerRequest
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		return req, false
	}

	if formErrs := requests.ValidateBannerForm(c.Request.PostForm); len(formErrs) > 0 {
		errors.Init()
		for key, msg := range formErrs {
			errors.Add(key, msg)
		}
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		return req, false
	}

	return req, true
}

// saveBannerImage stores the uploaded file and returns its stored name.
// On a validation problem it flashes the error and returns an empty name.
func (a *AdminHandler) saveBannerImage(c *gin.Context) string {
	imageFile, _ := c.FormFile("image")
	if imageFile == nil {
		return ""
	}

	extension := filepath.Ext(imageFile.Filename)
	if !slices.Contains(util.AllowImageExtensions(), extension) {
		errors.Init()
		errors.Add("image", domain_err.MustBeImage)
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		return ""
	}

	imageName := util.GenerateFilename(imageFile.Filename)

	if os.Getenv("STORAGE_STATUS") == "active" {
		go func() {
			imageHandle, _ := imageFile.Open()
			if err := a.dep.Storage.UploadFile(imageHandle, os.Getenv("STORAGE_BANNER_PATH")+imageName); err != nil {
				log.Println("-- failed to upload banner to s3: ", err)
			}
		}()
	}

	if err := c.SaveUploadedFile(imageFile, viper.GetString("Upload.Banners")+imageName); err != nil {
		_ = os.Remove(viper.GetString("Upload.Banners") + imageName)

		errors.Init()
		errors.Add("image", domain_err.StoreImageOnDiskFailed)
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		return ""
	}

	return imageName
}

func (a *AdminHandler) StoreBanner(c *gin.Context) {
	req, ok := a.bannerFormData(c)
	if !ok {
		c.Redirect(http.StatusFound, "/admins/banners/create")
		return
	}

	imageName := a.saveBannerImage(c)
	if imageName == "" {
		if errors.Get()["image"] == "" {
			errors.Init()
			errors.Add("image", domain_err.IsRequired)
			sessions.Set(c, "errors", errors.ToString())
			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())
		}
		c.Redirect(http.StatusFound, "/admins/banners/create")
		return
	}
	req.BannerImage = imageName

	if err := a.bannerSrv.Create(c, req); err != nil {
		_ = os.Remove(viper.GetString("Upload.Banners") + imageName)
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/banners/create")
		return
	}

	sessions.Set(c, "message", custom_messages.SuccessfullyCreatedBanner)
	c.Redirect(http.StatusFound, "/admins/banners")
	return
}

func (a *AdminHandler) UpdateBanner(c *gin.Context) {
	bannerID := util.StringToUint(c.Param("id"))
	if bannerID == 0 {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/banners")
		return
	}

	current, err := a.bannerSrv.Show(c, bannerID)
	if err != nil {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/banners")
		return
	}

	req, ok := a.bannerFormData(c)
	if !ok {
		c.Redirect(http.StatusFound, "/admins/banners/"+c.Param("id")+"/edit")
		return
	}

	// keep the current image unless a new file was uploaded
	if imageName := a.saveBannerImage(c); imageName != "" {
		req.BannerImage = imageName
	} else if errors.Get()["image"] != "" {
		c.Redirect(http.StatusFound, "/admins/banners/"+c.Param("id")+"/edit")
		return
	}

	if err := a.bannerSrv.Update(c, bannerID, req); err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/banners/"+c.Param("id")+"/edit")
		return
	}

	// drop the replaced file from disk (best effort)
	if req.BannerImage != "" && current.Image != "" && req.BannerImage != current.Image {
		_ = os.Remove(viper.GetString("Upload.Banners") + current.Image)
	}

	sessions.Set(c, "message", custom_messages.SuccessfullyUpdatedBanner)
	c.Redirect(http.StatusFound, "/admins/banners")
	return
}

func (a *AdminHandler) DeleteBanner(c *gin.Context) {
	bannerID := util.StringToUint(c.Param("id"))
	if bannerID == 0 {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/banners")
		return
	}

	current, err := a.bannerSrv.Show(c, bannerID)
	if err == nil {
		if delErr := a.bannerSrv.Delete(c, bannerID); delErr != nil {
			sessions.Set(c, "message", domain_err.SomethingWrongHappened)
			c.Redirect(http.StatusFound, "/admins/banners")
			return
		}
		if current.Image != "" {
			_ = os.Remove(viper.GetString("Upload.Banners") + current.Image)
		}
	} else {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/banners")
		return
	}

	sessions.Set(c, "message", custom_messages.DeleteSuccessfully)
	c.Redirect(http.StatusFound, "/admins/banners")
	return
}
