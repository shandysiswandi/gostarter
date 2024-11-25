package framework

import (
	"encoding/json"
	"mime"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// transportPOST implements a transport layer for handling POST requests
// in a GraphQL server. It supports JSON-encoded GraphQL queries and manages
// request execution.
type transportPOST struct{}

// Supports determines if the transport supports the incoming HTTP request.
// It checks if the method is POST, the "Upgrade" header is absent, and
// the Content-Type is "application/json".
func (transportPOST) Supports(r *http.Request) bool {
	if r.Header.Get("Upgrade") != "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get(keyContentType))
	if err != nil {
		return false
	}

	return r.Method == http.MethodPost && mediaType == "application/json"
}

// Do executes the GraphQL request using the provided executor. It reads the
// request body, decodes it into GraphQL parameters, and processes the operation
// context and response. Errors during decoding or execution are handled and
// written as JSON responses.
func (transportPOST) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	ctx := r.Context()
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
		data := exec.DispatchError(ctx, gqlerror.List{gqlerror.Errorf("invalid request body")})
		writeJSON(w, data, http.StatusBadRequest)

		return
	}

	rc, err := exec.CreateOperationContext(ctx, params)
	if err != nil {
		data := exec.DispatchError(graphql.WithOperationContext(ctx, rc), err)
		writeJSON(w, data, http.StatusUnprocessableEntity)

		return
	}

	responseFunc, ctx := exec.DispatchOperation(ctx, rc)
	resp := responseFunc(ctx)

	code := http.StatusOK
	var goErr *goerror.GoError
	if resp.Errors.As(&goErr) {
		code = goErr.StatusCode()
	}

	writeJSON(w, resp, code)
}
