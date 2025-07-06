import DebugErrorTrigger from "../../components/error/DebugErrorTrigger";

export default function TestErrorPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen gap-6">
      <h1 className="text-2xl font-bold">500 エラー表示テスト</h1>
      <p>このボタンを押すと意図的にエラーを発生させます。</p>
      <DebugErrorTrigger />
    </div>
  );
}
