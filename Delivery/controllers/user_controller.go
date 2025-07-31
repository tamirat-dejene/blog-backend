package controllers
import(
	"net/http"
	"g6/blog-api/Delivery/dto"
	usecases "g6/blog-api/Usecases"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	uc usecases.UserUsecase
}

func NewUserController(uc usecases.UserUsecase) *UserController {
	return &UserController{uc}
}

func (ctrl *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := ctrl.uc.Register(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	
	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}	

//Login
//Logout
func (ctrl *UserController) Logout(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if err := ctrl.uc.Logout(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

func (ctrl *UserController) ChangeRole(ctx *gin.Context) {
	var req dto.ChangeRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	initiator := ctx.GetString("role")
	target := ctx.Param("userID")
	if err := ctrl.uc.ChangeRole(initiator, target, req); err != nil{
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User role changed successfully"})
}