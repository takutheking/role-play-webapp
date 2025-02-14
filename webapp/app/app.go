package app

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mittz/role-play-webapp/webapp/database"
)

var dbHandler database.DatabaseHandler

// getUserID always returns the same id which is for scstore user.
// This should be updated when this app implements an authentication feature.
func getUserID() uint {
	return uint(2)
}

func getIndexEndpoint(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/products")
}

func postInitEndpoint(c *gin.Context) {
	if err := dbHandler.InitDatabase(); err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.String(http.StatusAccepted, "Initiatilized data.")
}

func getCheckoutsEndpoint(c *gin.Context) {
	userID := getUserID()
	checkouts, err := dbHandler.GetCheckouts(userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.HTML(http.StatusOK, "checkouts.html", gin.H{
		"title":     "Checkouts",
		"checkouts": checkouts,
	})
}

func postCheckoutEndpoint(c *gin.Context) {
	userID := getUserID()

	productID, err := strconv.Atoi(c.PostForm("product_id"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	productQuantity, err := strconv.Atoi(c.PostForm("product_quantity"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	checkoutID, err := dbHandler.CreateCheckout(uint(userID), uint(productID), uint(productQuantity))
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	checkout, err := dbHandler.GetCheckout(checkoutID)
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.HTML(http.StatusAccepted, "checkout.html", gin.H{
		"title":    "Checkout",
		"checkout": checkout,
	})
}

func getProductEndpoint(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	product, err := dbHandler.GetProduct(uint(productID))
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.HTML(http.StatusOK, "product.html", gin.H{
		"title":   "Product",
		"product": product,
	})
}

func getProductsEndpoint(c *gin.Context) {
	products, err := dbHandler.GetProducts()
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":    "Products",
		"products": products,
	})
}

func SetupRouter(dbh database.DatabaseHandler, assetsDir string, templatesDirMatch string) *gin.Engine {
	dbHandler = dbh

	router := gin.Default()
	router.Static("/assets", assetsDir)
	router.StaticFile("/favicon.ico", filepath.Join(assetsDir, "favicon.ico"))
	router.LoadHTMLGlob(templatesDirMatch)

	router.GET("/", getIndexEndpoint)

	router.POST("/admin/init", postInitEndpoint)

	router.GET("/product/:product_id", getProductEndpoint)
	router.GET("/products", getProductsEndpoint)
	router.GET("/checkouts", getCheckoutsEndpoint)
	router.POST("/checkout", postCheckoutEndpoint)

	return router
}
