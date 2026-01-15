import { FaSearch, FaUser, FaPlayCircle, FaBookOpen } from "react-icons/fa";
import { BiBookBookmark } from "react-icons/bi";
import { useState, useEffect } from "react";

export default function BottomNav() {
  const [activePath, setActivePath] = useState("/");

  useEffect(() => {
    setActivePath(window.location.pathname);
  }, []);

  const menus = [
    { name: "Anime", icon: FaPlayCircle, href: "/" },
    { name: "Manga", icon: FaBookOpen, href: "/manga" },
    { name: "Search", icon: FaSearch, href: "/search" },
    { name: "My List", icon: BiBookBookmark, href: "/bookmarks" },
    { name: "Profile", icon: FaUser, href: "/login" },
  ];

  return (
    <div className="fixed bottom-0 left-0 right-0 bg-[#09090b] border-t border-white/5 pb-safe z-50 shadow-[0_-5px_10px_rgba(0,0,0,0.3)]">
      <div className="flex justify-around items-center h-16">
        {menus.map((menu) => {
          const isActive = activePath === menu.href || (menu.href !== "/" && activePath.startsWith(menu.href));
          return (
            <a
              key={menu.name}
              href={menu.href}
              className={`flex flex-col items-center justify-center w-full h-full transition-colors duration-200 ${
                isActive ? "text-red-600" : "text-zinc-500 hover:text-red-500"
              }`}
            >
              <menu.icon size={20} className="mb-1" />
              <span className="text-[10px] font-medium tracking-wide">{menu.name}</span>
            </a>
          );
        })}
      </div>
    </div>
  );
}
