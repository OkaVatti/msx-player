// stores/library.ts
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { Song, Playlist, LibraryFilters } from "../types";

export const useLibraryStore = defineStore("library", () => {
  const songs = ref<Song[]>([]);
  const playlists = ref<Playlist[]>([]);
  const selectedPlaylistId = ref<number | null>(null);

  const filters = ref<LibraryFilters>({
    search: "",
    sortBy: "title",
    genre: "all",
    rating: null,
  });

  const filteredSongs = computed(() => {
    let res = [...songs.value];
    const q = filters.value.search?.trim().toLowerCase() ?? "";
    if (q) {
      res = res.filter(
        (s) =>
          (s.title || "").toLowerCase().includes(q) ||
          (s.artist || "").toLowerCase().includes(q) ||
          (s.album || "").toLowerCase().includes(q)
      );
    }

    switch (filters.value.sortBy) {
      case "title":
        res.sort((a, b) => a.title.localeCompare(b.title));
        break;
      case "artist":
        res.sort((a, b) => a.artist.localeCompare(b.artist));
        break;
      case "album":
        res.sort((a, b) => (a.album || "").localeCompare(b.album || ""));
        break;
      case "rating":
        res.sort((a, b) => (b.rating || 0) - (a.rating || 0));
        break;
      case "play_count":
        res.sort((a, b) => (b.play_count || 0) - (a.play_count || 0));
        break;
      case "last_played":
        res.sort((a, b) => {
          const A = a.last_played ? new Date(a.last_played).getTime() : 0;
          const B = b.last_played ? new Date(b.last_played).getTime() : 0;
          return B - A;
        });
        break;
      case "duration":
        res.sort((a, b) => a.duration - b.duration);
        break;
      case "year":
        res.sort((a, b) => (a.year || 0) - (b.year || 0));
        break;
      default:
        res.sort((a, b) => a.title.localeCompare(b.title));
    }

    return res;
  });

  async function fetchSongs() {
    try {
      const r = await fetch("/api/songs");
      if (!r.ok) throw new Error("fetch songs failed");
      songs.value = await r.json();
    } catch (err) {
      console.error("fetchSongs", err);
    }
  }

  async function fetchPlaylists() {
    try {
      const r = await fetch("/api/playlists");
      if (!r.ok) throw new Error("fetch playlists failed");
      playlists.value = await r.json();
    } catch (err) {
      console.error("fetchPlaylists", err);
    }
  }

  async function uploadSong(file: File) {
    const fd = new FormData();
    fd.append("file", file);
    const r = await fetch("/api/songs/upload", { method: "POST", body: fd });
    if (!r.ok) throw new Error("upload failed");
    const s: Song = await r.json();
    songs.value.unshift(s);
    return s;
  }

  async function createPlaylist(name: string, description?: string) {
    const r = await fetch("/api/playlists", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, description }),
    });
    if (!r.ok) throw new Error("create playlist failed");
    const pl: Playlist = await r.json();
    playlists.value.unshift(pl);
    return pl;
  }

  return {
    songs,
    playlists,
    selectedPlaylistId,
    filters,
    filteredSongs,
    fetchSongs,
    fetchPlaylists,
    uploadSong,
    createPlaylist,
  };
});
