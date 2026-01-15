import SkeletonCard from "./SkeletonCard";

export default function SkeletonRow() {
  return (
    <div className="py-6 px-4 md:px-12">
      {/* Header Skeleton */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <div className="w-1.5 h-6 bg-zinc-800 rounded-full animate-pulse" />
          <div className="h-6 w-48 bg-zinc-800 rounded animate-pulse" />
        </div>
      </div>

      {/* Cards Slider Skeleton */}
      <div className="flex gap-4 overflow-x-hidden">
        {Array.from({ length: 6 }).map((_, idx) => (
          <SkeletonCard key={idx} />
        ))}
      </div>
    </div>
  );
}
