// types/index.ts
export interface Song {
  id: number;
  title: string;
  artist: string;
  album?: string;
  genre?: string;
  year?: number;
  duration: number; // seconds
  file_path: string; // /files/<filename>
  file_name?: string; // internal filename
  orig_name?: string;
  rating: number;
  play_count: number;
  last_played?: string | null;
  created_at?: string;
  explicit?: boolean;
  clean?: boolean;
  size?: number;
}

export interface Playlist {
  id: number;
  name: string;
  description?: string;
  songs: Song[];
  created_at?: string;
}

export interface PlayerState {
  is_playing: boolean;
  current_song?: Song | null;
  current_time: number;
  duration: number;
  volume: number;
  speed: number;
}

export interface LibraryFilters {
  search: string;
  sortBy:
    | "title"
    | "artist"
    | "album"
    | "rating"
    | "play_count"
    | "last_played"
    | "duration"
    | "year";
  genre: string;
  rating: number | null;
}
