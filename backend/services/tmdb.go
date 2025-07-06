package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"go-movie-explorer/models"
)

const BaseURL = "https://api.themoviedb.org/3"

// シングルトンHTTPクライアント
var (
	httpClient *http.Client
	clientOnce sync.Once
)

// getHTTPClient はシングルトンパターンでHTTPクライアントを取得
func getHTTPClient() *http.Client {
	clientOnce.Do(func() {
		// コネクションプールとKeep-Alive設定
		transport := &http.Transport{
			MaxIdleConns:        100,              // 最大アイドル接続数
			MaxConnsPerHost:     10,               // ホスト毎の最大接続数
			MaxIdleConnsPerHost: 10,               // ホスト毎の最大アイドル接続数
			IdleConnTimeout:     90 * time.Second, // アイドル接続のタイムアウト
		}

		httpClient = &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second, // デフォルトタイムアウト
		}
	})
	return httpClient
}

// getPingHTTPClient はPing用の短いタイムアウトを持つクライアントを取得
func getPingHTTPClient() *http.Client {
	baseClient := getHTTPClient()
	// Ping用に5秒タイムアウトのクライアントを作成（Transportは共有）
	return &http.Client{
		Transport: baseClient.Transport,
		Timeout:   5 * time.Second,
	}
}

// TMDBのAPIキーを環境変数から取得
func GetTMDBApiKey() string {
	return os.Getenv("TMDB_API_KEY")
}

// --- TMDB疎通確認用Pinger（/healthz用）---
type TmdbPinger struct{}

func (t *TmdbPinger) Name() string { return "TMDB" }

func (t *TmdbPinger) Ping(ctx context.Context) error {
	apiKey := GetTMDBApiKey()
	if apiKey == "" {
		return fmt.Errorf("TMDB_API_KEYが設定されていません")
	}
	url := BaseURL + "/discover/movie?page=1"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("TMDB Pingリクエスト作成失敗: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	client := getPingHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("TMDB Pingリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TMDB Ping status not OK: %d", resp.StatusCode)
	}
	return nil
}

// --- 映画一覧取得（/discover/movie）---
func GetMoviesFromTMDB(page int) (*models.MoviesResponse, error) {
	apiKey := GetTMDBApiKey()
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEYが設定されていません")
	}

	// APIリクエストURL生成
	url := fmt.Sprintf("%s/discover/movie?page=%d", BaseURL, page)
	client := getHTTPClient()

	// HTTPリクエスト作成
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	// TMDB API呼び出し
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TMDB APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB APIエラー: status=%d", resp.StatusCode)
	}

	// TMDBレスポンスを直接MoviesResponseにデコード
	var moviesResp models.MoviesResponse
	if err := json.NewDecoder(resp.Body).Decode(&moviesResp); err != nil {
		return nil, fmt.Errorf("TMDBレスポンスのデコード失敗: %w", err)
	}

	return &moviesResp, nil
}

// --- 映画詳細取得（/movie/{id}）---
func GetMovieDetailFromTMDB(id int) (*models.MovieDetail, error) {
	apiKey := GetTMDBApiKey()
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEYが設定されていません")
	}

	// APIリクエストURL生成
	url := fmt.Sprintf("%s/movie/%d", BaseURL, id)
	client := getHTTPClient()

	// HTTPリクエスト作成
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	// TMDB API呼び出し
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TMDB APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB APIエラー: status=%d", resp.StatusCode)
	}
	var tmdbResp models.TmdbMovieDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		return nil, fmt.Errorf("TMDBレスポンスのデコード失敗: %w", err)
	}

	// TMDBのレスポンスを独自のMovieDetailに変換
	return &models.MovieDetail{
		ID:               tmdbResp.ID,
		Title:            tmdbResp.Title,
		OriginalTitle:    tmdbResp.OriginalTitle,
		Overview:         tmdbResp.Overview,
		ReleaseDate:      tmdbResp.ReleaseDate,
		PosterPath:       tmdbResp.PosterPath,
		BackdropPath:     tmdbResp.BackdropPath,
		Genres:           tmdbResp.Genres,
		Homepage:         tmdbResp.Homepage,
		IMDBID:           tmdbResp.IMDBID,
		Popularity:       tmdbResp.Popularity,
		Budget:           tmdbResp.Budget,
		OriginCountry:    tmdbResp.OriginCountry,
		OriginalLanguage: tmdbResp.OriginalLanguage,
	}, nil
}

// --- 映画検索（/search/movie）---
func SearchMoviesFromTMDB(query string, page int) (*models.MoviesResponse, error) {
	apiKey := GetTMDBApiKey()
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEYが設定されていません")
	}

	if query == "" {
		return nil, fmt.Errorf("検索クエリが指定されていません")
	}

	// APIリクエストURL生成（クエリパラメータをエスケープ）
	apiURL := fmt.Sprintf("%s/search/movie?query=%s&page=%d", BaseURL, url.QueryEscape(query), page)
	client := getHTTPClient()

	// HTTPリクエスト作成
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	// TMDB API呼び出し
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TMDB APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB APIエラー: status=%d", resp.StatusCode)
	}

	// TMDBレスポンスを直接MoviesResponseにデコード
	var moviesResp models.MoviesResponse
	if err := json.NewDecoder(resp.Body).Decode(&moviesResp); err != nil {
		return nil, fmt.Errorf("TMDBレスポンスのデコード失敗: %w", err)
	}

	return &moviesResp, nil
}

// --- 人気映画ランキング取得（/movie/popular）---
// func GetPopularMoviesFromTMDB(page int) (*models.MoviesResponse, error) {
// 	// TODO: 人気映画ランキングAPIの実装予定
// 	// 例: /movie/popular?page=...
// 	return nil, nil
// }

// --- ジャンル別映画取得（/discover/movie?with_genres=）---
func GetMoviesByGenreFromTMDB(genreID, page int) (*models.GenreMovieListResponse, error) {
	apiKey := GetTMDBApiKey()
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEYが設定されていません")
	}

	url := fmt.Sprintf("%s/discover/movie?with_genres=%d&page=%d", BaseURL, genreID, page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	client := getHTTPClient()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TMDb APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDb APIエラー: %s", resp.Status)
	}

	var tmdbResp models.TMDBGenreMovieList
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		return nil, fmt.Errorf("TMDbレスポンスのデコード失敗: %w", err)
	}

	movies := append([]models.MovieByGenre{}, tmdbResp.Results...)

	return &models.GenreMovieListResponse{
		GenreID:      genreID,
		Page:         tmdbResp.Page,
		PerPage:      len(movies),
		TotalPages:   tmdbResp.TotalPages,
		TotalResults: tmdbResp.TotalResults,
		Results:      movies,
	}, nil
}

// --- ジャンル一覧取得（/genre/movie/list）---
func GetGenresFromTMDB() (*models.GenreListResponse, error) {
	apiKey := GetTMDBApiKey()
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEYが設定されていません")
	}
	url := fmt.Sprintf("%s/genre/movie/list", BaseURL)
	client := getHTTPClient()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TMDB APPリクエスト失敗: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB APIエラー: status=%d", resp.StatusCode)
	}

	// TMDBレスポンスを直接GenreListResponseにデコード
	var tmdbResp models.GenreListResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		return nil, fmt.Errorf("TMDBレスポンスのデコード失敗: %w", err)
	}

	return &tmdbResp, nil
}
