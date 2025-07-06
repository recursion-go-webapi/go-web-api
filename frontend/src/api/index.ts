import type { MoviesResponse, MovieDetail, GenreMovieListResponse, APIError, GenreListResponse } from '@/types/movie';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

if (!import.meta.env.VITE_API_BASE_URL) {
  console.warn('VITE_API_BASE_URLが環境変数に設定されていません。デフォルト値を使用します:', API_BASE_URL);
}

const request = async <T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> => {
  const url = `${API_BASE_URL}${endpoint}`;

  try {
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      const errorData: APIError = await response.json();
      throw new Error(errorData.message || `HTTP ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    if (error instanceof Error) {
      throw error;
    }
    throw new Error('ネットワークエラーが発生しました');
  }
};

export const getMovies = (page: number = 1): Promise<MoviesResponse> => {
  return request<MoviesResponse>(`/api/movies?page=${page}`);
};

export const getMovieDetail = (id: number): Promise<MovieDetail> => {
  return request<MovieDetail>(`/api/movie/${id}`);
};

export const searchMovies = (query: string, page: number = 1): Promise<MoviesResponse> => {
  const encodedQuery = encodeURIComponent(query);
  return request<MoviesResponse>(`/api/movies/search?query=${encodedQuery}&page=${page}`);
};

export const getMoviesByGenre = (genreId: number, page: number = 1): Promise<GenreMovieListResponse> => {
  return request<GenreMovieListResponse>(`/api/movies/genre?genre_id=${genreId}&page=${page}`);
};

export const getGenres = (): Promise<{ genres: { id: number; name: string }[] }> => {
  return request<GenreListResponse>('/api/genres');
};

export const healthCheck = (): Promise<{ status: string }> => {
  return request<{ status: string }>('/healthz');
};
