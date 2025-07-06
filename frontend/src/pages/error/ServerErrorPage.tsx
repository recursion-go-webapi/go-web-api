import { ErrorLayout } from "@/components/error/ErrorLayout";

interface ServerErrorPageProps {
  error?: Error;
}

export default function ServerErrorPage({ error }: ServerErrorPageProps) {
  return (
    <ErrorLayout title="500 - サーバーエラー" imageSrc="/500.png">
      <p className="text-movie-gold font-serif italic text-lg md:text-xl leading-relaxed mb-4">
        サーバーで問題が発生しました。しばらくしてから再度お試しください。
        <br />
        必要であれば管理者にご連絡ください。
      </p>
      {error && (
        <pre className="bg-red-100 text-red-800 p-4 rounded text-sm text-left overflow-auto">
          {error.message}
        </pre>
      )}
      {/* 再読み込みボタンとホームリンク */}
      <div className="flex flex-row justify-center items-center gap-6">
        <a
          href="#"
          onClick={(e) => {
            e.preventDefault();
            window.location.reload();
          }}
          className="text-indigo-400 hover:text-indigo-200 underline transition"
        >
          🔄 再読み込み
        </a>
        <a
          href="/"
          className="text-indigo-400 hover:text-indigo-200 underline transition"
        >
          🏠 ホームに戻る
        </a>
      </div>
    </ErrorLayout>
  );
}
