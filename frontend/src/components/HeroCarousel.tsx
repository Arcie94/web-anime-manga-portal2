import { useState, useEffect } from 'react';
import { FaChevronLeft, FaChevronRight } from 'react-icons/fa';
import HeroBillboard from './HeroBillboard';
import type { FeaturedAnime } from '../data/featuredAnime';

interface HeroCarouselProps {
  featured: FeaturedAnime[];
  autoPlayInterval?: number;
}

export default function HeroCarousel({ 
  featured, 
  autoPlayInterval = 5000 
}: HeroCarouselProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isPaused, setIsPaused] = useState(false);

  // Auto-rotation logic
  useEffect(() => {
    if (!isPaused && featured.length > 1) {
      const timer = setInterval(() => {
        setCurrentIndex((prev) => (prev + 1) % featured.length);
      }, autoPlayInterval);
      
      return () => clearInterval(timer);
    }
  }, [currentIndex, isPaused, featured.length, autoPlayInterval]);

  // Navigation functions
  const goToNext = () => {
    setCurrentIndex((prev) => (prev + 1) % featured.length);
  };

  const goToPrevious = () => {
    setCurrentIndex((prev) => (prev - 1 + featured.length) % featured.length);
  };

  const goToSlide = (index: number) => {
    setCurrentIndex(index);
  };

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'ArrowLeft') goToPrevious();
      if (e.key === 'ArrowRight') goToNext();
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  if (featured.length === 0) {
    return null;
  }

  const currentAnime = featured[currentIndex];

  return (
    <div 
      className="relative w-full"
      onMouseEnter={() => setIsPaused(true)}
      onMouseLeave={() => setIsPaused(false)}
    >
      {/* Carousel Slides */}
      <div className="relative overflow-hidden">
        {featured.map((anime, index) => (
          <div
            key={anime.slug}
            className={`transition-opacity duration-1000 ${
              index === currentIndex 
                ? 'opacity-100 relative' 
                : 'opacity-0 absolute inset-0 pointer-events-none'
            }`}
          >
            <HeroBillboard 
              item={{
                title: anime.title,
                image: anime.poster, // Use poster as backdrop
                poster: anime.poster,
                cover: anime.poster,
                slug: anime.slug,
                type: anime.type,
                desc: anime.description,
                rating: anime.rating,
              }} 
            />
          </div>
        ))}
      </div>

      {/* Previous Button */}
      <button
        onClick={goToPrevious}
        className="absolute left-4 top-1/2 -translate-y-1/2 z-20 hidden md:flex items-center justify-center w-12 h-12 rounded-full bg-black/50 hover:bg-black/70 text-white border border-white/10 transition-all hover:scale-110 active:scale-95 backdrop-blur-sm"
        aria-label="Previous slide"
      >
        <FaChevronLeft className="text-xl" />
      </button>

      {/* Next Button */}
      <button
        onClick={goToNext}
        className="absolute right-4 top-1/2 -translate-y-1/2 z-20 hidden md:flex items-center justify-center w-12 h-12 rounded-full bg-black/50 hover:bg-black/70 text-white border border-white/10 transition-all hover:scale-110 active:scale-95 backdrop-blur-sm"
        aria-label="Next slide"
      >
        <FaChevronRight className="text-xl" />
      </button>


    </div>
  );
}
