// Featured Manga for Hero Carousel
// These manga will always be displayed in the hero section on /manga page

export interface FeaturedManga {
  title: string;
  slug: string;
  image: string;
  poster: string;
  description: string;
  rating: string;
  year: string;
  type: "manga";
}

export const FEATURED_MANGA: FeaturedManga[] = [
  {
    title: "One Piece",
    slug: "komik-one-piece-indo",
    image: "https://thumbnail.komiku.org/uploads/manga/komik-one-piece-indo/manga_thumbnail-Komik-One-Piece.jpg",
    poster: "https://thumbnail.komiku.org/uploads/manga/komik-one-piece-indo/manga_thumbnail-Komik-One-Piece.jpg",
    description: "Monkey D. Luffy bermimpi menjadi Raja Bajak Laut dengan menemukan harta karun legendaris One Piece. Bersama krunya, ia berlayar melintasi Grand Line menghadapi berbagai musuh dan petualangan epik.",
    rating: "9.2",
    year: "1997",
    type: "manga",
  },
  {
    title: "Naruto",
    slug: "naruto",
    image: "https://thumbnail.komiku.org/uploads/manga/naruto/manga_thumbnail-Manga-Naruto.jpg",
    poster: "https://thumbnail.komiku.org/uploads/manga/naruto/manga_thumbnail-Manga-Naruto.jpg",
    description: "Kisah Naruto Uzumaki, ninja muda yang bercita-cita menjadi Hokage. Dengan kekuatan rubah berekor sembilan yang tersegel dalam tubuhnya, ia berjuang mendapatkan pengakuan dari desa.",
    rating: "8.5",
    year: "1999",
    type: "manga",
  },
  {
    title: "Boruto: Two Blue Vortex",
    slug: "boruto-two-blue-vortex",
    image: "https://thumbnail.komiku.org/uploads/manga/boruto-two-blue-vortex/manga_thumbnail-Head-Boruto-Two-Blue-Vortex.jpg",
    poster: "https://thumbnail.komiku.org/uploads/manga/boruto-two-blue-vortex/manga_thumbnail-Head-Boruto-Two-Blue-Vortex.jpg",
    description: "Kelanjutan dari era Naruto. Boruto, putra Naruto, menghadapi ancaman baru yang lebih berbahaya sambil mencari jalan ninjanya sendiri di tengah bayang-bayang ayahnya.",
    rating: "7.8",
    year: "2023",
    type: "manga",
  },
  {
    title: "Shingeki no Kyojin",
    slug: "shingeki-no-kyojin",
    image: "https://thumbnail.komiku.org/uploads/manga/shingeki-no-kyojin/manga_thumbnail-Komik-Shingeki-no-Kyojin.jpg",
    poster: "https://thumbnail.komiku.org/uploads/manga/shingeki-no-kyojin/manga_thumbnail-Komik-Shingeki-no-Kyojin.jpg",
    description: "Umat manusia hidup di balik tembok raksasa untuk melindungi diri dari Titan. Eren Yeager bersumpah membasmi semua Titan setelah ibunya dimakan oleh salah satu dari mereka.",
    rating: "9.0",
    year: "2009",
    type: "manga",
  },
  {
    title: "Kimetsu no Yaiba",
    slug: "kimetsu-no-yaiba-indonesia",
    image: "https://thumbnail.komiku.org/uploads/manga/kimetsu-no-yaiba-indonesia/manga_thumbnail-Komik-Kimetsu-no-Yaiba.jpg",
    poster: "https://thumbnail.komiku.org/uploads/manga/kimetsu-no-yaiba-indonesia/manga_thumbnail-Komik-Kimetsu-no-Yaiba.jpg",
    description: "Tanjiro Kamado menjadi pemburu iblis setelah keluarganya dibantai dan adiknya Nezuko berubah menjadi iblis. Ia berjuang mencari cara mengembalikan Nezuko sambil membasmi iblis",
    rating: "8.7",
    year: "2016",
    type: "manga",
  },
  {
    title: "Jujutsu Kaisen",
    slug: "jujutsu-kaisen-indo",
    image: "https://thumbnail.komiku.org/uploads/manga/jujutsu-kaisen-indo/manga_thumbnail-Manhua-Jujutsu-Kaisen.jpg",
    poster: "https://thumbnail.komiku.org/uploads/manga/jujutsu-kaisen-indo/manga_thumbnail-Manhua-Jujutsu-Kaisen.jpg",
    description: "Yuji Itadori bergabung dengan sekolah jujutsu setelah menelan jari kutukan Sukuna. Ia harus mengendalikan kekuatan kutukan sambil mencari dan menghancurkan sisa jari Sukuna.",
    rating: "8.6",
    year: "2018",
    type: "manga",
  },
  {
    title: "Dragon Ball Super",
    slug: "dragon-ball-super",
    image: "https://thumbnail.komiku.org/uploads/manga/dragon-ball-super/manga_thumbnail-Manga-Dragon-Ball-Super.jpg",
    poster: "https://thumbnail.komiku.org/uploads/manga/dragon-ball-super/manga_thumbnail-Manga-Dragon-Ball-Super.jpg",
    description: "Kelanjutan dari Dragon Ball Z. Goku dan teman-temannya menghadapi musuh dari alam semesta lain yang lebih kuat, termasuk dewa-dewa dan makhluk dari dimensi berbeda.",
    rating: "8.1",
    year: "2015",
    type: "manga",
  },
  {
    title: "Bleach",
    slug: "bleach",
    image: "https://thumbnail.komiku.org/uploads/manga/bleach/manga_thumbnail-Manga-Bleach.jpg",
    poster: "https://thumbnail.komiku.org/uploads/manga/bleach/manga_thumbnail-Manga-Bleach.jpg",
    description: "Ichigo Kurosaki mendapatkan kekuatan Shinigami dan harus melindungi manusia dari roh jahat sambil mengungkap rahasia dunia roh dan konspirasi di Soul Society.",
    rating: "8.2",
    year: "2001",
    type: "manga",
  },
];
