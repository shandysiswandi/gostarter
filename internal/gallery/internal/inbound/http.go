package inbound

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/gallery/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/http/middleware"
	"github.com/shandysiswandi/gostarter/pkg/http/serve"
)

func RegisterHTTP(router *httprouter.Router, h *Endpoint) {
	serve := serve.New(
		serve.WithMiddlewares(middleware.Recovery),
		serve.WithErrorCodec(func(ctx context.Context, w http.ResponseWriter, err error) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			code := http.StatusInternalServerError
			if sc, ok := err.(serve.StatusCoder); ok {
				code = sc.StatusCode()
			}

			w.WriteHeader(code)
			msg := "internal server error"

			var e *goerror.GoError
			if errors.As(err, &e) {
				msg = e.Msg()
			}

			//nolint:errcheck // ignore for this, it never error
			_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
		}))

	router.GET("/galleries/:id", h.GetImage)
	router.Handler(http.MethodPost, "/galleries/upload", serve.Endpoint(h.Upload))
}

type Endpoint struct {
	UploadUC   domain.Upload
	GetImageUC domain.GetImage
}

func (e *Endpoint) Upload(ctx context.Context, r *http.Request) (any, error) {
	if err := r.ParseMultipartForm(5 << 20); /*5MB*/ err != nil {
		return nil, goerror.NewInvalidFormat("maximum file is 2MB", err)
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, goerror.NewInvalidFormat("required image field", err)
	}
	defer file.Close()

	bt, err := io.ReadAll(file)
	if err != nil {
		return nil, goerror.NewServer("internal", err)
	}

	resp, err := e.UploadUC.Call(ctx, domain.UploadInput{File: bt})
	if err != nil {
		return nil, err
	}

	return map[string]string{"image": "http://localhost:8081/galleries/" + resp.Path}, nil
}

func (e *Endpoint) GetImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.ParseUint(p.ByName("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	_, err = e.GetImageUC.Call(r.Context(), domain.GetImageInput{ID: id})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusOK)
}
