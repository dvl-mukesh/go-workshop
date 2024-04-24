package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Digivate-Labs-Pvt-Ltd/dvlutil"
	"github.com/dvl-mukesh/go-workshop/internal/comment"
)

var (
	MsgInvalidId         = "Invalid ID"
	MsgBadReq            = "Bad Request"
	MsgInternalServerErr = "Internal Server Error"
	MsgCreateSuccess     = "Comment Created Successfully"
	MsgFetchSuccess      = "Comment Fetched Successfully"
	MsgDelteSuccess      = "Comment Deleted Successfully"
	MsgUpdateSuccess     = "Comment Updated Successfully"
)

type Handler struct {
	Router  *http.ServeMux
	Service *comment.Service
	http.Server
}

func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) SetupRoutes() {
	log.Println("Setting up routes")

	h.Router = http.NewServeMux()
	h.Router.HandleFunc("/api/health", healthHandler)
	h.Router.HandleFunc("GET /api/comment", h.GetAllComments)
	h.Router.HandleFunc("GET /api/comment/{id}", h.GetComment)
	h.Router.HandleFunc("POST /api/comment", h.PostComment)
	h.Router.HandleFunc("PUT /api/comment/{id}", h.PutComment)
	// h.Router.HandleFunc("PATCH /api/comment/{id}", h.PatchComment)
	h.Router.HandleFunc("DELETE /api/comment/{id}", h.DeleteComment)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", h.Router))
	h.Router = v1
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I am alive!")
}

func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	i, err := strconv.ParseUint(id, 10, 64)

	if err != nil {

		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgInvalidId,
		})
		log.Println(err)
		return
	}

	comments, err := h.Service.GetComment(uint(i))

	if err != nil {
		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgInternalServerErr,
		})
		log.Println(err)
		return
	}

	dvlutil.WriteJSON(w, http.StatusOK, dvlutil.Response{
		Status: dvlutil.StatusCodeOK,
		Msg:    MsgFetchSuccess,
		Data:   comments,
	})

}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var comment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgBadReq,
		})
		log.Println(err)
		return
	}

	newComment, err := h.Service.PostComment(comment)

	if err != nil {
		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgBadReq,
		})
		log.Println(err)
		return
	}

	dvlutil.WriteJSON(w, http.StatusOK, dvlutil.Response{
		Status: dvlutil.StatusCodeOK,
		Msg:    MsgCreateSuccess,
		Data:   newComment,
	})

}

func (h *Handler) PutComment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var comment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgBadReq,
		})
		log.Println(err)
		return
	}

	id := r.PathValue("id")
	i, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgInvalidId,
		})
		log.Println(err)
		return
	}

	newComment, err := h.Service.UpdateComment(uint(i), comment)

	if err != nil {
		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgBadReq,
		})
		log.Println(err)
		return
	}

	dvlutil.WriteJSON(w, http.StatusOK, dvlutil.Response{
		Status: dvlutil.StatusCodeOK,
		Msg:    MsgUpdateSuccess,
		Data:   newComment,
	})
}

// func (h *Handler) PatchComment(w http.ResponseWriter, r *http.Request) {

// }
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	i, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgInvalidId,
		})
		log.Println(err)
		return
	}

	if err := h.Service.DeleteComment(uint(i)); err != nil {
		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgInternalServerErr,
		})
		log.Println(err)
		return
	}

	dvlutil.WriteJSON(w, http.StatusOK, dvlutil.Response{
		Status: dvlutil.StatusCodeOK,
		Msg:    MsgDelteSuccess,
	})

}
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.Service.GetAllComments()

	if err != nil {
		dvlutil.WriteJSON(w, http.StatusBadRequest, dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    MsgInternalServerErr,
		})
		log.Println(err)
		return
	}

	dvlutil.WriteJSON(w, http.StatusOK, dvlutil.Response{
		Status: dvlutil.StatusCodeOK,
		Msg:    MsgFetchSuccess,
		Data:   comments,
	})
}
