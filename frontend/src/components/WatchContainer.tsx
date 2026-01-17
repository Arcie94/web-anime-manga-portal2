import React, { useState } from 'react';
import VideoPlayer from './VideoPlayer';
import { FaChevronLeft, FaChevronRight, FaServer } from 'react-icons/fa';
import { MdHighQuality } from 'react-icons/md';

interface StreamServer {
  title: string;
  serverId: string;
  href: string;
}

interface Quality {
  title: string;
  serverList: StreamServer[];
}

interface WatchContainerProps {
  sources: Record<string, string>; // Dynamic server names from backend
  qualities?: Quality[];
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

export default function WatchContainer({ sources, qualities = [], info, nav, episodeList = [] }: WatchContainerProps) {
  // Get server entries and determine initial source
  const serverEntries = Object.entries(sources).filter(([key]) => key !== 'default');
  const firstServerName = serverEntries[0]?.[0] || 'default';
  const firstServerUrl = sources.default || serverEntries[0]?.[1] || '';
  
  const [currentSrc, setCurrentSrc] = useState(firstServerUrl);
  const [activeServer, setActiveServer] = useState(firstServerName);
  const [selectedQuality, setSelectedQuality] = useState<string>('');
  const [selectedQualityServer, setSelectedQualityServer] = useState<string>('');


  const handleServerChange = (src: string, name: string) => {
    setCurrentSrc(src);
    setActiveServer(name);
    setSelectedQuality(''); // Reset quality when server changes
    setSelectedQualityServer('');
  };

  const handleQualitySelect = (qualityTitle: string) => {
    setSelectedQuality(qualityTitle);
    
    // Auto-select server if there's only one
    const quality = qualities.find(q => q.title === qualityTitle);
    if (quality && quality.serverList.length === 1) {
       const server = quality.serverList[0];
       handleQualityServerSelect(server, qualityTitle);
    } else {
       setSelectedQualityServer(''); // Reset server when quality changes if multiple options
       setActiveServer(''); // Clear main server selection
    }
  };

  const handleQualityServerSelect = (server: StreamServer, qualityTitle: string) => {
    let videoUrl = server.href;
    
    // If it's a relative path (starts with /), prefix with Otakudesu base URL
    if (videoUrl.startsWith('/')) {
      videoUrl = `https://otakudesu.cloud${videoUrl}`;
    }
    
    setCurrentSrc(videoUrl);
    setSelectedQuality(qualityTitle);
    setSelectedQualityServer(server.serverId);
    setActiveServer(''); // Clear main server selection
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
       <div className="flex flex-col gap-4 bg-zinc-900/50 p-4 rounded-xl border border-white/5">
          {/* Server Selector */}
          <div className="flex flex-wrap gap-2 items-center w-full">
             <span className="text-zinc-500 text-sm flex items-center gap-2 mr-2">
                <FaServer /> Servers:
             </span>
             {Object.entries(sources)
               .filter(([key]) => key !== 'default')
               .map(([serverName, url]) => (
                 <button
                   key={serverName}
                   onClick={() => handleServerChange(url, serverName)}
                   className={`px-3 py-1.5 rounded text-sm font-medium transition-colors ${
                     activeServer === serverName 
                       ? 'bg-red-600 text-white' 
                       : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700 hover:text-white'
                   }`}
                 >
                   {serverName}
                 </button>
               ))
             }
          </div>

          {/* Quality Selector (NEW) */}
          {qualities.length > 0 && (
             <div className="flex flex-col gap-2 w-full">
                <div className="flex flex-wrap gap-2 items-center w-full">
                   <span className="text-zinc-500 text-sm flex items-center gap-2 mr-2">
                      <MdHighQuality /> Quality:
                   </span>
                   {qualities.map((quality) => (
                      <button
                         key={quality.title}
                         onClick={() => handleQualitySelect(quality.title)}
                         className={`px-3 py-1.5 rounded text-sm font-medium transition-colors ${selectedQuality === quality.title ? 'bg-red-600 text-white' : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700 hover:text-white'}`}
                      >
                         {quality.title}
                      </button>
                   ))}
                </div>

                {/* Show servers for selected quality if multiple */}
                {selectedQuality && (() => {
                   const selectedQ = qualities.find(q => q.title === selectedQuality);
                   if (selectedQ && selectedQ.serverList.length > 1) {
                      return (
                         <div className="flex flex-wrap gap-2 items-center w-full pl-2">
                            <span className="text-zinc-500 text-xs mr-2">Choose Server:</span>
                            {selectedQ.serverList.map((server) => (
                               <button
                                  key={server.serverId}
                                  onClick={() => handleQualityServerSelect(server, selectedQuality)}
                                  className={`px-2 py-1 rounded text-xs font-medium transition-colors ${selectedQualityServer === server.serverId ? 'bg-red-600 text-white' : 'bg-zinc-700 text-zinc-400 hover:bg-zinc-600 hover:text-white'}`}
                               >
                                  {server.title}
                               </button>
                            ))}
                         </div>
                      );
                   }
                   return null;
                })()}
             </div>
          )}
       </div>

       {/* Episode Navigation */}
       <div className="flex gap-3 justify-center items-center bg-zinc-900/50 p-4 rounded-xl border border-white/5">
          {/* Previous Episode */}
          {nav.prevSlug ? (
             <a 
                href={`/anime/watch/${nav.prevSlug}`}
                className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-white rounded-lg transition-colors"
             >
                <FaChevronLeft />
                <span>Prev</span>
             </a>
          ) : (
             <button disabled className="flex items-center gap-2 px-4 py-2 bg-zinc-800/50 text-zinc-600 rounded-lg cursor-not-allowed">
                <FaChevronLeft />
                <span>Prev</span>
             </button>
          )}

          {/* Episode List Dropdown (if available) */}
          {episodeList && episodeList.length > 0 && (
             <div className="relative">
                <select
                  onChange={(e) => {
                     const selectedEp = e.target.value;
                     if (selectedEp) window.location.href = `/anime/watch/${selectedEp}`;
                  }}
                  className="px-4 py-2 bg-zinc-800 text-white rounded-lg appearance-none pr-10 cursor-pointer hover:bg-zinc-700 transition-colors"
                  defaultValue=""
                >
                  <option value="" disabled>Episode List</option>
                  {episodeList.map((ep: any, idx: number) => (
                     <option key={idx} value={ep.slug || ep.href?.split('/').pop()}>
                        {ep.title}
                     </option>
                  ))}
                </select>
                <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-white">
                   <svg className="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                      <path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z" />
                   </svg>
                </div>
             </div>
          )}

          {/* Next Episode */}
          {nav.nextSlug ? (
             <a 
                href={`/anime/watch/${nav.nextSlug}`}
                className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-white rounded-lg transition-colors"
             >
                <span>Next</span>
                <FaChevronRight />
             </a>
          ) : (
             <button disabled className="flex items-center gap-2 px-4 py-2 bg-zinc-800/50 text-zinc-600 rounded-lg cursor-not-allowed">
                <span>Next</span>
                <FaChevronRight />
             </button>
          )}
       </div>
    </div>
  );
}
