interface ErrorLayoutProps {
  title: string;
  children: React.ReactNode;
  imageSrc?: string;
}

export function ErrorLayout({ title, children, imageSrc }: ErrorLayoutProps) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-800 relative">
      <div className="absolute inset-0 bg-gradient-to-br from-slate-950/50 via-slate-900/30 to-slate-800/20" />
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_30%_20%,rgba(251,191,36,0.1),transparent_60%)] opacity-60" />
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_70%_80%,rgba(251,191,36,0.05),transparent_60%)] opacity-40" />

      <div className="min-h-screen flex items-center justify-center p-4 relative text-white z-10">
        {/* メインコンテンツカード */}
        <div className="relative w-full max-w-6xl">
          <div className="backdrop-blur-xl bg-white/10 rounded-3xl border border-white/20 shadow-2xl p-6 md:p-8 text-center">
            {/* レイアウト*/}
            <div className="flex flex-col lg:landscape:flex-row lg:landscape:items-center lg:landscape:text-left lg:landscape:gap-12">
              {/* タイトルとコンテンツ */}
              <div className="flex-1 lg:landscape:order-1">
                <h1 className="text-white text-2xl md:text-3xl lg:text-4xl xl:landscape:text-5xl font-extrabold tracking-wider mb-6 lg:landscape:mb-8 relative">
                  <span className="bg-gradient-to-r from-white via-white to-white/80 bg-clip-text text-transparent drop-shadow-2xl">
                    {title}
                  </span>
                  <div className="absolute -bottom-2 left-1/2 lg:landscape:left-0 transform -translate-x-1/2 lg:landscape:translate-x-0 w-24 h-1 bg-gradient-to-r from-transparent via-white/50 to-transparent lg:landscape:from-white/50 lg:landscape:via-white/50 lg:landscape:to-transparent rounded-full"></div>
                </h1>

                {/* Content with enhanced styling - moved here for landscape */}
                <div className="relative block lg:landscape:block">
                  <div className="text-movie-gold font-serif italic text-lg md:text-xl leading-relaxed relative z-10">
                    <div className="absolute inset-0 bg-neutral-900/40 rounded-2xl backdrop-blur-sm border border-white/10"></div>
                    <div className="relative p-6 text-left">{children}</div>
                  </div>
                </div>
              </div>

              {/* {右側} */}
              {imageSrc && (
                <div className="flex-shrink-0 mb-6 lg:landscape:mb-0 lg:landscape:order-2 flex justify-center">
                  <div className="relative group">
                    <div className="absolute inset-0 rounded-2xl opacity-75 group-hover:opacity-100"></div>
                    <div className="relative w-48 md:w-56 lg:landscape:w-64 xl:landscape:w-72 border-2 border-white/30 rounded-2xl shadow-2xl bg-black/20 backdrop-blur-sm overflow-hidden">
                      <img
                        src={imageSrc || "/placeholder.svg"}
                        alt="Error visual"
                        className="w-full object-cover"
                      />
                      <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent"></div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
