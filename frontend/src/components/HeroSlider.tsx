import { Swiper, SwiperSlide } from 'swiper/react';
import { Autoplay, EffectFade, Pagination } from 'swiper/modules';
// Import Swiper styles
import 'swiper/css';
import 'swiper/css/effect-fade';
import 'swiper/css/pagination';

interface HeroSliderProps {
    items: {
        title: string;
        image: string;
        slug: string;
        type: 'anime' | 'manga';
        desc?: string;
    }[];
}

export default function HeroSlider({ items }: HeroSliderProps) {
  return (
    <div className="h-[50vh] min-h-[400px] w-full relative group">
      <Swiper
        modules={[Autoplay, EffectFade, Pagination]}
        effect="fade"
        autoplay={{ delay: 5000, disableOnInteraction: false }}
        pagination={{ clickable: true }}
        loop={true}
        className="h-full w-full"
      >
        {items.map((item, idx) => (
          <SwiperSlide key={idx} className="relative w-full h-full bg-gray-900">
             {/* Background Image with Gradient Overlay */}
            <div className="absolute inset-0 z-0">
                <img 
                    src={item.image} 
                    alt={item.title} 
                    className="w-full h-full object-cover opacity-80"
                    onError={(e) => {
                        (e.target as HTMLImageElement).src = 'https://placehold.co/1920x1080/1a1a1a/FFF?text=No+Image';
                    }}
                />
                <div className="absolute inset-0 bg-gradient-to-t from-gray-900 via-gray-900/40 to-transparent" />
            </div>

            {/* Content */}
            <div className="absolute bottom-0 left-0 right-0 p-4 md:p-8 z-10 flex flex-col items-start gap-2 md:gap-4 max-w-2xl px-6 pb-12">
                <span className="px-2 py-0.5 bg-pink-600 text-white text-[10px] md:text-xs font-bold rounded-sm uppercase tracking-wider shadow-sm">
                    Featured {item.type}
                </span>
                <h2 className="text-3xl md:text-5xl font-black text-white leading-tight drop-shadow-md line-clamp-2">
                    {item.title}
                </h2>
                {item.desc && (
                    <p className="text-gray-200 line-clamp-2 text-xs md:text-sm drop-shadow-sm hidden sm:block">
                        {item.desc}
                    </p>
                )}
                
                <div className="flex gap-2 mt-2 w-full">
                    <a 
                        href={`/${item.type}/watch/${item.slug}`} 
                        className="flex-1 sm:flex-none text-center bg-pink-600 text-white px-6 py-2 rounded-lg text-sm font-bold hover:bg-pink-700 transition-colors shadow-lg"
                    >
                        Watch Now
                    </a>
                    <a 
                        href={`/${item.type}/${item.slug}`} 
                        className="flex-1 sm:flex-none text-center bg-white/10 backdrop-blur-sm border border-white/20 text-white px-6 py-2 rounded-lg text-sm font-bold hover:bg-white/20 transition-colors"
                    >
                        Detail
                    </a>
                </div>
            </div>
          </SwiperSlide>
        ))}
      </Swiper>
    </div>
  );
}
