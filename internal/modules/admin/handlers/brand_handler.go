package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/old"
	"shop/internal/pkg/sessions"
	"shop/internal/pkg/util"
	"slices"
	"strconv"
	"strings"
)

func (a *AdminHandler) ShowCreateBrand(c *gin.Context) {
	html.Render(c, 200, "create-brand", gin.H{"TITLE": "create new brand"})
	return
}

func (a *AdminHandler) IndexBrand(c *gin.Context) {
	//categories, err := a.categorySrv.Index(c)
	brands, err := a.brandSrv.Index(c)

	if err.Code == 404 {
		sessions.Set(c, "message", "برندی موجود نمی باشد")
		c.Redirect(http.StatusFound, "/admins/brands/")
		return

	} else if err.Code == 500 {
		html.Error500(c)
		return
	}

	html.Render(c, 200, "modules/admin/html/admin_index_brand", gin.H{
		"TITLE":  "Index Brand",
		"BRANDS": brands,
	})
	return
}

func (a *AdminHandler) ShowBrand(c *gin.Context) {
	brandID, brand, done := checkIDAndExistence(c, a)
	if done {
		return
	}

	html.Render(c, http.StatusOK, "show-brand", gin.H{
		"TITLE":    "show brand",
		"BRAND_ID": brandID,
		"BRAND":    brand,
	})
	return
}

func checkIDAndExistence(c *gin.Context, a *AdminHandler) (int, *responses.Brand, bool) {
	brandID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", "ID برند صحیح نمی باشد.")
		c.Redirect(http.StatusFound, "/admins/brands")
	}
	brand, bErr := a.brandSrv.Show(c, brandID)
	if bErr.Code == 404 {
		sessions.Set(c, "message", "برند موجود نمی باشد.")
		c.Redirect(http.StatusFound, "/admins/brands")
		return 0, nil, true
	}
	if bErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/brands")
		return 0, nil, true
	}
	return brandID, brand, false
}

func (a *AdminHandler) EditBrand(c *gin.Context) {
	brandID, brand, done := checkIDAndExistence(c, a)
	if done {
		return
	}

	html.Render(c, http.StatusFound, "edit-brand", gin.H{
		"TITLE":    "update brand",
		"BRAND_ID": brandID,
		"BRAND":    brand,
	})
	return
}

func (a *AdminHandler) UpdateBrand(c *gin.Context) {
	brandID, brand, done := checkIDAndExistence(c, a)
	if done {
		return
	}

	url := fmt.Sprintf("/admins/brands/%d/edit", brandID)

	var req requests.UpdateBrandRequest
	_ = c.Request.ParseForm()
	req.Slug = strings.TrimSpace(req.Slug)
	req.Title = strings.TrimSpace(req.Title)
	req.Image = strings.TrimSpace(req.Image)
	err := c.ShouldBind(&req)
	if err != nil {
		sessions.Set(c, "message", "خطایی رخ داده است")
		c.Redirect(http.StatusFound, "/admins/brands")
		return
	}

	if req.Slug != brand.Slug {
		if a.brandSrv.CheckSlugUniqueness(c, req.Slug) {
			//slug is not unique

			errors.Init()
			errors.Add("slug", custom_error.MustBeUnique)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}
	}

	imageFile, imageErr := c.FormFile("image")
	oldImagePath := brand.Image
	pathToUpload := brand.Image

	if imageErr == nil {

		// file extension validation
		fileExtension := filepath.Ext(imageFile.Filename)
		ok := slices.Contains(util.AllowImageExtensions(), fileExtension)
		if !ok {
			errors.Init()
			errors.Add("image", custom_error.MustBeImage)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}

		//upload image and store on disk
		newImageName := util.GenerateFilename(imageFile.Filename)
		pathToUpload = viper.GetString("Upload.Brands") + newImageName

	}

	if req.Title == "" {
		req.Title = brand.Title
	}
	if req.Slug == "" {
		req.Slug = brand.Slug
	}
	req.Image = pathToUpload

	_, uErr := a.brandSrv.Update(c, brandID, &req)
	if uErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, url)
		return
	}
	if uErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, url)
		return
	}

	if imageErr == nil {
		uploadErr := c.SaveUploadedFile(imageFile, pathToUpload)

		if uploadErr != nil {
			errors.Init()
			errors.Add("image", custom_error.UploadImageError)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}
	}

	_ = os.Remove(oldImagePath)

	sessions.Set(c, "message", "با موفقیت ویرایش گردید")
	c.Redirect(http.StatusFound, url)
	return
}

func (a *AdminHandler) StoreBrand(c *gin.Context) {

	var req requests.CreateBrandRequest
	_ = c.Request.ParseForm()
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("error ---------- :", err)
		errors.Init()
		errors.SetFromErrors(err)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/brands/create")
		return
	}

	//slug unique validation
	if ok := a.brandSrv.CheckSlugUniqueness(c.Request.Context(), req.Slug); ok {
		errors.Init()
		errors.Add("slug", custom_error.MustBeUnique)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/brands/create")
		return
	}

	//------------------------
	//	upload and save image
	//------------------------
	imageFile, imageErr := c.FormFile("image")
	//require validation on image
	if imageErr != nil {
		//validation on required tag
		errors.Init()
		errors.Add("image", custom_error.IsRequired)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/brands/create")

		return
	}

	pathToUpload := ""
	if imageErr == nil {
		// file extension validation
		fileExtension := filepath.Ext(imageFile.Filename)
		ok := slices.Contains(util.AllowImageExtensions(), fileExtension)
		if !ok {
			errors.Init()
			errors.Add("image", custom_error.MustBeImage)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, "/admins/brands/create")
			return
		}

		//upload image and store on disk
		newImageName := util.GenerateFilename(imageFile.Filename)
		pathToUpload := viper.GetString("Upload.brands") + newImageName
		uploadErr := c.SaveUploadedFile(imageFile, pathToUpload)
		if uploadErr != nil {
			fmt.Println("upload error:", uploadErr)
			errors.Init()
			errors.Add("image", custom_error.UploadImageError)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, "/admins/brands/create")

			return
		}
		req.Image = newImageName
	}

	//---------------------------
	//	end upload and save image
	//---------------------------

	newBrand, err := a.brandSrv.Create(c.Request.Context(), &req)
	if err != nil || newBrand.ID <= 0 {
		_ = os.Remove(pathToUpload)
		fmt.Println("----- error in creating brand ---- : ", err)
		sessions.Set(c, "message", "خطا در ایجاد برند")
		c.Redirect(http.StatusFound, "/admins/categories/create")
		return
	}

	sessions.Set(c, "message", "ایجاد برند با موفقیت انجام شد")
	c.Redirect(http.StatusFound, "/admins/brands/create")
	return

}
