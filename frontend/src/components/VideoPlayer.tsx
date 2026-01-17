import React, { useState, useEffect } from 'react';
import ReactPlayer from 'react-player';
import { FaPlay, FaSpinner, FaExclamationTriangle } from "react-icons/fa";

interface VideoPlayerProps {
    src: string;
    poster?: string;
    title?: string;
}

export default function VideoPlayer({ src, poster, title }: VideoPlayerProps) {
    const [hasMounted, setHasMounted] = useState(false);
    const [error, setError] = useState(false);
    const [loading, setLoading] = useState(true);

    // Prevent hydration issues
    useEffect(() => {
        setHasMounted(true);
    }, []);

    if (!hasMounted) {
        return (
            <div className="w-full aspect-video bg-zinc-900 rounded-xl animate-pulse flex items-center justify-center">
                <FaSpinner className="animate-spin text-4xl text-zinc-700" />
            </div>
        );
    }

    if (!src) {
        return (
            <div className="w-full aspect-video bg-zinc-900 rounded-xl flex flex-col items-center justify-center text-zinc-500 gap-4">
                <FaExclamationTriangle className="text-4xl text-yellow-600" />
                <p>Maaf, video tidak tersedia.</p>
            </div>
        );
    }

    return (
        <div className="w-full aspect-video bg-black rounded-xl overflow-hidden shadow-2xl relative group">
            {/* Player Wrapper */}
            <div className="absolute inset-0">
                {/* Fallback to Iframe for Embed URLs */}
                {/* Oploverz player URLs (acefile, filedon, akirabox) need iframe */}
                {/* Oploverz player URLs (acefile, filedon, akirabox) need iframe */
                /* Also support standard Otakudesu .php streams */
                (src.includes('.php') || 
                  src.includes('/anime/server/') ||
                  src.includes('acefile.co') ||
                  src.includes('filedon.co') ||
                  src.includes('akirabox') ||
                  src.includes('gdplayer.to') ||
                  src.includes('anime-indo.lol') ||
                  src.includes('my.mail.ru') ||
                  src.includes('sankavollerei.com')) ? (
                    <iframe 
                        src={src} 
                        className="w-full h-full border-0" 
                        allowFullScreen 
                        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                    />
                ) : (
                    <ReactPlayer
                        url={src}
                        width="100%"
                        height="100%"
                        controls
                        playing
                        light={poster} 
                        playIcon={
                            <div className="bg-red-600 w-20 h-20 rounded-full flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform">
                                <FaPlay className="text-white text-2xl ml-1" />
                            </div>
                        }
                        onReady={() => setLoading(false)}
                        onError={() => {
                            setError(true);
                            setLoading(false);
                        }}
                        config={{
                            file: {
                                attributes: {
                                    crossOrigin: 'anonymous', 
                                }
                            }
                        }}
                    />
                )}
            </div>

            {/* Error Overlay */}
            {error && (
                <div className="absolute inset-0 bg-black/80 flex flex-col items-center justify-center text-white z-20">
                    <FaExclamationTriangle className="text-4xl text-red-500 mb-2" />
                    <p>Gagal memutar video.</p>
                    <p className="text-sm text-zinc-400 mt-2">Format tidak didukung atau server error.</p>
                </div>
            )}
        </div>
    );
}
