package mapper 
 
import (
	domain "g6/blog-api/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenModel struct {
	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`
}

type UserModel struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Username       string             `bson:"username"`
	Email          string             `bson:"email"`
	FirstName      string             `bson:"first_name"`
	LastName       string             `bson:"last_name"`
	Password       string             `bson:"password"`
	Role           string             `bson:"role"`
	Tokens         TokenModel         `bson:"tokens"`
	Bio            string             `bson:"bio"`	
	ProfilePicture string             `bson:"profile_picture"`
	CreatedAt      time.Time          `bson:"created_at"`	
	UpdatedAt      time.Time          `bson:"updated_at"`
}

// Convert to domain
func UserToDomain(user *UserModel) *domain.User {
	return &domain.User{
		ID:             user.ID.Hex(),
		Username:       user.Username,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Password:       user.Password,
		Role:           domain.Role(user.Role),
		Tokens:         domain.TokenPair{
			AccessToken: user.Tokens.AccessToken, 
			RefreshToken: user.Tokens.RefreshToken,
		},
		Bio:            user.Bio,
		ProfilePicture: user.ProfilePicture,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func UserFromDomain(user *domain.User) (*UserModel, error) {
	var objectID primitive.ObjectID
	var err error
	// if user.ID is not empty, convert it to ObjectID
	if user.ID != "" {
		objectID, err = primitive.ObjectIDFromHex(user.ID)
		if err != nil {
			return nil, err
		}
	} else{
		objectID = primitive.NewObjectID() // generate new ObjectID if ID is empty
	}

	return &UserModel{
		ID:             objectID,
		Username:       user.Username,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Password:       user.Password,
		Role:           string(user.Role),
		Tokens: TokenModel{
			AccessToken:  user.Tokens.AccessToken,
			RefreshToken: user.Tokens.RefreshToken,
		},
		Bio:            user.Bio,
		ProfilePicture: user.ProfilePicture,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}
