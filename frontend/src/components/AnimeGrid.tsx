import { FaPlay } from "react-icons/fa";

interface Props {
    items: any[];
}

export default function AnimeGrid({ items }: Props) {
    if (!items || items.length === 0) {
        return (
            <div className="flex flex-col items-center justify-center py-20 text-zinc-500">
                <p className="text-lg">No anime found for this category.</p>
            </div>
        );
    }

    return (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 md:gap-6">
            {items.map((item, idx) => (
                <a 
                    key={idx} 
                    href={`/anime/${item.slug || item.animeId}`}
                    className="relative group/card bg-zinc-900 rounded-xl overflow-hidden shadow-lg hover:shadow-2xl transition-all duration-300"
                >
                    <div className="aspect-[2/3] w-full relative overflow-hidden">
                        <img 
                            src={item.poster || item.cover || item.image || 'https://placehold.co/200x300/222/999?text=No+Image'} 
                            alt={item.title} 
                            className="w-full h-full object-cover transition-transform duration-500 group-hover/card:scale-110"
                            loading="lazy"
                        />
                        
                        {/* Overlay Gradient */}
                        <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/20 to-transparent opacity-60 group-hover/card:opacity-40 transition-opacity" />

                        {/* Badge Top Left */}
                        <div className="absolute top-2 left-2 bg-zinc-900/80 backdrop-blur text-white text-[10px] font-bold px-2 py-0.5 rounded border border-white/10">
                            TV
                        </div>

                        {/* Badge Top Right (Episode) */}
                        <div className="absolute top-2 right-2 bg-red-600 text-white text-[10px] font-bold px-2 py-0.5 rounded shadow-lg">
                            {item.totalEpisodes ? `EP ${item.totalEpisodes}` : 'Ongoing'}
                        </div>

                        {/* Play Button Overlay */}
                        <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover/card:opacity-100 transition-all duration-300 transform scale-50 group-hover/card:scale-100">
                            <div className="w-12 h-12 bg-red-600 rounded-full flex items-center justify-center shadow-lg shadow-red-600/40 text-white">
                                <FaPlay className="pl-1 text-lg" />
                            </div>
                        </div>
                    </div>
                    
                    <div className="p-3">
                         <h3 className="text-sm md:text-base font-bold text-zinc-100 line-clamp-2 group-hover/card:text-red-500 transition-colors leading-tight">
                            {item.title}
                         </h3>
                         <div className="flex items-center gap-2 text-[11px] text-zinc-500 mt-1.5">
                            <span>2024</span>
                            <span className="w-1 h-1 bg-zinc-600 rounded-full" />
                            <span>Anime</span>
                         </div>
                    </div>
                </a>
            ))}
        </div>
    );
}
