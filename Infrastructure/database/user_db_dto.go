package database

// import (
// 	domain "g6/blog-api/Domain"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type UserDB struct {
// 	ID        primitive.ObjectID `bson:"_id"`
// 	Username  string             `bson:"username"`
// 	Email     string             `bson:"email"`
// 	Password  string             `bson:"password,omitempty"`
// 	FirstName string             `bson:"first_name"`
// 	LastName  string             `bson:"last_name"`
// 	Role      string             `bson:"role"`
// 	Bio       string             `bson:"bio"`
// 	AvatarURL string             `bson:"avatar_url"`
// 	CreatedAt primitive.DateTime `bson:"created_at"`
// 	UpdatedAt primitive.DateTime `bson:"updated_at"`
// }

// func FromUserEntityToDB(user *domain.User) *UserDB {
// 	return &UserDB{
// 		ID:        user.ID,
// 		Username:  user.Username,
// 		Email:     user.Email,
// 		Password:  user.Password, // Password should be hashed before storing
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 		Role:      user.Role,
// 		Bio:       user.Bio,
// 		AvatarURL: user.AvatarURL,
// 		CreatedAt: primitive.NewDateTimeFromTime(user.CreatedAt),
// 		UpdatedAt: primitive.NewDateTimeFromTime(user.UpdatedAt),
// 	}
// }

// func FromUserDBToEntity(userDB *UserDB) *domain.User {
// 	return &domain.User{
// 		ID:        userDB.ID,
// 		Username:  userDB.Username,
// 		Email:     userDB.Email,
// 		Password:  userDB.Password, // Password should be handled securely
// 		FirstName: userDB.FirstName,
// 		LastName:  userDB.LastName,
// 		Role:      userDB.Role,
// 		Bio:       userDB.Bio,
// 		AvatarURL: userDB.AvatarURL,
// 		CreatedAt: userDB.CreatedAt.Time(),
// 		UpdatedAt: userDB.UpdatedAt.Time(),
// 	}
// }
// func FromUserDBListToEntityList(userDBs []*UserDB) []*domain.User {
// 	var users []*domain.User
// 	for _, userDB := range userDBs {
// 		users = append(users, FromUserDBToEntity(userDB))
// 	}
// 	return users
// }
