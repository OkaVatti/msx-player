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
            new Date(a.last_played || 0).getTime()
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
  const fetchSongs = async () => {
    loading.value = true;
    error.value = null;

    try {
      const url = new URL("http://localhost:1323/api/songs");

      // Add filters to query params
      if (filters.value.search) {
        url.searchParams.append("search", filters.value.search);
      }
      if (filters.value.sortBy) {
        url.searchParams.append("sort_by", filters.value.sortBy);
      }
      if (filters.value.genre) {
        url.searchParams.append("genre", filters.value.genre);
      }
      if (filters.value.rating !== null) {
        url.searchParams.append("rating", filters.value.rating.toString());
      }

      const response = await fetch(url);

      if (!response.ok) {
        throw new Error(`Failed to fetch songs: ${response.status}`);
      }

      songs.value = await response.json();
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Unknown error";
      console.error("Error fetching songs:", err);
    } finally {
      loading.value = false;
    }
  };

  const fetchPlaylists = async () => {
    try {
      const response = await fetch("http://localhost:1323/api/playlists");
      if (!response.ok) {
        throw new Error(`Failed to fetch playlists: ${response.status}`);
      }
      playlists.value = await response.json();
    } catch (err) {
      console.error("Error fetching playlists:", err);
    }
  };

  const createPlaylist = async (name: string, description: string = "") => {
    try {
      const response = await fetch("http://localhost:1323/api/playlists", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, description }),
      });

      if (!response.ok) {
        throw new Error("Failed to create playlist");
      }

      const newPlaylist = await response.json();
      playlists.value.push(newPlaylist);
      return newPlaylist;
    } catch (err) {
      console.error("Error creating playlist:", err);
      throw err;
    }
  };

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

  const deletePlaylist = async (id: number) => {
    try {
      const response = await fetch(
        `http://localhost:1323/api/playlists/${id}`,
        {
          method: "DELETE",
        }
      );

      if (!response.ok) throw new Error("Failed to delete playlist");

      playlists.value = playlists.value.filter((p) => p.id !== id);
    } catch (err) {
      console.error("Error deleting playlist:", err);
      throw err;
    }
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
