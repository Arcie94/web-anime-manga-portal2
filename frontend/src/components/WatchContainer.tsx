import React, { useState } from 'react';
import VideoPlayer from './VideoPlayer';
import { FaChevronLeft, FaChevronRight, FaServer } from 'react-icons/fa';

interface WatchContainerProps {
  sources: {
    default: string;
    backup: string;
    direct: string;
  };
  info: {
    animeTitle: string;
    episodeTitle: string;
    episodeNumber: string;
  };
  nav: {
    prevSlug?: string;
    nextSlug?: string;
  };
  episodeList?: any[];
}

export default function WatchContainer({ sources, info, nav, episodeList = [] }: WatchContainerProps) {
  // Determine initial source (prioritize default -> backup -> direct)
  const [currentSrc, setCurrentSrc] = useState(sources.default || sources.backup || sources.direct);
  const [activeServer, setActiveServer] = useState(sources.default ? 'Server 1' : sources.backup ? 'Server 2' : 'Server 3');

  const handleServerChange = (src: string, name: string) => {
    setCurrentSrc(src);
    setActiveServer(name);
  };

  return (
    <div className="space-y-6">
       {/* Player */}
       <VideoPlayer src={currentSrc} title={`Episode ${info.episodeNumber}`} />

       {/* Title Header (Moved Below Video) */}
       <div>
          <h1 className="text-xl md:text-3xl font-bold text-white mb-1">{info.animeTitle}</h1>
          <div className="flex items-center gap-3">
            <p className="text-zinc-400 text-sm md:text-base line-clamp-1">{info.episodeTitle}</p>
          </div>
       </div>

       {/* Controls */}
       <div className="flex flex-col md:flex-row gap-4 justify-between items-start md:items-center bg-zinc-900/50 p-4 rounded-xl border border-white/5">
          {/* Server Selector */}
          <div className="flex flex-wrap gap-2 items-center w-full md:w-auto">
             <span className="text-zinc-500 text-sm flex items-center gap-2 mr-2">
                <FaServer /> Servers:
             </span>
             {sources.default && (
                <button
                   onClick={() => handleServerChange(sources.default, 'Server 1')}
                   className={`px-3 py-1.5 rounded text-sm font-medium transition-colors ${activeServer === 'Server 1' ? 'bg-red-600 text-white' : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700 hover:text-white'}`}
                >
                   Server 1
                </button>
             )}
             {sources.backup && (
                <button
                   onClick={() => handleServerChange(sources.backup, 'Server 2')}
                   className={`px-3 py-1.5 rounded text-sm font-medium transition-colors ${activeServer === 'Server 2' ? 'bg-red-600 text-white' : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700 hover:text-white'}`}
                >
                   Server 2
                </button>
             )}
             {sources.direct && (
                <button
                   onClick={() => handleServerChange(sources.direct, 'Server 3')}
                   className={`px-3 py-1.5 rounded text-sm font-medium transition-colors ${activeServer === 'Server 3' ? 'bg-red-600 text-white' : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700 hover:text-white'}`}
                >
                   Server 3
                </button>
             )}
          </div>

          {/* Navigation */}
          <div className="flex gap-2 w-full md:w-auto border-t md:border-t-0 border-white/5 pt-4 md:pt-0 mt-2 md:mt-0">
             {nav.nextSlug ? (
                 <a
                    href={`/anime/watch/${nav.nextSlug}`}
                    className="flex-1 md:flex-none flex items-center justify-center gap-2 px-6 py-3 rounded-lg font-bold bg-red-600 text-white hover:bg-red-700 transition-all shadow-lg hover:shadow-red-600/20"
                 >
                    Next
                 </a>
             ) : (
                <button disabled className="flex-1 md:flex-none flex items-center justify-center gap-2 px-6 py-3 rounded-lg font-medium bg-zinc-900 text-zinc-600 cursor-not-allowed opacity-50 border border-white/5">
                    Next
                </button>
             )}

             {nav.prevSlug ? (
                <a
                    href={`/anime/watch/${nav.prevSlug}`}
                    className="flex-1 md:flex-none flex items-center justify-center gap-2 px-6 py-3 rounded-lg font-bold bg-zinc-800 hover:bg-zinc-700 text-white transition-all border border-white/10 hover:border-red-500/50"
                >
                    Prev
                </a>
             ) : (
                <button disabled className="flex-1 md:flex-none flex items-center justify-center gap-2 px-6 py-3 rounded-lg font-medium bg-zinc-900 text-zinc-600 cursor-not-allowed opacity-50 border border-white/5">
                    Prev
                </button>
             )}
          </div>
       </div>

       {/* Episode Queue / List */}
       <div className="space-y-4 pt-4 border-t border-white/10">
          <div className="flex items-center justify-between">
             <h2 className="text-xl font-bold text-white">More Episodes</h2>
             <span className="text-sm text-zinc-400">{episodeList.length} Available</span>
          </div>
          
          {episodeList.length === 0 ? (
             <p className="text-zinc-500 text-sm">No episodes found.</p>
          ) : (
             /* Horizontal Scroll Container */
             <div className="flex gap-3 overflow-x-auto pb-4 scrollbar-hide snap-x">
                {episodeList.map((ep, idx) => {
                   // Try to extract episode number from episodeId first, then eps, then slug
                   let epNum = ep.eps;
                   if (!epNum && ep.episodeId) {
                     const match = ep.episodeId.match(/episode-(\d+)/);
                     epNum = match ? match[1] : null;
                   }
                   if (!epNum && ep.slug) {
                     const match = ep.slug.match(/episode-(\d+)/);
                     epNum = match ? match[1] : null;
                   }
                   epNum = epNum || (episodeList.length - idx).toString(); // Fallback to reverse index
                   
                   const isActive = info.episodeNumber === epNum.toString();
                   
                   return (
                      <a
                         key={idx}
                         href={`/anime/watch/${ep.episodeId || ep.slug}`}
                         className={`flex-shrink-0 snap-start w-32 md:w-40 bg-zinc-900/50 hover:bg-zinc-800 border ${isActive ? 'border-red-600 ring-1 ring-red-600' : 'border-white/5 hover:border-white/20'} rounded-lg p-3 group transition-all`}
                      >
                         <div className="aspect-video bg-zinc-950 rounded mb-2 flex items-center justify-center relative overflow-hidden">
                            <span className={`text-lg font-bold ${isActive ? 'text-red-500' : 'text-zinc-500 group-hover:text-white'}`}>
                               {epNum}
                            </span>
                            {/* Visual indicator for active */}
                            {isActive && <div className="absolute inset-0 bg-red-600/10" />}
                         </div>
                         <p className="text-xs text-zinc-400 line-clamp-2 leading-tight group-hover:text-zinc-200">
                            {ep.title || `Episode ${epNum}`}
                         </p>
                      </a>
                   );
                })}
             </div>
          )}
       </div>
    </div>
  );
}
