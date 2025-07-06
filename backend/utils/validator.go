package utils

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

// バリデーションエラーを表すカスタムエラー型
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// 定数でバリデーションの上限値を定義
const (
	MaxPageValue       = 1000  // ページ番号の上限
	MaxLimitValue      = 100   // 1ページあたりの件数上限
	MaxMovieID         = 9999999 // 映画IDの上限
	MaxGenreID         = 999    // ジャンルIDの上限
	MaxSearchQueryLen  = 100    // 検索クエリの最大文字数
	MinSearchQueryLen  = 1      // 検索クエリの最小文字数
)

// ページ番号のバリデーション
func ValidatePage(pageStr string) (int, error) {
	if pageStr == "" {
		return 1, nil // デフォルト値
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, &ValidationError{
			Field:   "page",
			Message: "ページ番号は数値で指定してください",
		}
	}

	if page < 1 {
		return 0, &ValidationError{
			Field:   "page",
			Message: "ページ番号は1以上で指定してください",
		}
	}

	if page > MaxPageValue {
		return 0, &ValidationError{
			Field:   "page",
			Message: "ページ番号が上限を超えています",
		}
	}

	return page, nil
}

// 件数制限のバリデーション
func ValidateLimit(limitStr string) (int, error) {
	if limitStr == "" {
		return 20, nil // デフォルト値
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, &ValidationError{
			Field:   "limit",
			Message: "件数は数値で指定してください",
		}
	}

	if limit < 1 {
		return 0, &ValidationError{
			Field:   "limit",
			Message: "件数は1以上で指定してください",
		}
	}

	if limit > MaxLimitValue {
		return 0, &ValidationError{
			Field:   "limit",
			Message: "件数が上限を超えています",
		}
	}

	return limit, nil
}

// 映画IDのバリデーション
func ValidateMovieID(idStr string) (int, error) {
	if idStr == "" {
		return 0, &ValidationError{
			Field:   "id",
			Message: "映画IDが指定されていません",
		}
	}

	// スラッシュが含まれていないかチェック
	if strings.Contains(idStr, "/") {
		return 0, &ValidationError{
			Field:   "id",
			Message: "無効な映画IDです",
		}
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, &ValidationError{
			Field:   "id",
			Message: "映画IDは数値で指定してください",
		}
	}

	if id < 1 {
		return 0, &ValidationError{
			Field:   "id",
			Message: "映画IDは1以上で指定してください",
		}
	}

	if id > MaxMovieID {
		return 0, &ValidationError{
			Field:   "id",
			Message: "映画IDが上限を超えています",
		}
	}

	return id, nil
}

// ジャンルIDのバリデーション
func ValidateGenreID(idStr string) (int, error) {
	if idStr == "" {
		return 0, &ValidationError{
			Field:   "genre_id",
			Message: "ジャンルIDが指定されていません",
		}
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, &ValidationError{
			Field:   "genre_id",
			Message: "ジャンルIDは数値で指定してください",
		}
	}

	if id < 1 {
		return 0, &ValidationError{
			Field:   "genre_id",
			Message: "ジャンルIDは1以上で指定してください",
		}
	}

	if id > MaxGenreID {
		return 0, &ValidationError{
			Field:   "genre_id",
			Message: "ジャンルIDが上限を超えています",
		}
	}

	return id, nil
}

// 検索クエリのバリデーション
func ValidateSearchQuery(query string) error {
	if query == "" {
		return &ValidationError{
			Field:   "query",
			Message: "検索クエリが指定されていません",
		}
	}

	// 文字数チェック（UTF-8文字として計算）
	queryLen := utf8.RuneCountInString(query)
	if queryLen < MinSearchQueryLen {
		return &ValidationError{
			Field:   "query",
			Message: "検索クエリは1文字以上で指定してください",
		}
	}

	if queryLen > MaxSearchQueryLen {
		return &ValidationError{
			Field:   "query",
			Message: "検索クエリが長すぎます",
		}
	}

	// 危険な文字列パターンのチェック
	dangerousPatterns := []string{
		"<script",
		"javascript:",
		"data:",
		"vbscript:",
		"onload=",
		"onerror=",
	}

	queryLower := strings.ToLower(query)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(queryLower, pattern) {
			return &ValidationError{
				Field:   "query",
				Message: "検索クエリに不正な文字列が含まれています",
			}
		}
	}

	return nil
}

// 汎用的な数値バリデーション
func ValidatePositiveInt(value string, fieldName string, max int) (int, error) {
	if value == "" {
		return 0, &ValidationError{
			Field:   fieldName,
			Message: fieldName + "が指定されていません",
		}
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, &ValidationError{
			Field:   fieldName,
			Message: fieldName + "は数値で指定してください",
		}
	}

	if num < 1 {
		return 0, &ValidationError{
			Field:   fieldName,
			Message: fieldName + "は1以上で指定してください",
		}
	}

	if max > 0 && num > max {
		return 0, &ValidationError{
			Field:   fieldName,
			Message: fieldName + "が上限を超えています",
		}
	}

	return num, nil
}

// バリデーションエラーかどうかを判定
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}
