import { useState, useEffect } from 'react';

const categories = [
    "All", "Latest", "Trending", "Action", "Romance", "Comedy", 
    "Drama", "Fantasy", "Sci-Fi", "Thriller", "Horror", "Mystery"
];

export default function CategoryRow() {
    const [activeCat, setActiveCat] = useState('All');

    useEffect(() => {
        // Simple logic to highlight active category based on URL
        const path = window.location.pathname;
        if (path === '/') {
            setActiveCat('All');
        } else if (path.startsWith('/category/')) {
            const catSlug = path.split('/category/')[1];
            // Capitalize first letter to match array
            const formattedCat = catSlug.charAt(0).toUpperCase() + catSlug.slice(1);
            setActiveCat(formattedCat);
        }
    }, []);

    const getLink = (cat: string) => {
        if (cat === "All") return "/";
        return `/category/${cat.toLowerCase()}`;
    };

    return (
        <div className="sticky top-16 md:top-20 z-40 bg-zinc-950/95 backdrop-blur border-b border-white/5 space-y-2 py-3 shadow-xl">
            {/* Category Pills */}
            <div className="w-full overflow-x-auto scrollbar-hide px-4 md:px-12 pb-2">
                <div className="flex gap-2.5">
                    {categories.map((cat, idx) => (
                        <a 
                            key={idx}
                            href={getLink(cat)}
                            className={`
                                px-4 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition-all duration-300 border block
                                ${activeCat === cat
                                    ? 'bg-white text-black border-white shadow-lg scale-105' 
                                    : 'bg-zinc-900/50 text-zinc-400 border-zinc-800 hover:border-zinc-600 hover:text-white'
                                }
                            `}
                        >
                            {cat}
                        </a>
                    ))}
                </div>
            </div>
        </div>
    );
}
