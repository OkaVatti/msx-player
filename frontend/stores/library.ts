import { defineStore } from "pinia";
import type { Song, Playlist } from "../types";
import { $fetch } from "ofetch";

export const useLibraryStore = defineStore("library", {
  state: () => ({
    songs: [] as Song[],
    playlists: [] as Playlist[],
    searchQuery: "",
    filters: {
      genre: "",
      rating: 0,
      sortBy: "title",
    },
  }),

  getters: {
    filteredSongs: (state) => {
      let songs = state.songs;

      // Search filter
      if (state.searchQuery) {
        const query = state.searchQuery.toLowerCase();
        songs = songs.filter(
          (song) =>
            song.title.toLowerCase().includes(query) ||
            song.artist.toLowerCase().includes(query) ||
            song.album.toLowerCase().includes(query)
        );
      }

      // Genre filter
      if (state.filters.genre) {
        songs = songs.filter((song) => song.genre === state.filters.genre);
      }

      // Rating filter
      if (state.filters.rating > 0) {
        songs = songs.filter((song) => song.rating >= state.filters.rating);
      }

      // Sort
      songs.sort((a, b) => {
        const aVal = a[state.filters.sortBy as keyof Song];
        const bVal = b[state.filters.sortBy as keyof Song];

        if (typeof aVal === "string" && typeof bVal === "string") {
          return aVal.localeCompare(bVal);
        }
        return (aVal as number) - (bVal as number);
      });

      return songs;
    },

    genres: (state) => {
      return [
        ...new Set(state.songs.map((song) => song.genre).filter(Boolean)),
      ];
    },
  },

  actions: {
    async fetchSongs() {
      const query = new URLSearchParams();
      if (this.searchQuery) query.append("search", this.searchQuery);
      if (this.filters.genre) query.append("genre", this.filters.genre);
      if (this.filters.rating)
        query.append("rating", this.filters.rating.toString());
      if (this.filters.sortBy) query.append("sort_by", this.filters.sortBy);

      const { data } = await $fetch(`/api/songs?${query}`);
      this.songs = data;
    },

    async uploadSong(file: File) {
      const formData = new FormData();
      formData.append("file", file);

      await $fetch("/api/songs", {
        method: "POST",
        body: formData,
      });
      await this.fetchSongs();
    },

    async deleteSong(songId: number) {
      await $fetch(`/api/songs/${songId}`, { method: "DELETE" });
      await this.fetchSongs();
    },

    async updateSong(songId: number, updates: Partial<Song>) {
      await $fetch(`/api/songs/${songId}`, {
        method: "PUT",
        body: updates,
      });
      await this.fetchSongs();
    },

    async fetchPlaylists() {
      const { data } = await $fetch("/api/playlists");
      this.playlists = data;
    },

    async createPlaylist(playlist: Omit<Playlist, "id" | "createdAt">) {
      await $fetch("/api/playlists", {
        method: "POST",
        body: playlist,
      });
      await this.fetchPlaylists();
    },
  },
});
