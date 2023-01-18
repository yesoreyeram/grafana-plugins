package restds

type AuthType string

const (
	AuthTypeNone         AuthType = "none"
	AuthTypeBasic        AuthType = "basicAuth"
	AuthTypeBearerToken  AuthType = "bearerToken"
	AuthTypeApiKey       AuthType = "apiKey"
	AuthTypeForwardOauth AuthType = "oauthPassThru"
	AuthTypeDigestAuth   AuthType = "digestAuth"
	AuthTypeOAuth2       AuthType = "oauth2"
)

type APIKeyType string

const (
	ApiKeyTypeHeader APIKeyType = "header"
	ApiKeyTypeQuery  APIKeyType = "query"
)

type OAuth2Type string

const (
	OAuth2TypeClientCredentials OAuth2Type = "client_credentials"
	OAuth2TypeJWT               OAuth2Type = "jwt"
)

type Config struct {
	BaseURL              string
	AuthenticationMethod AuthType
	BasicAuthUser        string
	BasicAuthPassword    string
	ApiKeyType           APIKeyType
	ApiKeyKey            string
	ApiKeyValue          string
	BearerToken          string
	Headers              map[string]string
	QueryParams          map[string]string
	OAuth2Settings       struct {
		Type           OAuth2Type
		TokenURL       string
		ClientID       string
		ClientSecret   string
		Email          string
		PrivateKeyID   string
		PrivateKey     string
		Subject        string
		Scopes         []string
		EndpointParams map[string]string
	}
}

func (c *Config) Validate() error {
	return nil
}
