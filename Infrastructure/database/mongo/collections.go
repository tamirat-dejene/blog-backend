package mongo

// Collections holds the names of MongoDB collections used in the application.
// User of this struct can pass any collection names they want to use, leaving the
// unnecessary collections empty.
type Collections struct {
	BlogPosts         string
	BlogComments      string
	BlogUserReactions string

	Users         string
	RefreshTokens string
	PasswordReset string

	AIBlogPosts string
}

// func NewCollections(blogPosts, blogComments, blogUserReactions string) *collections {
// 	return &collections{
// 		BlogPosts:         blogPosts,
// 		BlogComments:      blogComments,
// 		BlogUserReactions: blogUserReactions,
// 	}
// }
