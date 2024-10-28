package server

import (
	"fmt"
	"forum-go/internal/models"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := s.db.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, r, "../posts", map[string]interface{}{"Posts": posts})
}

func (s *Server) PostNewPostsHandler(w http.ResponseWriter, r *http.Request) {
	newPost := models.Post{
		PostId:                strconv.Itoa(rand.Intn(math.MaxInt32)),
		Title:                 r.FormValue("title"),
		Content:               r.FormValue("content"),
		UserID:                r.FormValue("UserId"),
		CreationDate:          time.Now(),
		FormattedCreationDate: time.Now().Format("Jan 02, 2006 - 15:04:05"),
	}
	err := s.db.AddPost(newPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func (s *Server) DeletePostsHandler(w http.ResponseWriter, r *http.Request) {
	PostID := r.FormValue("postId")
	fmt.Println(PostID)
	err := s.db.DeletePost(PostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
func (s *Server) GetNewPostHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	render(w, r, "createPost", nil)
}
func (s *Server) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := strings.Split(r.URL.Path, "/")
	if len(vars) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	postID := vars[2]
	post, err := s.db.GetPost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if post.PostId == "" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	render(w, r, "detailsPost", map[string]interface{}{"Post": post})
}
func (s *Server) PostCommentHandler(w http.ResponseWriter, r *http.Request) {
}

func IsUniquePost(posts []models.Post, post string) bool {
	for _, existingPost := range posts {
		if strings.EqualFold(existingPost.PostId, post) {
			return false
		}
	}
	return true
}

const MaxChar = 1000

func (s *Server) ValidatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if len(post.Content) > MaxChar {
		render(w, r, "../posts", map[string]interface{}{"Posts": s.posts, "Error": "Post content is too long"})
		return
	}
	if len(post.Categories) < 1 {
		//return (Wainting for cat selector)
	}
}
