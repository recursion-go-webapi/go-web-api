package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-movie-explorer/middleware"
	"go-movie-explorer/services"
	"go-movie-explorer/utils"
)

// 映画一覧取得APIハンドラー /api/movies
func MoviesHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// クエリパラメータ取得とバリデーション
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	
	// ページ番号のバリデーション
	page, err := utils.ValidatePage(pageStr)
	if err != nil {
		if utils.IsValidationError(err) {
			return middleware.NewBadRequestError(err.Error())
		}
		return middleware.NewInternalServerError(fmt.Sprintf("ページ番号の検証中にエラーが発生しました: %v", err))
	}
	
	// 件数制限のバリデーション
	limit, err := utils.ValidateLimit(limitStr)
	if err != nil {
		if utils.IsValidationError(err) {
			return middleware.NewBadRequestError(err.Error())
		}
		return middleware.NewInternalServerError(fmt.Sprintf("件数制限の検証中にエラーが発生しました: %v", err))
	}
	
	// 現在はlimitは使用していないが、将来的に使用する予定
	_ = limit

	// サービス層でTMDB APIから映画一覧を取得（API仕様変更や他サービス連携時はここを編集）
	moviesResp, err := services.GetMoviesFromTMDB(page)
	if err != nil {
		return middleware.NewInternalServerError(fmt.Sprintf("TMDB API呼び出し失敗: %v", err))
	}

	// レスポンスをJSONで返却（レスポンス形式を変更したい場合はここを編集）
	if err := json.NewEncoder(w).Encode(moviesResp); err != nil {
		return middleware.NewInternalServerError(fmt.Sprintf("JSONレスポンスのエンコードに失敗しました: %v", err))
	}
	return nil
}

// 映画詳細取得ハンドラー /api/movies/{id}
func MovieDetailHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	prefix := "/api/movie/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return middleware.NewBadRequestError(fmt.Sprintf("無効なパス: %s", r.URL.Path))
	}
	id := strings.TrimPrefix(r.URL.Path, prefix)
	
	// 映画IDのバリデーション
	movieID, err := utils.ValidateMovieID(id)
	if err != nil {
		if utils.IsValidationError(err) {
			return middleware.NewBadRequestError(err.Error())
		}
		return middleware.NewInternalServerError(fmt.Sprintf("映画IDの検証中にエラーが発生しました: %v", err))
	}
	// サービス層でTMDB APIから映画詳細を取得
	movieDetail, err := services.GetMovieDetailFromTMDB(movieID)

	if err != nil {
		return middleware.NewInternalServerError(fmt.Sprintf("映画詳細取得失敗: %v", err))
	}

	w.WriteHeader(http.StatusOK)
	// レスポンスをJSONで返却
	if err := json.NewEncoder(w).Encode(movieDetail); err != nil {
		return middleware.NewInternalServerError(fmt.Sprintf("JSONレスポンスのエンコードに失敗しました: %v", err))
	}
	return nil
}

// 映画検索APIハンドラー /api/movies/search
func SearchMoviesHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	// クエリパラメータ取得とバリデーション
	query := r.URL.Query().Get("query")
	
	// 検索クエリのバリデーション
	if err := utils.ValidateSearchQuery(query); err != nil {
		if utils.IsValidationError(err) {
			return middleware.NewBadRequestError(err.Error())
		}
		return middleware.NewInternalServerError(fmt.Sprintf("検索クエリの検証中にエラーが発生しました: %v", err))
	}

	// ページ番号のバリデーション
	pageStr := r.URL.Query().Get("page")
	page, err := utils.ValidatePage(pageStr)
	if err != nil {
		if utils.IsValidationError(err) {
			return middleware.NewBadRequestError(err.Error())
		}
		return middleware.NewInternalServerError(fmt.Sprintf("ページ番号の検証中にエラーが発生しました: %v", err))
	}

	// サービス層でTMDB APIから映画検索結果を取得
	moviesResp, err := services.SearchMoviesFromTMDB(query, page)
	if err != nil {
		return middleware.NewInternalServerError(fmt.Sprintf("TMDB 検索API呼び出し失敗: %v", err))
	}

	w.WriteHeader(http.StatusOK)
	// レスポンスをJSONで返却
	if err := json.NewEncoder(w).Encode(moviesResp); err != nil {
		return middleware.NewInternalServerError(fmt.Sprintf("JSONレスポンスのエンコードに失敗しました: %v", err))
	}
	return nil
}

// - /api/movies/popular : 人気映画ランキング（今後追加予定）
//
// 新しいエンドポイントを追加する場合は、このファイルにハンドラー関数を追記してください。

func ListMoviesByGenreHandler(w http.ResponseWriter, r *http.Request) error {
	genreIDStr := r.URL.Query().Get("genre_id")
	pageStr := r.URL.Query().Get("page")

	// ジャンルIDのバリデーション
	genreID, err := utils.ValidateGenreID(genreIDStr)
	if err != nil {
		if utils.IsValidationError(err) {
			return middleware.NewBadRequestError(err.Error())
		}
		return middleware.NewInternalServerError(fmt.Sprintf("ジャンルIDの検証中にエラーが発生しました: %v", err))
	}

	// ページ番号のバリデーション
	page, err := utils.ValidatePage(pageStr)
	if err != nil {
		if utils.IsValidationError(err) {
			return middleware.NewBadRequestError(err.Error())
		}
		return middleware.NewInternalServerError(fmt.Sprintf("ページ番号の検証中にエラーが発生しました: %v", err))
	}

	result, err := services.GetMoviesByGenreFromTMDB(genreID, page)
	if err != nil {
		return middleware.NewInternalServerError(fmt.Sprintf("ジャンルの取得に失敗しました。: %v", err))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	return nil
}

// ジャンル一覧取得APIハンドラー  /api/genres
func GenresHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	// サービス層でTMDB APIからジャンル一覧を取得
	genresResp, err := services.GetGenresFromTMDB()
	if err != nil {
		return middleware.NewInternalServerError(fmt.Sprintf("TMDB ジャンル一覧の取得に呼び出し失敗しました: %v", err))
	}

	w.WriteHeader(http.StatusOK)
	// レスポンスをJSONで返却
	if err := json.NewEncoder(w).Encode(genresResp); err != nil {
		return middleware.NewInternalServerError(fmt.Sprintf("JSONレスポンスのエンコードに失敗しました: %v", err))
	}
	return nil
}
