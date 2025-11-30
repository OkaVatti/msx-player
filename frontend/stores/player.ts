import { defineStore } from "pinia";
import type { Song, PlayerState } from "../types";
import { $fetch } from "ofetch";

export const usePlayerStore = defineStore("player", {
  state: (): PlayerState => ({
    isPlaying: false,
    currentSong: null,
    currentTime: 0,
    duration: 0,
    volume: 0.7,
    speed: 1.0,
  }),

  actions: {
    async playSong(song: Song) {
      const { data } = await $fetch("/api/player/play", {
        method: "POST",
        body: { song_id: song.id },
      });
      this.$patch(data);
    },

    async pause() {
      const { data } = await $fetch("/api/player/pause", { method: "POST" });
      this.$patch(data);
    },

    async stop() {
      const { data } = await $fetch("/api/player/stop", { method: "POST" });
      this.$patch(data);
    },

    async seek(time: number) {
      const { data } = await $fetch("/api/player/seek", {
        method: "POST",
        body: { time },
      });
      this.$patch(data);
    },

    async setVolume(volume: number) {
      const { data } = await $fetch("/api/player/volume", {
        method: "POST",
        body: { volume },
      });
      this.$patch(data);
    },

    async setSpeed(speed: number) {
      const { data } = await $fetch("/api/player/speed", {
        method: "POST",
        body: { speed },
      });
      this.$patch(data);
    },

    async rateSong(songId: number, rating: number) {
      await $fetch("/api/player/rate", {
        method: "POST",
        body: { song_id: songId, rating },
      });
    },
  },
});
