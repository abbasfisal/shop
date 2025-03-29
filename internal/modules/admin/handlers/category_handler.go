package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
	"shop/internal/modules/admin/requests"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/custom_messages"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/old"
	"shop/internal/pkg/sessions"
	"shop/internal/pkg/util"
	"slices"
	"strconv"
)

//-------------------------------
//		category routes
//-------------------------------

func (a *AdminHandler) IndexCategory(c *gin.Context) {

	categories, err := a.categorySrv.Index(c)

	if err.Code == 400 {
		c.JSON(200, gin.H{
			"data": "empty",
		})
	} else if err.Code == 500 {
		html.Error500(c)
		return
	}
	html.Render(c, 200, "modules/admin/html/admin_index_category", gin.H{
		"TITLE":      "لیست دسته بندی",
		"CATEGORIES": categories,
	})
	return
}

func (a *AdminHandler) CreateCategory(c *gin.Context) {
	categories, _ := a.categorySrv.GetAllCategories(c)
	html.Render(c, http.StatusFound, "modules/admin/html/admin_create_category",
		gin.H{
			"TITLE":      "ایجاد کتگوری",
			"CATEGORIES": categories,
		})
	return
}

func (a *AdminHandler) StoreCategory(c *gin.Context) {
	var req requests.CreateCategoryRequest
	_ = c.Request.ParseForm()
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/categories/create")
		return
	}

	//slug unique validation
	if ok := a.categorySrv.CheckSlugUniqueness(context.TODO(), req.Slug); ok {
		errors.Init()
		errors.Add("slug", custom_error.MustBeUnique)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/categories/create")
		return
	}

	//------------------------
	//	upload and save image
	//------------------------
	imageFile, imageErr := c.FormFile("image")
	//require validation on image
	//if imageErr != nil {
	//	//validation on required tag
	//	errors.Init()
	//	errors.Add("image", custom_error.IsRequired)
	//	sessions.Set(c, "errors", errors.ToString())
	//
	//	old.Init()
	//	old.Set(c)
	//	sessions.Set(c, "olds", old.ToString())
	//
	//	c.Redirect(http.StatusFound, "/admins/categories/create")
	//
	//	return
	//}

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

			c.Redirect(http.StatusFound, "/admins/categories/create")
			return
		}

		//upload image and store on disk
		newImageName := util.GenerateFilename(imageFile.Filename)
		pathToUpload := viper.GetString("Upload.Categories") + newImageName
		uploadErr := c.SaveUploadedFile(imageFile, pathToUpload)
		if uploadErr != nil {
			fmt.Println("upload error:", uploadErr)
			errors.Init()
			errors.Add("image", custom_error.UploadImageError)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, "/admins/categories/create")

			return
		}
		req.Image = newImageName
	}

	//---------------------------
	//	end upload and save image
	//---------------------------

	newCategory, err := a.categorySrv.Create(c.Request.Context(), &req)
	if err != nil || newCategory.ID <= 0 {
		_ = os.Remove(pathToUpload)
		sessions.Set(c, "message", "خطا در ایجاد دسته بندی")
		c.Redirect(http.StatusFound, "/admins/categories/create")
		return
	}

	sessions.Set(c, "message", "ایجاد دسته بندی با موفقیت انجام شد")
	c.Redirect(http.StatusFound, "/admins/categories/create")
	return
}

func (a *AdminHandler) ShowCategory(c *gin.Context) {

	catID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		catID = 1
	}
	cat, catErr := a.categorySrv.Show(context.TODO(), catID)
	if catErr.Code == 404 {
		html.Render(c, http.StatusOK, "modules/admin/html/admin_index_category", gin.H{
			"MESSAGE": custom_error.RecordNotFound,
		})
		return
	}
	if catErr.Code == 500 {
		html.Error500(c)
		return
	}

	html.Render(c, http.StatusFound, "modules/admin/html/admin_show_category", gin.H{
		"TITLE":    "show a category",
		"CATEGORY": cat,
	})
	return
}

func (a *AdminHandler) EditCategory(c *gin.Context) {
	catID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}

	cat, cErr := a.categorySrv.Show(c.Request.Context(), catID)
	categories, _ := a.categorySrv.GetAllCategories(c)

	if cErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}
	if cErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}

	html.Render(c, http.StatusFound, "admin_edit_category",
		gin.H{
			"TITLE":      "ویرایش دسته بندی",
			"CATEGORY":   cat,
			"CATEGORIES": categories,
		})
	return
}

func (a *AdminHandler) UpdateCategory(c *gin.Context) {

	catID, catIDErr := strconv.Atoi(c.Param("id"))
	if catIDErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}

	url := fmt.Sprintf("/admins/categories/%d/edit", catID)

	var req requests.UpdateCategoryRequest
	_ = c.Request.ParseForm()
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, url)
		return
	}

	categoryData, catErr := a.categorySrv.Show(c, catID)
	if catErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}
	if catErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}

	//slug unique validation
	if categoryData.Slug != req.Slug {
		if ok := a.categorySrv.CheckSlugUniqueness(context.TODO(), req.Slug); ok {
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

	//------------------------
	//	upload and save image
	//------------------------
	imageFile, imageErr := c.FormFile("image")

	pathToUpload := ""
	if imageErr != nil {
		req.Image = categoryData.Image
	}
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
		pathToUpload = viper.GetString("Upload.Categories") + newImageName
		uploadErr := c.SaveUploadedFile(imageFile, pathToUpload)
		if uploadErr != nil {
			fmt.Println("upload error:", uploadErr)
			errors.Init()
			errors.Add("image", custom_error.UploadImageError)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)

			return
		}
		req.Image = newImageName
	}

	//---------------------------
	//	end upload and save image
	//---------------------------

	categoryUpdateErr := a.categorySrv.Edit(c, catID, &req)

	if categoryUpdateErr.Code > 0 {
		if categoryUpdateErr.Code == 404 {
			sessions.Set(c, "message", custom_error.RecordNotFound)
			c.Redirect(http.StatusFound, "/admins/categories/")
			return
		}
		if categoryUpdateErr.Code == 500 {
			sessions.Set(c, "message", custom_error.InternalServerError)
			c.Redirect(http.StatusFound, "/admins/categories/")
			return
		}
	}

	if pathToUpload != "" {
		_ = os.Remove(viper.GetString("upload.categories") + categoryData.Image)
	}
	sessions.Set(c, "message", custom_messages.CategoryUpdatedSucc)
	c.Redirect(http.StatusFound, "/admins/categories/")
	return

}

func (a *AdminHandler) CategoryProducts(c *gin.Context) {
	c.JSON(200,
		gin.H{
			"category_id": c.Param("id"),
			"msg":         "implement me",
		})
	return
}
