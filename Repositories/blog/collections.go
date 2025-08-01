package repository

type collections struct {
	BlogPosts         string
	BlogComments      string
	BlogUserReactions string
}

func NewCollections(blogPosts, blogComments, blogUserReactions string) *collections {
	return &collections{
		BlogPosts:         blogPosts,
		BlogComments:      blogComments,
		BlogUserReactions: blogUserReactions,
	}
}
