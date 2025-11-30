export interface Song {
  id: number;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number;
  filePath: string;
  rating: number;
  playCount: number;
  lastPlayed: string;
  createdAt: string;
  explicit: boolean;
  clean: boolean;
}

export interface Playlist {
  id: number;
  name: string;
  description: string;
  songs: Song[];
  createdAt: string;
}

export interface PlayerState {
  isPlaying: boolean;
  currentSong: Song | null;
  currentTime: number;
  duration: number;
  volume: number;
  speed: number;
}
