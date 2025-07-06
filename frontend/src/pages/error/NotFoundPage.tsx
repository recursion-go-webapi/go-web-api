import { ErrorLayout } from "@/components/error/ErrorLayout";
import { Link } from "react-router-dom";

// 映画の名言をもじって作ったメッセージ
const messages = [
  { text: "🤖 404 Not Found... I'll be back. 🔥", movieId: 87101 }, // ターミネーター
  {
    text: "💃 City of stars… are you shining just for 404? 🕺",
    movieId: 313369,
  }, // ララランド
  {
    text: "🧙‍♂️ This page missed the Hogwarts Express—so you’ve arrived at 404. 🦉",
    movieId: 671,
  }, // ハリーポッター
  {
    text: "❄️ “Let it go—this page never bothered us anyway. 404!” ⛄️",
    movieId: 109445,
  }, // アナ雪
  {
    text: "🏃Pages are like a box of chocolates—this one’s missing. 404. 🍫",
    movieId: 13,
  }, // フォレスト・ガンプ
];

const random = messages[Math.floor(Math.random() * messages.length)];

export function NotFoundPage() {
  return (
    <ErrorLayout title="404 - ページが見つかりません" imageSrc="/404.png">
      <div className="space-y-6">
        <p className="text-lg">
          <Link
            to={`/movie/${random.movieId}`}
            className="underline hover:text-indigo-300 transition"
          >
            {random.text}
          </Link>
        </p>

        <div className="flex flex-row items-start space-x-4 mt-4 text-base">
          <Link
            to="/"
            className="text-indigo-400 hover:text-indigo-200 underline transition"
          >
            ← ホームに戻る
          </Link>
          <Link
            to="/movies"
            className="text-indigo-400 hover:text-indigo-200 underline transition"
          >
            🔍 映画を検索する
          </Link>
        </div>
      </div>
    </ErrorLayout>
  );
}
