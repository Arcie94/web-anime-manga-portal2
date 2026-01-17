import { useState } from 'react';
import { FaPlay, FaInfoCircle, FaStar, FaBook } from "react-icons/fa";

interface HeroProps {
    item: {
        title: string;
        image: string; // Backdrop
        poster?: string; // Poster from API
        cover?: string; // Poster (vertical), usually we might need separate fields or just use image for both if API doesn't provide
        slug: string;
        type: 'anime' | 'manga';
        desc?: string;
        rating?: string;
        year?: string;
    };
}

export default function HeroBillboard({ item }: HeroProps) {
    const [imageLoaded, setImageLoaded] = useState(false);

    // Configurable fallback
    const backdropImage = item.image;
    // If we had a separate poster image, we'd use it here. 
    // For now, let's use the same image but assume we might want a different aspect ratio container
    // Or if the API provides 'cover' (vertical) vs 'image' (horizontal).
    const posterImage = item.poster || item.cover || item.image; 

    return (
        <div className="relative w-full min-h-[85vh] md:h-[85vh] md:min-h-[600px] overflow-hidden group">
            {/* Blurred Backdrop */}
            <div className="absolute inset-0 w-full h-full">
                <img 
                    src={backdropImage} 
                    alt="Backdrop"
                    className={`w-full h-full object-cover blur-xl scale-110 transition-opacity duration-1000 ${imageLoaded ? 'opacity-60' : 'opacity-0'}`}
                    onLoad={() => setImageLoaded(true)}
                    onError={() => setImageLoaded(true)}
                />
                <div className="absolute inset-0 bg-gradient-to-t from-zinc-950 via-zinc-950/50 to-zinc-950/30" />
                <div className="absolute inset-0 bg-gradient-to-r from-zinc-950 via-zinc-950/60 to-transparent" />
            </div>

            {/* Content Container */}
            <div className="relative z-10 w-full h-full container mx-auto px-4 md:px-12 flex items-start md:items-center pt-28 md:pt-0 pb-20 md:pb-0">
                <div className="flex flex-col md:flex-row items-center gap-6 md:gap-8 w-full">
                    
                    {/* Poster - Enhanced size for mobile prominence */}
                    <div className="w-56 sm:w-64 md:w-72 lg:w-[240px] aspect-[2/3] flex-shrink-0 shadow-2xl rounded-lg overflow-hidden border border-white/10 transform transition-all duration-300 mx-auto md:mx-0">
                         <img 
                            src={posterImage} 
                            alt={item.title} 
                            className={`w-full h-full object-cover`}
                            onLoad={() => setImageLoaded(true)}
                        />
                    </div>

                    {/* Info */}
                    <div className="flex-1 space-y-6 max-w-2xl text-center md:text-left mx-auto md:mx-0 pt-4 md:pt-0">
                        <div className="space-y-2">
                             <div className="flex items-center justify-center md:justify-start gap-3 text-xs md:text-sm font-semibold text-yellow-500 mb-2">
                                {(item.rating) && (
                                    <span className="flex items-center gap-1 bg-yellow-500/10 px-2 py-1 rounded">
                                        <FaStar /> {item.rating}
                                    </span>
                                )}
                                {item.year && <span className="text-zinc-400">{item.year}</span>}
                                {/* <span className="bg-zinc-800 text-zinc-300 px-2 py-1 rounded">13+</span> */}
                                <span className="text-red-500 font-bold uppercase">{item.type}</span>
                            </div>

                            <h1 className="text-4xl md:text-6xl font-black text-white leading-tight drop-shadow-2xl line-clamp-1">
                                {item.title}
                            </h1>
                        </div>

                        {item.desc && (
                            <p className="text-zinc-300 text-sm md:text-base line-clamp-2 leading-relaxed max-w-xl mx-auto md:mx-0">
                                {item.desc}
                            </p>
                        )}

                        <div className="flex flex-col sm:flex-row items-center justify-center md:justify-start gap-4 pt-4">
                            <a 
                                href={`/${item.type}/${item.slug}`} 
                                className="w-64 sm:w-auto flex items-center justify-center gap-2 bg-red-600 text-white px-8 py-3.5 rounded-full font-bold hover:bg-red-700 transition-all shadow-lg shadow-red-600/20 active:scale-95"
                            >
                                {item.type === 'manga' ? <FaBook className="text-sm" /> : <FaPlay className="text-sm" />}
                                {item.type === 'manga' ? 'Read Now' : 'Watch Now'}
                            </a>
                            <a 
                                href={`/${item.type}/${item.slug}`} 
                                className="w-64 sm:w-auto flex items-center justify-center gap-2 bg-white/10 backdrop-blur-md text-white border border-white/10 px-8 py-3.5 rounded-full font-bold hover:bg-white/20 transition-all active:scale-95"
                            >
                                <FaInfoCircle />
                                More Info
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
