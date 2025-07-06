import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Layout } from "@/components/Layout";
import { HomePage } from "@/pages/HomePage";
import { MoviesPage } from "@/pages/MoviesPage";
import { MovieDetailPage } from "@/pages/MovieDetailPage";
import { GenrePage } from "@/pages/GenrePage";
import { NotFoundPage } from "@/pages/error/NotFoundPage";
import ServerErrorPage from "@/pages/error/ServerErrorPage";
import TestErrorPage from "@/pages/error/TestErrorPage";
import { ErrorBoundary } from "./pages/error/ErrorBoundary";

function App() {
  return (
    <BrowserRouter>
      <ErrorBoundary>
        <Routes>
          <Route
            path="/"
            element={<Layout />}
            errorElement={<ServerErrorPage />}
          >
            <Route index element={<HomePage />} />
            <Route path="movies" element={<MoviesPage />} />
            <Route path="genre" element={<GenrePage />} />
            <Route path="movie/:id" element={<MovieDetailPage />} />
            <Route path="genre/:id" element={<GenrePage />} />
            <Route path="*" element={<NotFoundPage />} />
            <Route path="/test-error" element={<TestErrorPage />} />
          </Route>
        </Routes>
      </ErrorBoundary>
    </BrowserRouter>
  );
}

export default App;
