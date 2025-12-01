export interface Song {
  id: number;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number;
  rating: number;
  playCount: number;
  explicit: boolean;
  clean: boolean;
  trackNumber: number;
  filePath: string;
  lastPlayed?: string;
  fileSize?: number;
  bitrate?: number;
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
