package api

import (
	"fmt"
	"net/http"

	"gopkg.in/dgrijalva/jwt-go.v3"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/treeverse/lakefs/api/gen/models"
	"github.com/treeverse/lakefs/api/gen/restapi"
	"github.com/treeverse/lakefs/api/gen/restapi/operations"
	"github.com/treeverse/lakefs/auth"
	"github.com/treeverse/lakefs/block"
	"github.com/treeverse/lakefs/db"
	"github.com/treeverse/lakefs/httputil"
	"github.com/treeverse/lakefs/index"
	"github.com/treeverse/lakefs/logging"
	_ "github.com/treeverse/lakefs/statik"
	"github.com/treeverse/lakefs/stats"
)

const (
	RequestIdHeaderName        = "X-Request-ID"
	LoggerServiceName          = "rest_api"
	JWTAuthorizationHeaderName = "X-JWT-Authorization"
)

var (
	ErrAuthenticationFailed = errors.New(http.StatusUnauthorized, "error authenticating request")
)

type Handler struct {
	meta        auth.MetadataManager
	index       index.Index
	blockStore  block.Adapter
	authService auth.Service
	stats       stats.Collector
	migrator    db.Migrator
	apiServer   *restapi.Server
	handler     *http.ServeMux
	server      *http.Server
	logger      logging.Logger
}

func NewHandler(
	index index.Index,
	blockStore block.Adapter,
	authService auth.Service,
	meta auth.MetadataManager,
	stats stats.Collector,
	migrator db.Migrator,
	logger logging.Logger,
) http.Handler {
	logger.Info("initialized OpenAPI handler")
	s := &Handler{
		index:       index,
		blockStore:  blockStore,
		authService: authService,
		meta:        meta,
		stats:       stats,
		migrator:    migrator,
		logger:      logger,
	}
	s.buildAPI()
	return s.handler
}

// JwtTokenAuth decodes, validates and authenticates a user that exists
// in the X-JWT-Authorization header.
// This header either exists natively, or is set using a token
func (s *Handler) JwtTokenAuth() func(string) (*models.User, error) {
	logger := logging.Default().WithField("auth", "jwt")
	return func(tokenString string) (*models.User, error) {
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return s.authService.SecretStore().SharedSecret(), nil
		})
		if err != nil {
			return nil, ErrAuthenticationFailed
		}
		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok || !token.Valid {
			return nil, ErrAuthenticationFailed
		}
		userData, err := s.authService.GetUser(claims.Subject)
		if err != nil {
			logger.WithField("subject", claims.Subject).Warn("could not find user for token")
			return nil, ErrAuthenticationFailed
		}
		return &models.User{
			ID: userData.DisplayName,
		}, nil
	}
}

// BasicAuth returns a function that hooks into Swagger's basic Auth provider
// it uses the Auth.Service provided to ensure credentials are valid
func (s *Handler) BasicAuth() func(accessKey, secretKey string) (user *models.User, err error) {
	logger := logging.Default().WithField("auth", "basic")
	return func(accessKey, secretKey string) (user *models.User, err error) {
		credentials, err := s.authService.GetCredentials(accessKey)
		if err != nil {
			logger.WithError(err).WithField("access_key", accessKey).Warn("could not get access key for login")
			return nil, ErrAuthenticationFailed
		}
		if secretKey != credentials.AccessSecretKey {
			logger.WithField("access_key", accessKey).Warn("access key secret does not match")
			return nil, ErrAuthenticationFailed
		}
		userData, err := s.authService.GetUserById(credentials.UserId)
		if err != nil {
			logger.WithField("access_key", accessKey).Warn("could not find user for key pair")
			return nil, ErrAuthenticationFailed
		}
		return &models.User{
			ID: userData.DisplayName,
		}, nil
	}
}

func (s *Handler) setupHandler(api http.Handler, ui http.Handler, setup http.Handler) {
	mux := http.NewServeMux()
	// api handler
	mux.Handle("/api/", api)
	// swagger
	mux.Handle("/swagger.json", api)
	// setup system
	mux.Handle(SetupLakeFSRoute, setup)
	// otherwise, serve  UI
	mux.Handle("/", ui)

	s.handler = mux
}

// buildAPI wires together the JWT and basic authenticator and registers all relevant API handlers
func (s *Handler) buildAPI() {
	swaggerSpec, _ := loads.Analyzed(restapi.SwaggerJSON, "")

	api := operations.NewLakefsAPI(swaggerSpec)
	api.Logger = func(msg string, ctx ...interface{}) {
		logging.Default().WithField("logger", "swagger").Debugf(msg, ctx)
	}
	api.BasicAuthAuth = s.BasicAuth()
	api.JwtTokenAuth = s.JwtTokenAuth()

	// bind our handlers to the server
	NewController(s.index, s.authService, s.blockStore, s.stats, s.logger).Configure(api)

	// setup host/port
	s.apiServer = restapi.NewServer(api)
	s.apiServer.ConfigureAPI()

	s.setupHandler(
		// api handler
		httputil.LoggingMiddleware(
			RequestIdHeaderName,
			logging.Fields{"service_name": LoggerServiceName},
			cookieToAPIHeader(s.apiServer.GetHandler()),
		),

		// ui handler
		UIHandler(s.authService),

		// setup handler
		httputil.LoggingMiddleware(
			RequestIdHeaderName,
			logging.Fields{"service_name": LoggerServiceName},
			setupLakeFSHandler(s.authService, s.meta, s.migrator, s.stats),
		),
	)
}

func cookieToAPIHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read cookie (no need to validate, this will be done in the API
		cookie, err := r.Cookie(JWTCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		// header found
		r.Header.Set(JWTAuthorizationHeaderName, cookie.Value)
		next.ServeHTTP(w, r)
	})
}