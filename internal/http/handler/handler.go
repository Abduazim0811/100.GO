package handler

import (
	"log"
	"net/http"

	_ "100.GO/docs"

	"100.GO/internal/entity/origin"
	"100.GO/internal/entity/user"
	"100.GO/internal/infrastructura/repository/redis"
	pkg "100.GO/internal/pkg/email"
	"100.GO/internal/pkg/token"
	"100.GO/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

// @title Artisan Connect
// @version 1.0
// @description This is a sample server for a restaurant reservation system.
// @securityDefinitions.apikey Bearer
// @in 				header
// @name Authorization
// @description Enter the token in the format `Bearer {token}`
// @host localhost:8888
// @BasePath /

type UserHandler struct {
	service *service.UserService
	redis   *redis.RedisClient
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body user.CreateUser true "User request body"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /register [post]
func (u *UserHandler) CreateUser(c *gin.Context) {
	var req user.CreateUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Hashedpassword error:", err)
	}
	hashedpassword := string(bytes)
	code := 10000 + rand.Intn(90000)
	err = pkg.SendEmail(req.Email, pkg.SendClientCode(code, req.Firstname))
	if err != nil {
		log.Println("ERROR: sending email to user !!", err)
	}

	userdata := map[string]interface{}{
		"firstname": req.Firstname,
		"lastname":  req.Lastname,
		"email":     req.Email,
		"password":  hashedpassword,
	}

	if u.redis == nil {
		log.Println("Redis client is not initialized")
		return
	}

	err = u.redis.SetHash(req.Email, userdata)
	if err != nil {
		log.Printf("Failed to save user data in Redis: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "User created"})
}

// VerifyCode godoc
// @Summary      Verify the user code
// @Description  Verify the code sent to the user's email
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        code  body user.VerifyCode true "Verification code and email"
// @Success      200   {object} string
// @Failure      400   {object} string
// @Failure      500   {object} string
// @Router       /verify [post]
func (u *UserHandler) VerifyCode(c *gin.Context) {
	var req user.VerifyCode

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}
	user, err := u.redis.VerifyEmail(req.Email, int64(req.Code))
	if err != nil {
		log.Println("error code")
		return
	}

	err = u.service.Createuser(user)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Verify Successfully"})
}

// Login godoc
// @Summary      Login a user
// @Description  Login a user with email and password
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        credentials  body user.Login true "User credentials"
// @Success      200          {object} map[string]string
// @Failure      400          {object} string
// @Failure      500          {object} string
// @Router       /login [post]
func (u *UserHandler) Login(c *gin.Context) {
	var req user.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	user, err := u.service.GetuserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})
		return
	}

	token, err := token.GenerateJWTToken(user.Email)
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"Token": token})
}

/*
 Origin Handlers
*/

// CreateOrigin godoc
// @Summary      Create a new origin
// @Description  Create a new origin entity
// @Tags         origin
// @Accept       json
// @Produce      json
// @Param        origin  body origin.CreateOrigin true "Origin request body"
// @Success      200     {object} string
// @Failure      400     {object} string
// @Failure      500     {object} string
// @Router       /origins [post]
func (u *UserHandler) CreateOrigin(c *gin.Context) {
	var req origin.CreateOrigin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	err := u.service.Createorigin(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Origin created"})
}

// GetbyIdOrigin godoc
// @Summary      Get origin by ID
// @Description  Retrieve a specific origin entity by its ID
// @Tags         origin
// @Produce      json
// @Param        id  path string true "Origin ID"
// @Success      200  {object} origin.GetOrigin
// @Failure      400  {object} string
// @Failure      500  {object} string
// @Router       /origins/{id} [get]
func (u *UserHandler) GetbyIdOrigin(c *gin.Context) {
	id := c.Param("id")

	origin, err := u.service.GetoriginById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return
	}
	c.JSON(http.StatusOK, origin)
}

// GetOrigin godoc
// @Summary      Get all origins
// @Description  Retrieve a list of all origin entities
// @Tags         origin
// @Produce      json
// @Success      200  {array} origin.GetOrigin
// @Failure      400  {object} string
// @Failure      500  {object} string
// @Router       /origins [get]
func (u *UserHandler) GetOrigin(c *gin.Context) {
	origins, err := u.service.GetAllorigins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, origins)
}

// UpdateOrigin godoc
// @Summary      Update an existing origin
// @Description  Update the details of an existing origin
// @Tags         origin
// @Accept       json
// @Produce      json
// @Param        id      path string true "Origin ID"
// @Param        origin  body origin.CreateOrigin true "Updated origin data"
// @Success      200     {object} string
// @Failure      400     {object} string
// @Failure      500     {object} string
// @Router       /origins/{id} [put]
func (u *UserHandler) UpdateOrigin(c *gin.Context) {
	id := c.Param("id")

	reqid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}
	var req origin.CreateOrigin
	err = u.service.Updateorigin(reqid, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Origin updated"})
}

// DeleteOrigin godoc
// @Summary      Delete an origin
// @Description  Delete an origin entity by its ID
// @Tags         origin
// @Param        id  path string true "Origin ID"
// @Success      200  {object} string
// @Failure      400  {object} string
// @Failure      500  {object} string
// @Router       /origins/{id} [delete]
func (u *UserHandler) DeleteOrigin(c *gin.Context) {
	id := c.Param("id")

	reqid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}

	err = u.service.Deleteorigin(reqid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Origin deleted"})
}

func (u *UserHandler) EnableCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		res, err := u.service.OriginGetall()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			for _, v := range res {
				if origin == v.Origin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		c.Next()
	}
}

func (u *UserHandler) CorsMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "Congrats buddy, you are now one of the trusted persons",
	})
}
