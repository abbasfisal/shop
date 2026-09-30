package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	responses "shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/infrastructure/messages"
	"shop/interfaces/http/requests/admin"
	"shop/interfaces/http/response"
	"shop/pkg/errors"
	"shop/pkg/old"
	"shop/pkg/sessions"
	"shop/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// IndexSiteBanner lists the fixed site banners (?placement=header filters).
func (a *AdminHandler) IndexSiteBanner(c *gin.Context) {
	placement := c.Query("placement")
	banners, err := a.siteBannerSrv.Index(c, placement)
	if err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/site-banners")
		return
	}

	response.Render(c, http.StatusOK, "admin_index_site_banner", gin.H{
		"TITLE":      "لیست بنرهای سایت",
		"BANNERS":    responses.ToSiteBanners(banners),
		"FILTER":     placement,
		"MEDIA_PATH": util.GetSiteBannerStoragePath(),
		"CREAT_URL":  "/admins/site-banners/create",
		"INDEX_URL":  "/admins/site-banners",
	})
	return
}

// CreateSiteBanner renders the empty site-banner form.
func (a *AdminHandler) CreateSiteBanner(c *gin.Context) {
	response.Render(c, http.StatusOK, "create_site_banner", gin.H{
		"TITLE":     "ایجاد بنر سایت",
		"PLACEMENT": c.DefaultQuery("placement", "header"),
		"STATUS":    "1",
	})
	return
}

// EditSiteBanner renders the site-banner form pre-filled with the banner.
func (a *AdminHandler) EditSiteBanner(c *gin.Context) {
	banner, err := a.siteBannerSrv.Show(c, uint(util.StringToUint(c.Param("id"))))
	if err != nil {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/site-banners")
		return
	}

	response.Render(c, http.StatusOK, "create_site_banner", gin.H{
		"TITLE":      "ویرایش بنر سایت",
		"BANNER":     responses.ToSiteBanner(banner),
		"IS_EDIT":    true,
		"MEDIA_PATH": util.GetSiteBannerStoragePath(),
	})
	return
}

// siteBannerFormData is the common validation + old-input flash for
// store/update. It returns false when the request must not continue.
func (a *AdminHandler) siteBannerFormData(c *gin.Context) (requests.CreateSiteBannerRequest, bool) {
	_ = c.Request.ParseMultipartForm(32 << 20)

	var req requests.CreateSiteBannerRequest
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		return req, false
	}

	if formErrs := requests.ValidateSiteBannerForm(c.Request.PostForm); len(formErrs) > 0 {
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

// saveSiteBannerUpload stores one uploaded file (field = form input name) and
// returns its stored name. On a validation problem it flashes the error under
// `errKey` and returns "".
func (a *AdminHandler) saveSiteBannerUpload(c *gin.Context, field, errKey string) string {
	// errors.errorList is process-wide: reset it so the caller's
	// errors.Get()[errKey] checks only see what THIS call produced.
	errors.Init()

	file, _ := c.FormFile(field)
	if file == nil {
		return ""
	}

	extension := filepath.Ext(file.Filename)
	if !slices.Contains(util.AllowImageExtensions(), extension) {
		errors.Init()
		errors.Add(errKey, domain_err.MustBeImage)
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		return ""
	}

	name := util.GenerateFilename(file.Filename)

	if os.Getenv("STORAGE_STATUS") == "active" {
		go func() {
			handle, _ := file.Open()
			if err := a.dep.Storage.UploadFile(handle, siteBannerS3Path()+name); err != nil {
				log.Println("-- failed to upload site banner to s3: ", err)
			}
		}()
	}

	if err := c.SaveUploadedFile(file, viper.GetString("Upload.site_banners")+name); err != nil {
		_ = os.Remove(viper.GetString("Upload.site_banners") + name)

		errors.Init()
		errors.Add(errKey, domain_err.StoreImageOnDiskFailed)
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		return ""
	}

	return name
}

// siteBannerS3Path is the S3 key prefix of a site banner upload. Falls back
// to the promotion-banner prefix so a deploy that only configured
// STORAGE_BANNER_PATH keeps working.
func siteBannerS3Path() string {
	if p := os.Getenv("STORAGE_SITE_BANNER_PATH"); p != "" {
		return p
	}
	return os.Getenv("STORAGE_BANNER_PATH")
}

func (a *AdminHandler) StoreSiteBanner(c *gin.Context) {
	req, ok := a.siteBannerFormData(c)
	if !ok {
		c.Redirect(http.StatusFound, "/admins/site-banners/create")
		return
	}

	imageName := a.saveSiteBannerUpload(c, "image", "image")
	if imageName == "" {
		if errors.Get()["image"] == "" {
			errors.Init()
			errors.Add("image", domain_err.IsRequired)
			sessions.Set(c, "errors", errors.ToString())
			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())
		}
		c.Redirect(http.StatusFound, "/admins/site-banners/create")
		return
	}
	req.SiteBannerImage = imageName

	// mobile image is optional — keep whatever arrived ("" → fallback to desktop)
	if mobileName := a.saveSiteBannerUpload(c, "mobile_image", "mobile_image"); mobileName != "" {
		req.MobileImage = mobileName
	} else if errors.Get()["mobile_image"] != "" {
		c.Redirect(http.StatusFound, "/admins/site-banners/create")
		return
	}

	if err := a.siteBannerSrv.Create(c, req); err != nil {
		_ = os.Remove(viper.GetString("Upload.site_banners") + imageName)
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/site-banners/create")
		return
	}

	sessions.Set(c, "message", custom_messages.SuccessfullyCreatedSiteBanner)
	c.Redirect(http.StatusFound, "/admins/site-banners")
	return
}

func (a *AdminHandler) UpdateSiteBanner(c *gin.Context) {
	bannerID := util.StringToUint(c.Param("id"))
	if bannerID == 0 {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/site-banners")
		return
	}

	current, err := a.siteBannerSrv.Show(c, bannerID)
	if err != nil {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/site-banners")
		return
	}

	req, ok := a.siteBannerFormData(c)
	if !ok {
		c.Redirect(http.StatusFound, "/admins/site-banners/"+c.Param("id")+"/edit")
		return
	}

	// keep the current desktop image unless a new file was uploaded
	if imageName := a.saveSiteBannerUpload(c, "image", "image"); imageName != "" {
		req.SiteBannerImage = imageName
	} else if errors.Get()["image"] != "" {
		c.Redirect(http.StatusFound, "/admins/site-banners/"+c.Param("id")+"/edit")
		return
	}

	// mobile image: new file > explicit removal > keep current
	req.MobileImage = current.MobileImage
	if req.RemoveMobileImage {
		req.MobileImage = ""
	} else if mobileName := a.saveSiteBannerUpload(c, "mobile_image", "mobile_image"); mobileName != "" {
		req.MobileImage = mobileName
	} else if errors.Get()["mobile_image"] != "" {
		c.Redirect(http.StatusFound, "/admins/site-banners/"+c.Param("id")+"/edit")
		return
	}

	if err := a.siteBannerSrv.Update(c, bannerID, req); err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/site-banners/"+c.Param("id")+"/edit")
		return
	}

	// drop the replaced files from disk (best effort)
	if req.SiteBannerImage != "" && current.Image != "" && req.SiteBannerImage != current.Image {
		_ = os.Remove(viper.GetString("Upload.site_banners") + current.Image)
	}
	if req.MobileImage != current.MobileImage && current.MobileImage != "" {
		_ = os.Remove(viper.GetString("Upload.site_banners") + current.MobileImage)
	}

	sessions.Set(c, "message", custom_messages.SuccessfullyUpdatedSiteBanner)
	c.Redirect(http.StatusFound, "/admins/site-banners")
	return
}

func (a *AdminHandler) DeleteSiteBanner(c *gin.Context) {
	bannerID := util.StringToUint(c.Param("id"))
	if bannerID == 0 {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/site-banners")
		return
	}

	current, err := a.siteBannerSrv.Show(c, bannerID)
	if err != nil {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/site-banners")
		return
	}

	if delErr := a.siteBannerSrv.Delete(c, bannerID); delErr != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/site-banners")
		return
	}

	if current.Image != "" {
		_ = os.Remove(viper.GetString("Upload.site_banners") + current.Image)
	}
	if current.MobileImage != "" && current.MobileImage != current.Image {
		_ = os.Remove(viper.GetString("Upload.site_banners") + current.MobileImage)
	}

	sessions.Set(c, "message", custom_messages.DeleteSuccessfully)
	c.Redirect(http.StatusFound, "/admins/site-banners")
	return
}
