import React, { useState, useEffect, useRef } from 'react';
import { apiFetch } from '../lib/api';

interface ChapterData {
    slug: string;
    title: string;
    images: string[];
    nextSlug?: string;
    prevSlug?: string;
}

interface MangaReaderProps {
    initialSlug: string;
    initialTitle: string;
    initialImages: string[];
    initialNextSlug?: string;
    initialPrevSlug?: string;
}

export default function MangaReader({ 
    initialSlug, 
    initialTitle, 
    initialImages,
    initialNextSlug 
}: MangaReaderProps) {
    const [chapters, setChapters] = useState<ChapterData[]>([
        {
            slug: initialSlug,
            title: initialTitle,
            images: initialImages,
            nextSlug: initialNextSlug
        }
    ]);
    const [loading, setLoading] = useState(false);
    const observerTarget = useRef<HTMLDivElement>(null);

    // Optimized Image Loader to prevent lag
    const loadNextChapter = async (nextSlug: string | undefined) => {
        if (!nextSlug || loading) return;

        setLoading(true);
        try {
            console.log("Fetching next chapter:", nextSlug);
            const res = await apiFetch<any>(`/manga/chapter/${nextSlug}`);
            
            // Format slug as title if API title is missing or generic
            // e.g. "one-piece-chapter-1163" -> "One Piece Chapter 1163"
            const formattedTitle = nextSlug
                .split('-')
                .map(word => word.charAt(0).toUpperCase() + word.slice(1))
                .join(' ');

            const newChapter: ChapterData = {
                slug: nextSlug,
                // Prefer formatted slug to ensure Series Name is present
                // (API sometimes returns just "Chapter X")
                title: res.title && res.title.length > 10 ? res.title : formattedTitle, 
                images: res.images || [],
                nextSlug: res.nextSlug, // Navigation from API
            };

            setChapters(prev => [...prev, newChapter]);
            
            // Note: We do NOT update URL or Header here anymore.
            // The Scroll Observer will handle it when the user actually scrolls into the new chapter.

        } catch (error) {
            console.error("Failed to load next chapter", error);
        } finally {
            setLoading(false);
        }
    };

    // Observer 1: Infinite Scroll Trigger (Bottom of list)
    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting) {
                    const lastChapter = chapters[chapters.length - 1];
                    if (lastChapter.nextSlug) {
                        loadNextChapter(lastChapter.nextSlug);
                    }
                }
            },
            { threshold: 0.1, rootMargin: "400px" } 
        );

        if (observerTarget.current) {
            observer.observe(observerTarget.current);
        }

        return () => observer.disconnect();
    }, [chapters, loading]);

    // Observer 2: Active Chapter Tracking (Update Header & URL)
    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    // When a chapter crosses the center line of the viewport
                    if (entry.isIntersecting) {
                        const slug = entry.target.id.replace('chapter-', '');
                        
                        // Debounce/Check if it's actually the dominant one?
                        // With the tight rootMargin, usually only one intersects at a time.
                        const chapter = chapters.find((c) => c.slug === slug);
                        
                        if (chapter) {
                            // 1. Update Sticky Header
                            const headerTitle = document.getElementById('chapter-header-title');
                            if (headerTitle) {
                                headerTitle.textContent = chapter.title;
                            }

                            // 2. Update URL silently 
                            if (!window.location.pathname.includes(slug)) {
                                window.history.replaceState(null, "", `/manga/read/${slug}`);
                            }
                        }
                    }
                });
            },
            { 
                // Detection Line: A thin strip in the middle of the viewport
                rootMargin: "-45% 0px -45% 0px", 
                threshold: 0 
            }
        );

        // Observe all chapter containers
        chapters.forEach((chapter) => {
            const el = document.getElementById(`chapter-${chapter.slug}`);
            if (el) observer.observe(el);
        });

        return () => observer.disconnect();
    }, [chapters]);

    return (
        <div className="flex flex-col space-y-8 min-h-screen pb-20">
            {chapters.map((chapter, index) => (
                <div key={chapter.slug} className="chapter-container" id={`chapter-${chapter.slug}`}>
                    {/* Chapter Divider (Only for subsequent chapters) */}
                    {index > 0 && (
                        <div className="py-8 text-center">
                            <div className="inline-block bg-gray-800 text-gray-300 px-6 py-2 rounded-full text-sm font-semibold shadow-lg border border-gray-700">
                                Reading: {chapter.title}
                            </div>
                        </div>
                    )}

                    {/* Images */}
                    <div className="space-y-0"> 
                        {/* space-y-0 handles seamless vertical aspect, usually images stick together */}
                        {chapter.images.map((img, idx) => (
                            <img
                                key={`${chapter.slug}-${idx}`}
                                src={img}
                                className="max-w-full mx-auto shadow-lg md:rounded-lg"
                                loading="lazy"
                                alt={`${chapter.title} - Page ${idx + 1}`}
                            />
                        ))}
                    </div>
                </div>
            ))}

            {/* Loading Indicator / Sentinel */}
            <div ref={observerTarget} className="h-24 flex items-center justify-center">
                {loading ? (
                    <div className="flex flex-col items-center gap-2 text-gray-400">
                         <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-red-500"></div>
                         <span>Loading Next Chapter...</span>
                    </div>
                ) : (
                   chapters[chapters.length - 1].nextSlug ? (
                       <p className="text-gray-600 animate-pulse">Scroll for next chapter</p>
                   ) : (
                       <p className="text-gray-500 font-medium">End of Series (Latest Chapter)</p>
                   )
                )}
            </div>
        </div>
    );
}
