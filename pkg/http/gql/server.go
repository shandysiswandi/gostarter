package gql

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/vektah/gqlparser/v2/ast"
)

func ServerDefault(es graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(es)

	srv.AddTransport(DefaultPOST{})
	// srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

type DefaultPOST struct{}

func (DefaultPOST) Supports(r *http.Request) bool {
	if r.Header.Get("Upgrade") != "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return false
	}

	return r.Method == http.MethodPost && mediaType == "application/json"
}

func (DefaultPOST) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := &graphql.RawParams{}
	start := graphql.Now()
	params.Headers = r.Header
	params.ReadTime = graphql.TraceTiming{
		Start: start,
		End:   graphql.Now(),
	}

	// body, _ := io.ReadAll(r.Body)
	// log.Println(string(body))

	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	if err := dec.Decode(params); err != nil {
		log.Println("decode request body", err)

		return
	}

	rc, err := exec.CreateOperationContext(ctx, params)
	if err != nil {
		log.Println("CreateOperationContext", err)

		return
	}

	var responses graphql.ResponseHandler
	responses, ctx = exec.DispatchOperation(ctx, rc)

	if err := json.NewEncoder(w).Encode(responses(ctx)); err != nil {
		log.Println("json.NewEncoder(w).Encode", err)
	}
}
