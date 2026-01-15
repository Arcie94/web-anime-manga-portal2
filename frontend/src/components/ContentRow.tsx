import { useRef, useEffect, useState } from 'react';
import { FaChevronLeft, FaChevronRight, FaPlay } from "react-icons/fa";
import SkeletonRow from "./skeletons/SkeletonRow";
import { apiFetch } from '../lib/api';

interface Props {
    title: string;
    items?: any[];
    endpoint?: string;
    type: 'anime' | 'manga';
    dataKey?: string; // Optional key to extract from response (e.g., "ongoing", "completed")
}

export default function ContentRow({ title, items: initialItems = [], endpoint, type, dataKey }: Props) {
    const rowRef = useRef<HTMLDivElement>(null);
    const [items, setItems] = useState<any[]>(initialItems);
    const [loading, setLoading] = useState(!!endpoint && initialItems.length === 0);
    const [error, setError] = useState(false);

    useEffect(() => {
        if (!endpoint || initialItems.length > 0) return;

        const fetchData = async () => {
            try {
                const res = await fetch(`/api${endpoint}`);
                if (!res.ok) throw new Error('Network response was not ok');
                const json = await res.json();
                
                // Extract data based on dataKey if provided
                let extractedData = json.data || [];
                if (dataKey && json.data && json.data[dataKey]) {
                    extractedData = json.data[dataKey];
                } else if (dataKey && json[dataKey]) {
                    extractedData = json[dataKey];
                }
                
                setItems(extractedData);
            } catch (err) {
                console.error("Failed to fetch row data", err);
                setError(true);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [endpoint, dataKey]);

    const scroll = (direction: 'left' | 'right') => {
        if (rowRef.current) {
            const { scrollLeft, clientWidth } = rowRef.current;
            const scrollTo = direction === 'left' ? scrollLeft - clientWidth : scrollLeft + clientWidth;
            rowRef.current.scrollTo({ left: scrollTo, behavior: 'smooth' });
        }
    };

    if (loading) {
        return <SkeletonRow />;
    }

    if (error || items.length === 0) {
        // Don't render empty rows
        return null; 
    }

    return (
        <div className="py-6 px-4 md:px-12 group relative">
            <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl md:text-2xl font-bold text-white flex items-center gap-3">
                    <div className="w-1.5 h-6 bg-red-600 rounded-full"></div>
                    {title}
                </h2>
                <a href={`/search?type=${type}`} className="text-xs font-semibold text-zinc-500 hover:text-red-500 transition-colors uppercase tracking-wider">
                    View All
                </a>
            </div>

            <div className="relative group/slider">
                {/* Left Arrow */}
                <button 
                    className="absolute -left-4 md:-left-6 top-1/2 -translate-y-1/2 z-30 w-10 h-10 md:w-12 md:h-12 bg-zinc-900/90 hover:bg-red-600 text-white rounded-full shadow-xl flex items-center justify-center opacity-0 group-hover/slider:opacity-100 transition-all duration-300 disabled:opacity-0"
                    onClick={() => scroll('left')}
                >
                     <FaChevronLeft />
                </button>

                {/* Slider Container */}
                <div 
                    ref={rowRef}
                    className="flex gap-4 overflow-x-auto scrollbar-hide scroll-smooth snap-x pb-4"
                    style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}
                >
                    {items.map((item, idx) => (
                        <a 
                            key={idx} 
                            href={`/${type}/${item.slug || item.animeId}`} // Support both data structures
                            className="flex-none w-[140px] md:w-[180px] snap-start relative group/card"
                        >
                            <div className="aspect-[2/3] w-full relative rounded-xl overflow-hidden shadow-lg bg-zinc-900">
                                <img 
                                    src={item.poster || item.cover || item.image} 
                                    alt={item.title} 
                                    className="w-full h-full object-cover transition-transform duration-500 group-hover/card:scale-110"
                                    loading="lazy"
                                    onError={(e) => {
                                        (e.target as HTMLImageElement).src = 'https://placehold.co/200x300/222/999?text=No+Image';
                                    }}
                                />
                                
                                {/* Overlay Gradient */}
                                <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/20 to-transparent opacity-60 group-hover/card:opacity-40 transition-opacity" />

                                {/* Badge Top Left */}
                                <div className="absolute top-2 left-2 bg-zinc-900/80 backdrop-blur text-white text-[10px] font-bold px-2 py-0.5 rounded border border-white/10">
                                    {type === 'anime' ? 'TV' : 'MANGA'} 
                                </div>

                                {/* Badge Top Right (Episode) */}
                                <div className="absolute top-2 right-2 bg-red-600 text-white text-[10px] font-bold px-2 py-0.5 rounded shadow-lg">
                                    {item.totalEpisodes ? `EP ${item.totalEpisodes}` : 'Ongoing'}
                                </div>

                                {/* Play Button Overlay */}
                                <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover/card:opacity-100 transition-all duration-300 transform scale-50 group-hover/card:scale-100">
                                    <div className="w-10 h-10 md:w-12 md:h-12 bg-red-600 rounded-full flex items-center justify-center shadow-lg shadow-red-600/40 text-white">
                                        <FaPlay className="pl-1" />
                                    </div>
                                </div>
                            </div>
                            
                            <div className="mt-2.5">
                                 <h3 className="text-sm font-semibold text-zinc-100 line-clamp-1 group-hover/card:text-red-500 transition-colors">
                                    {item.title}
                                 </h3>
                                 <div className="flex items-center gap-2 text-[11px] text-zinc-500 mt-0.5">
                                    <span>2024</span>
                                    <span className="w-1 h-1 bg-zinc-600 rounded-full" />
                                    <span>{type === 'anime' ? 'Anime' : 'Manga'}</span>
                                 </div>
                            </div>
                        </a>
                    ))}
                </div>

                 {/* Right Arrow */}
                 <button 
                    className="absolute -right-4 md:-right-6 top-1/2 -translate-y-1/2 z-30 w-10 h-10 md:w-12 md:h-12 bg-zinc-900/90 hover:bg-red-600 text-white rounded-full shadow-xl flex items-center justify-center opacity-0 group-hover/slider:opacity-100 transition-all duration-300"
                    onClick={() => scroll('right')}
                >
                     <FaChevronRight />
                </button>
            </div>
        </div>
    );
}
