// stores/library.ts - Extended version
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { Song, LibraryFilters, Playlist } from "../types";

export const useLibraryStore = defineStore("library", () => {
  // State
  const songs = ref<Song[]>([]);
  const playlists = ref<Playlist[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const selectedPlaylistId = ref<number | null>(null);

  const filters = ref<LibraryFilters>({
    search: "",
    sortBy: "title",
    genre: "",
    rating: null,
  });

  // Computed
  const filteredSongs = computed(() => {
    let result = [...songs.value];

    // Apply search filter
    if (filters.value.search) {
      const search = filters.value.search.toLowerCase();
      result = result.filter(
        (song) =>
          song.title.toLowerCase().includes(search) ||
          song.artist.toLowerCase().includes(search) ||
          song.album.toLowerCase().includes(search)
      );
    }

    // Apply genre filter
    if (filters.value.genre) {
      result = result.filter((song) => song.genre === filters.value.genre);
    }

    // Apply rating filter
    if (filters.value.rating !== null) {
      result = result.filter((song) => song.rating >= filters.value.rating!);
    }

    // Apply sorting
    result.sort((a, b) => {
      switch (filters.value.sortBy) {
        case "title":
          return a.title.localeCompare(b.title);
        case "artist":
          return a.artist.localeCompare(b.artist);
        case "album":
          return a.album.localeCompare(b.album);
        case "rating":
          return (b.rating || 0) - (a.rating || 0);
        case "play_count":
          return (b.play_count || 0) - (a.play_count || 0);
        case "last_played":
          return (
            new Date(b.last_played || 0).getTime() -
            new Date(b.last_played || 0).getTime()
          );
        default:
          return a.title.localeCompare(b.title);
      }
    });

    return result;
  });

  const genres = computed(() => {
    const genreSet = new Set(
      songs.value.map((song) => song.genre).filter(Boolean)
    );
    return Array.from(genreSet).sort();
  });

  // Getters for external access
  const searchQuery = computed({
    get: () => filters.value.search,
    set: (value: string) => {
      filters.value.search = value;
    },
  });

  // Actions
  const updatePlaylist = async (id: number, updates: Partial<Playlist>) => {
    try {
      const response = await fetch(
        `http://localhost:1323/api/playlists/${id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(updates),
        }
      );

      if (!response.ok) throw new Error("Failed to update playlist");

      const updatedPlaylist = await response.json();
      const index = playlists.value.findIndex((p) => p.id === id);
      if (index !== -1) {
        playlists.value[index] = updatedPlaylist;
      }
    } catch (err) {
      console.error("Error updating playlist:", err);
      throw err;
    }
  };

  const fetchSongs = async () => {
    try {
      const res = await fetch("http://localhost:1323/api/library/songs");
      if (!res.ok) throw new Error("Failed to fetch songs");
      const data: Song[] = await res.json();
      songs.value = data;
    } catch (err) {
      console.error("fetchSongs error:", err);
    }
  };

  const fetchPlaylists = async () => {
    try {
      const res = await fetch("http://localhost:1323/api/library/playlists");
      if (!res.ok) throw new Error("Failed to fetch playlists");
      const data: Playlist[] = await res.json();
      playlists.value = data;
    } catch (err) {
      console.error("fetchPlaylists error:", err);
    }
  };

  const createPlaylist = (p: Playlist) => {
    playlists.value.push(p);
  };

  const deletePlaylist = (id: number) => {
    playlists.value = playlists.value.filter((pl) => pl.id !== id);
    if (selectedPlaylistId.value === id) selectedPlaylistId.value = null;
  };

  const addSongToPlaylist = async (playlistId: number, songId: number) => {
    try {
      const response = await fetch(
        `http://localhost:1323/api/playlists/${playlistId}/songs/${songId}`,
        {
          method: "POST",
        }
      );

      if (!response.ok) throw new Error("Failed to add song to playlist");

      // Refresh playlists
      await fetchPlaylists();
    } catch (err) {
      console.error("Error adding song to playlist:", err);
      throw err;
    }
  };

  const removeSongFromPlaylist = async (playlistId: number, songId: number) => {
    try {
      const response = await fetch(
        `http://localhost:1323/api/playlists/${playlistId}/songs/${songId}`,
        {
          method: "DELETE",
        }
      );

      if (!response.ok) throw new Error("Failed to remove song from playlist");

      // Refresh playlists
      await fetchPlaylists();
    } catch (err) {
      console.error("Error removing song from playlist:", err);
      throw err;
    }
  };

  const updateSong = async (id: number, updates: Partial<Song>) => {
    try {
      const response = await fetch(`http://localhost:1323/api/songs/${id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(updates),
      });

      if (!response.ok) throw new Error("Failed to update song");

      // Update local state
      const index = songs.value.findIndex((song) => song.id === id);
      if (index !== -1) {
        songs.value[index] = { ...songs.value[index], ...updates };
      }

      // Also update in playlists
      playlists.value.forEach((playlist) => {
        const songIndex = playlist.songs.findIndex((song) => song.id === id);
        if (songIndex !== -1) {
          playlist.songs[songIndex] = {
            ...playlist.songs[songIndex],
            ...updates,
          };
        }
      });
    } catch (err) {
      console.error("Error updating song:", err);
      throw err;
    }
  };

  const deleteSong = async (id: number) => {
    try {
      const response = await fetch(`http://localhost:1323/api/songs/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) throw new Error("Failed to delete song");

      // Remove from local state
      songs.value = songs.value.filter((song) => song.id !== id);

      // Remove from playlists
      playlists.value.forEach((playlist) => {
        playlist.songs = playlist.songs.filter((song) => song.id !== id);
      });
    } catch (err) {
      console.error("Error deleting song:", err);
      throw err;
    }
  };

  const uploadSong = async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch("http://localhost:1323/api/songs", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Upload failed");
      }

      const data = await response.json();

      // Add to local state
      if (data.song) {
        songs.value.push(data.song);
      }

      return data;
    } catch (err) {
      console.error("Error uploading song:", err);
      throw err;
    }
  };

  const scanForDuplicates = () => {
    const duplicates: {
      title: string;
      artist: string;
      count: number;
      ids: number[];
    }[] = [];
    const seen = new Map<string, { count: number; ids: number[] }>();

    songs.value.forEach((song) => {
      const key = `${song.title.toLowerCase()}-${song.artist.toLowerCase()}`;
      if (seen.has(key)) {
        const existing = seen.get(key)!;
        existing.count++;
        existing.ids.push(song.id);
      } else {
        seen.set(key, { count: 1, ids: [song.id] });
      }
    });

    seen.forEach((value, key) => {
      if (value.count > 1) {
        const [title, artist] = key.split("-");
        duplicates.push({
          title,
          artist,
          count: value.count,
          ids: value.ids,
        });
      }
    });

    return duplicates;
  };

  return {
    // State
    songs,
    playlists,
    loading,
    error,
    filters,

    // Computed
    filteredSongs,
    genres,
    searchQuery,
    selectedPlaylistId,

    // Actions
    fetchSongs,
    fetchPlaylists,
    createPlaylist,
    updatePlaylist,
    deletePlaylist,
    addSongToPlaylist,
    removeSongFromPlaylist,
    updateSong,
    deleteSong,
    uploadSong,
    scanForDuplicates,
  };
});
