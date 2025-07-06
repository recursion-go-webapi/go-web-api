// src/components/error/DebugErrorTrigger.tsx
import { useState } from "react";

export default function DebugErrorTrigger() {
  const [shouldThrow, setShouldThrow] = useState(false);

  // shouldThrow が true になったらレンダー中に例外を投げる
  if (shouldThrow) {
    throw new Error("手動で発生させた500エラー");
  }

  return (
    <button
      onClick={() => setShouldThrow(true)}
      className="bg-red-600 text-white px-4 py-2 rounded shadow hover:bg-red-700"
    >
      💥 エラーを発生させる
    </button>
  );
}
