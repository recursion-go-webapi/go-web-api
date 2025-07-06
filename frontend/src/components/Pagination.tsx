import { motion } from "framer-motion";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  loading?: boolean;
}

export function Pagination({ currentPage, totalPages, onPageChange, loading }: PaginationProps) {
  if (totalPages <= 1) return null;

  const getVisiblePages = () => {
    const delta = 2;
    const range = [];
    const rangeWithDots = [];

    for (let i = Math.max(2, currentPage - delta); i <= Math.min(totalPages - 1, currentPage + delta); i++) {
      range.push(i);
    }

    // 最初のページを追加
    if (currentPage - delta > 2) {
      rangeWithDots.push(1, '...');
    } else {
      rangeWithDots.push(1);
    }

    // 中間のページを追加（重複を避けるため、1を除外）
    rangeWithDots.push(...range.filter(page => page !== 1 && page !== totalPages));

    // 最後のページを追加
    if (currentPage + delta < totalPages - 1) {
      rangeWithDots.push('...', totalPages);
    } else if (totalPages > 1) {
      rangeWithDots.push(totalPages);
    }

    // 重複を削除し、順序を保持
    const uniquePages = [];
    const seen = new Set();

    for (const page of rangeWithDots) {
      const key = typeof page === 'string' ? `${page}-${uniquePages.length}` : page;
      if (!seen.has(key)) {
        seen.add(key);
        uniquePages.push(page);
      }
    }

    return uniquePages;
  };

  const visiblePages = getVisiblePages();

  return (
    <div className="flex justify-center items-center gap-3 mt-12 mb-8">
      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage <= 1 || loading}
        className="px-4 py-2 rounded-lg bg-slate-800 border border-slate-600 text-slate-300 hover:text-amber-400 hover:border-amber-400 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 font-medium"
      >
        ← 前へ
      </motion.button>

      <div className="flex gap-2">
        {visiblePages.map((page, index) => (
          page === '...' ? (
            <span key={`dots-${index}-${currentPage}`} className="px-3 py-2 text-slate-400 font-bold">
              ...
            </span>
          ) : (
            <motion.button
              key={`page-${page}-${index}`}
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.9 }}
              onClick={() => onPageChange(page as number)}
              disabled={loading}
              className={`min-w-[44px] h-[44px] rounded-lg font-semibold transition-all duration-300 ${currentPage === page
                  ? 'bg-amber-500 text-slate-900 shadow-lg shadow-amber-500/30 border-2 border-amber-400'
                  : 'bg-slate-800 text-slate-300 border-2 border-slate-600 hover:text-amber-400 hover:border-amber-400 hover:bg-slate-700'
                } ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
            >
              {page}
            </motion.button>
          )
        ))}
      </div>

      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage >= totalPages || loading}
        className="px-4 py-2 rounded-lg bg-slate-800 border border-slate-600 text-slate-300 hover:text-amber-400 hover:border-amber-400 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 font-medium"
      >
        次へ →
      </motion.button>
    </div>
  );
}
