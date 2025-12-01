// types/index.ts
export interface Song {
  id: number;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number;
  file_path: string;
  rating: number;
  play_count: number;
  last_played: string;
  created_at: string;
  explicit: boolean;
  clean: boolean;
}

export interface Playlist {
  id: number;
  name: string;
  description: string;
  songs: Song[];
  created_at: string;
}

export interface PlayerState {
  is_playing: boolean;
  current_song: Song | null;
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

export interface UploadProgress {
  id: string;
  file: File;
  name: string;
  size: number;
  status: "pending" | "uploading" | "completed" | "error";
  progress: number;
  error?: string;
}

export interface VisualizerMode {
  id: "bars" | "wave" | "particles" | "circular" | "spiral" | "waveform";
  name: string;
}
