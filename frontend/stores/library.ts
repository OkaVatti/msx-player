import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { Song } from "../types";

export const useLibraryStore = defineStore("library", () => {
  const songs = ref<Song[]>([]);
  const searchQuery = ref("");
  const filters = ref({
    genre: "",
    sortBy: "title",
  });
  const loading = ref(false);

  const filteredSongs = computed(() => {
    let filtered = [...songs.value];

    // Apply search
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      filtered = filtered.filter(
        (song) =>
          song.title?.toLowerCase().includes(query) ||
          song.artist?.toLowerCase().includes(query) ||
          song.album?.toLowerCase().includes(query)
      );
    }

    // Apply genre filter
    if (filters.value.genre) {
      filtered = filtered.filter((song) => song.genre === filters.value.genre);
    }

    // Apply sorting
    filtered.sort((a, b) => {
      const aVal = a[filters.value.sortBy as keyof Song];
      const bVal = b[filters.value.sortBy as keyof Song];

      if (typeof aVal === "string" && typeof bVal === "string") {
        return aVal.localeCompare(bVal);
      }
      if (typeof aVal === "number" && typeof bVal === "number") {
        return bVal - aVal; // Descending for numbers
      }
      return 0;
    });

    return filtered;
  });

  const genres = computed(() => {
    const genreSet = new Set(
      songs.value.map((song) => song.genre).filter(Boolean)
    );
    return Array.from(genreSet).sort();
  });

  const fetchSongs = async () => {
    try {
      loading.value = true;
      const params = new URLSearchParams();
      if (searchQuery.value) params.append("search", searchQuery.value);
      if (filters.value.genre) params.append("genre", filters.value.genre);
      params.append("sort_by", filters.value.sortBy);

      const response = await fetch(`http://localhost:1323/api/songs?${params}`);
      if (response.ok) {
        const songsData = await response.json();
        songs.value = songsData.map((song: any) => ({
          ...song,
          // Ensure consistent field names
          playCount: song.play_count || song.playCount || 0,
          filePath: song.file_path || song.filePath,
          lastPlayed: song.last_played || song.lastPlayed,
          createdAt: song.created_at || song.createdAt,
        }));
        console.log("📚 Loaded songs:", songs.value.length);
      } else {
        console.error("Failed to fetch songs");
      }
    } catch (error) {
      console.error("Error fetching songs:", error);
    } finally {
      loading.value = false;
    }
  };

  const deleteSong = async (id: number) => {
    try {
      const response = await fetch(`http://localhost:1323/api/songs/${id}`, {
        method: "DELETE",
      });

      if (response.ok) {
        songs.value = songs.value.filter((song) => song.id !== id);
        console.log("🗑️ Deleted song:", id);
        return true;
      } else {
        console.error("Failed to delete song");
        return false;
      }
    } catch (error) {
      console.error("Error deleting song:", error);
      return false;
    }
  };

  const updateSong = async (id: number, updates: Partial<Song>) => {
    try {
      const response = await fetch(`http://localhost:1323/api/songs/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(updates),
      });

      if (response.ok) {
        const updatedSong = await response.json();
        const index = songs.value.findIndex((song) => song.id === id);
        if (index !== -1) {
          songs.value[index] = {
            ...updatedSong,
            playCount: updatedSong.play_count || updatedSong.playCount || 0,
            filePath: updatedSong.file_path || updatedSong.filePath,
          };
        }
        console.log("✏️ Updated song:", id);
        return true;
      } else {
        console.error("Failed to update song");
        return false;
      }
    } catch (error) {
      console.error("Error updating song:", error);
      return false;
    }
  };

  return {
    songs,
    searchQuery,
    filters,
    loading,
    filteredSongs,
    genres,
    fetchSongs,
    deleteSong,
    updateSong,
  };
});
