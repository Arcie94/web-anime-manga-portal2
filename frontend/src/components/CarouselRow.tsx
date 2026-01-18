import { useRef, useState } from 'react';
import { Swiper, SwiperSlide } from 'swiper/react';
import { Autoplay, Navigation } from 'swiper/modules';
import { FaPlay, FaChevronLeft, FaChevronRight } from "react-icons/fa";
import { getAnimeImage } from '../lib/utils';

// Import Swiper styles
import 'swiper/css';
import 'swiper/css/navigation';

interface Props {
    title: string;
    items?: any[];
    endpoint?: string;
    type: 'anime' | 'manga';
    autoPlay?: boolean;
}

export default function CarouselRow({ title, items = [], type, autoPlay = true }: Props) {
    if (items.length === 0) return null;

    const [prevEl, setPrevEl] = useState<HTMLElement | null>(null);
    const [nextEl, setNextEl] = useState<HTMLElement | null>(null);

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
                {/* Navigation Buttons */}
                <button 
                    ref={setPrevEl}
                    className="absolute -left-4 md:-left-6 top-1/2 -translate-y-1/2 z-30 w-10 h-10 md:w-12 md:h-12 bg-zinc-900/90 hover:bg-red-600 text-white rounded-full shadow-xl flex items-center justify-center opacity-0 group-hover/slider:opacity-100 transition-all duration-300 disabled:opacity-0 disabled:cursor-not-allowed"
                >
                     <FaChevronLeft />
                </button>
                <button 
                    ref={setNextEl}
                    className="absolute -right-4 md:-right-6 top-1/2 -translate-y-1/2 z-30 w-10 h-10 md:w-12 md:h-12 bg-zinc-900/90 hover:bg-red-600 text-white rounded-full shadow-xl flex items-center justify-center opacity-0 group-hover/slider:opacity-100 transition-all duration-300 disabled:opacity-0 disabled:cursor-not-allowed"
                >
                     <FaChevronRight />
                </button>

                <Swiper
                    modules={[Autoplay, Navigation]}
                    spaceBetween={16}
                    slidesPerView={2.2}
                    loop={true}
                    autoplay={autoPlay ? {
                        delay: 3000,
                        disableOnInteraction: false,
                        pauseOnMouseEnter: true
                    } : false}
                    navigation={{
                        prevEl,
                        nextEl
                    }}
                    breakpoints={{
                        640: { slidesPerView: 3.2, spaceBetween: 16 },
                        768: { slidesPerView: 4.2, spaceBetween: 20 },
                        1024: { slidesPerView: 5.2, spaceBetween: 24 },
                        1280: { slidesPerView: 6.2, spaceBetween: 24 },
                    }}
                    className="w-full !overflow-visible"
                >
                    {items.map((item, idx) => (
                        <SwiperSlide key={idx} className="!h-auto">
                            <a 
                                href={`/${type}/${item.slug || item.animeId}`}
                                className="block h-full relative group/card"
                            >
                                <div className="aspect-[2/3] w-full relative rounded-xl overflow-hidden shadow-lg bg-zinc-900">
                                    <img 
                                        src={getAnimeImage(item)} 
                                        alt={item.title} 
                                        className="w-full h-full object-cover transition-transform duration-500 group-hover/card:scale-110"
                                        loading="lazy"
                                    />
                                    
                                    {/* Overlay */}
                                    <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/20 to-transparent opacity-60 group-hover/card:opacity-40 transition-opacity" />

                                    {/* Top Metadata */}
                                    <div className="absolute top-2 left-2 right-2 flex justify-between items-start">
                                        <span className="bg-zinc-900/80 backdrop-blur text-white text-[9px] font-bold px-1.5 py-0.5 rounded border border-white/10">
                                            {item.type?.toUpperCase() || (type === 'anime' ? 'ANIME' : 'MANGA')} 
                                        </span>
                                        <span className="bg-red-600 text-white text-[9px] font-bold px-1.5 py-0.5 rounded shadow-lg">
                                            {item.totalEpisodes ? `EP ${item.totalEpisodes}` : 'Ongoing'}
                                        </span>
                                    </div>

                                    {/* Play Button */}
                                    <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover/card:opacity-100 transition-all duration-300 transform scale-50 group-hover/card:scale-100">
                                        <div className="w-10 h-10 bg-red-600 rounded-full flex items-center justify-center shadow-lg shadow-red-600/40 text-white">
                                            <FaPlay className="pl-1 text-sm" />
                                        </div>
                                    </div>
                                </div>
                                
                                <div className="mt-2.5">
                                     <h3 className="text-sm font-semibold text-zinc-100 line-clamp-1 group-hover/card:text-red-500 transition-colors">
                                        {item.title}
                                     </h3>
                                     <div className="flex items-center gap-2 text-[10px] text-zinc-500 mt-0.5">
                                        {(item.releaseDate || item.time_ago || item.status) && (
                                            <>
                                                <span className="truncate max-w-[80px]">{item.releaseDate || item.time_ago || item.status}</span>
                                                <span className="w-1 h-1 bg-zinc-600 rounded-full flex-shrink-0" />
                                            </>
                                        )}
                                        <span>{type === 'anime' ? 'Anime' : 'Manga'}</span>
                                     </div>
                                </div>
                            </a>
                        </SwiperSlide>
                    ))}
                </Swiper>
            </div>
        </div>
    );
}
