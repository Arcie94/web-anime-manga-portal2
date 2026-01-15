export default function SkeletonCard() {
  return (
    <div className="flex-none w-[140px] md:w-[180px] snap-start relative">
      <div className="aspect-[2/3] w-full bg-zinc-900 rounded-xl overflow-hidden animate-pulse">
        {/* Poster Placeholder */}
        <div className="w-full h-full bg-zinc-800/50" />
      </div>
      
      <div className="mt-2.5 space-y-2">
        {/* Title Placeholder */}
        <div className="h-4 bg-zinc-800 rounded w-3/4 animate-pulse" />
        
        {/* Metadata Placeholder */}
        <div className="flex items-center gap-2">
          <div className="h-3 bg-zinc-800 rounded w-8 animate-pulse" />
          <div className="w-1 h-1 bg-zinc-800 rounded-full" />
          <div className="h-3 bg-zinc-800 rounded w-10 animate-pulse" />
        </div>
      </div>
    </div>
  );
}
