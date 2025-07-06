package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// SecurityConfig はセキュリティミドルウェアの設定
type SecurityConfig struct {
	// CORS設定
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool

	// セキュリティヘッダー設定
	EnableXSSProtection      bool
	EnableFrameOptions       bool
	EnableContentTypeNoSniff bool
	EnableHSTS               bool

	// Content Security Policy
	CSPDirectives map[string]string
}

// DefaultSecurityConfig はデフォルトのセキュリティ設定を返す
func DefaultSecurityConfig() *SecurityConfig {
	// 環境変数からフロントエンドURLを取得（デフォルトはlocalhost:3000）
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Println("FRONTEND_URLが設定されていません")
	}

	return &SecurityConfig{
		// CORS設定
		AllowedOrigins: []string{
			frontendURL,
			"http://localhost:8081", // Swagger UIのCORS対応
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
			"X-Requested-With", "X-HTTP-Method-Override",
		},
		AllowCredentials: true,

		// セキュリティヘッダー
		EnableXSSProtection:      true,
		EnableFrameOptions:       true,
		EnableContentTypeNoSniff: true,
		EnableHSTS:               false, // 開発環境ではfalse、本番ではtrue

		// Content Security Policy
		CSPDirectives: map[string]string{
			"default-src": "'self'",
			"script-src":  "'self' 'unsafe-inline' 'unsafe-eval'", // React開発用
			"style-src":   "'self' 'unsafe-inline'",               // React開発用
			"img-src":     "'self' data: https://image.tmdb.org",  // TMDB画像
			"connect-src": "'self' https://api.themoviedb.org",    // TMDB API
			"font-src":    "'self' data:",
			"object-src":  "'none'",
			"base-uri":    "'self'",
			"form-action": "'self'",
		},
	}
}

// SecurityMiddleware はセキュリティヘッダーを設定するミドルウェア
func SecurityMiddleware(config *SecurityConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultSecurityConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS ヘッダーの設定
			setCORSHeaders(w, r, config)

			// プリフライトリクエスト（OPTIONS）の処理
			if r.Method == "OPTIONS" {
				handlePreflight(w, r, config)
				return
			}

			// セキュリティヘッダーの設定
			setSecurityHeaders(w, config)

			// 次のハンドラーに処理を渡す
			next.ServeHTTP(w, r)
		})
	}
}

// setCORSHeaders はCORSヘッダーを設定
func setCORSHeaders(w http.ResponseWriter, r *http.Request, config *SecurityConfig) {
	origin := r.Header.Get("Origin")

	// オリジンが許可リストに含まれているかチェック
	if isOriginAllowed(origin, config.AllowedOrigins) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// その他のCORSヘッダー
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))

	if config.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// プリフライトリクエストのキャッシュ時間（24時間）
	w.Header().Set("Access-Control-Max-Age", "86400")
}

// handlePreflight はプリフライトリクエストを処理
func handlePreflight(w http.ResponseWriter, r *http.Request, config *SecurityConfig) {
	// プリフライトリクエストで要求されたメソッドをチェック
	requestedMethod := r.Header.Get("Access-Control-Request-Method")
	if requestedMethod != "" && !isMethodAllowed(requestedMethod, config.AllowedMethods) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// プリフライトリクエストで要求されたヘッダーをチェック
	requestedHeaders := r.Header.Get("Access-Control-Request-Headers")
	if requestedHeaders != "" {
		headers := strings.Split(requestedHeaders, ",")
		for _, header := range headers {
			header = strings.TrimSpace(header)
			if !isHeaderAllowed(header, config.AllowedHeaders) {
				http.Error(w, "Header not allowed", http.StatusForbidden)
				return
			}
		}
	}

	// プリフライトリクエスト成功
	w.WriteHeader(http.StatusNoContent)
}

// setSecurityHeaders はセキュリティヘッダーを設定
func setSecurityHeaders(w http.ResponseWriter, config *SecurityConfig) {
	// XSS Protection
	if config.EnableXSSProtection {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
	}

	// Content Type Options (MIMEタイプスニッフィング防止)
	if config.EnableContentTypeNoSniff {
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	// Frame Options (クリックジャッキング防止)
	if config.EnableFrameOptions {
		w.Header().Set("X-Frame-Options", "DENY")
	}

	// HTTP Strict Transport Security (HTTPS強制) - 本番環境のみ
	if config.EnableHSTS {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}

	// Content Security Policy
	if len(config.CSPDirectives) > 0 {
		csp := buildCSP(config.CSPDirectives)
		w.Header().Set("Content-Security-Policy", csp)
	}

	// その他のセキュリティヘッダー
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
}

// isOriginAllowed はオリジンが許可されているかチェック
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
		// ワイルドカード対応（簡易版）
		if strings.Contains(allowed, "*") {
			// *.example.com のような形式に対応
			if strings.HasPrefix(allowed, "*.") {
				domain := strings.TrimPrefix(allowed, "*.")
				if strings.HasSuffix(origin, "."+domain) || origin == domain {
					return true
				}
			}
		}
	}
	return false
}

// isMethodAllowed はHTTPメソッドが許可されているかチェック
func isMethodAllowed(method string, allowedMethods []string) bool {
	for _, allowed := range allowedMethods {
		if strings.EqualFold(allowed, method) {
			return true
		}
	}
	return false
}

// isHeaderAllowed はヘッダーが許可されているかチェック
func isHeaderAllowed(header string, allowedHeaders []string) bool {
	header = strings.ToLower(header)
	for _, allowed := range allowedHeaders {
		if strings.EqualFold(allowed, header) {
			return true
		}
	}
	return false
}

// buildCSP はCSPディレクティブから文字列を構築
func buildCSP(directives map[string]string) string {
	var parts []string
	for directive, value := range directives {
		if value != "" {
			parts = append(parts, directive+" "+value)
		}
	}
	return strings.Join(parts, "; ")
}

// ProductionSecurityConfig は本番環境用のセキュリティ設定を返す
func ProductionSecurityConfig(frontendURL string) *SecurityConfig {
	config := DefaultSecurityConfig()

	// 本番環境では特定のオリジンのみ許可
	if frontendURL != "" {
		config.AllowedOrigins = []string{frontendURL}
	}

	// HTTPS強制を有効化
	config.EnableHSTS = true

	// より厳格なCSP設定
	config.CSPDirectives = map[string]string{
		"default-src":               "'self'",
		"script-src":                "'self'",                              // unsafe-inlineを削除
		"style-src":                 "'self' 'unsafe-inline'",              // CSSは許可
		"img-src":                   "'self' data: https://image.tmdb.org", // TMDB画像
		"connect-src":               "'self' https://api.themoviedb.org",   // TMDB API
		"font-src":                  "'self'",
		"object-src":                "'none'",
		"base-uri":                  "'self'",
		"form-action":               "'self'",
		"upgrade-insecure-requests": "", // HTTPSにアップグレード
	}

	return config
}
