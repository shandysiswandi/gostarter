package gql

import (
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type transportPOST struct{}

func (transportPOST) Supports(r *http.Request) bool {
	if r.Header.Get("Upgrade") != "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return false
	}

	return r.Method == http.MethodPost && mediaType == "application/json"
}

func (transportPOST) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := &graphql.RawParams{}
	start := graphql.Now()
	params.Headers = r.Header
	params.ReadTime = graphql.TraceTiming{
		Start: start,
		End:   graphql.Now(),
	}

	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	if err := dec.Decode(params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		gqlErr := gqlerror.Errorf("invalid request body")
		writeJSON(w, exec.DispatchError(ctx, gqlerror.List{gqlErr}))

		return
	}

	rc, err := exec.CreateOperationContext(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeJSON(w, exec.DispatchError(graphql.WithOperationContext(ctx, rc), err))

		return
	}

	responseFunc, ctx := exec.DispatchOperation(ctx, rc)
	resp := responseFunc(ctx)

	var goErr *goerror.GoError
	if resp.Errors.As(&goErr) {
		w.WriteHeader(goErr.StatusCode())
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, resp)
}

func writeJSON(w io.Writer, data *graphql.Response) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("gql.server.json.NewEncoder(w).Encode", err)
	}
}
