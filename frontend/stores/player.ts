// stores/player.ts
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { Song, PlayerState } from "../types";

export const usePlayerStore = defineStore("player", () => {
  const currentSong = ref<Song | null>(null);
  const isPlaying = ref(false);
  const currentTime = ref(0);
  const duration = ref(0);
  const volume = ref(0.7);
  const speed = ref(1.0);
  const repeatMode = ref<"none" | "all" | "one">("none");
  const shuffle = ref(false);
  const queue = ref<Song[]>([]);
  const currentQueueIndex = ref(-1);
  const audioElement = ref<HTMLAudioElement | null>(null);
  const isVisualizerConnected = ref(false);

  const progressPercentage = computed(() => {
    if (!duration.value || duration.value === 0) return 0;
    return Math.max(
      0,
      Math.min(100, (currentTime.value / duration.value) * 100)
    );
  });

  const hasNext = computed(() => {
    if (!queue.value || queue.value.length === 0) return false;
    if (shuffle.value) return queue.value.length > 1;
    if (repeatMode.value === "all") return queue.value.length > 0;
    return currentQueueIndex.value < queue.value.length - 1;
  });

  const hasPrevious = computed(() => {
    if (!queue.value || queue.value.length === 0) return false;
    if (shuffle.value) return queue.value.length > 1;
    if (repeatMode.value === "all") return queue.value.length > 0;
    return currentQueueIndex.value > 0;
  });

  function setAudioElement(el: HTMLAudioElement) {
    audioElement.value = el;
    el.addEventListener("timeupdate", () => {
      currentTime.value = el.currentTime;
    });
    el.addEventListener("loadedmetadata", () => {
      duration.value = el.duration;
    });
    el.addEventListener("ended", () => {
      handleSongEnded();
    });
  }

  async function fetchPlayerState() {
    try {
      const r = await fetch("/api/player/state");
      if (!r.ok) return;
      const st: PlayerState = await r.json();
      isPlaying.value = st.is_playing;
      currentSong.value = st.current_song ?? null;
      currentTime.value = st.current_time;
      duration.value = st.duration;
      volume.value = st.volume;
      speed.value = st.speed;
      if (audioElement.value && currentSong.value) {
        const expected = `${location.origin}${currentSong.value.file_path}`;
        if (audioElement.value.src !== expected) {
          audioElement.value.src = expected;
        }
        if (st.is_playing) {
          await audioElement.value.play().catch(() => {});
        }
      }
    } catch (err) {
      // swallow; client will play locally
      console.warn("fetchPlayerState error", err);
    }
  }

  async function playSong(song: Song) {
    try {
      const r = await fetch("/api/player/play", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ song_id: song.id }),
      });
      if (!r.ok) throw new Error("server play failed");
      await fetchPlayerState();
    } catch {
      // fallback to local
      await playSongLocally(song);
    }
  }

  async function playSongLocally(song: Song) {
    if (!audioElement.value) {
      audioElement.value = new Audio();
      setAudioElement(audioElement.value);
    }
    audioElement.value.src = song.file_path.startsWith("http")
      ? song.file_path
      : `${location.origin}${song.file_path}`;
    audioElement.value.volume = volume.value;
    audioElement.value.playbackRate = speed.value;
    await audioElement.value.play().catch(console.error);
    currentSong.value = song;
    isPlaying.value = true;
    currentTime.value = 0;

    const idx = queue.value.findIndex((s) => s.id === song.id);
    if (idx === -1) {
      queue.value.push(song);
      currentQueueIndex.value = queue.value.length - 1;
    } else {
      currentQueueIndex.value = idx;
    }
  }

  async function pause() {
    try {
      await fetch("/api/player/pause", { method: "POST" });
      isPlaying.value = false;
      if (audioElement.value) audioElement.value.pause();
    } catch {
      if (audioElement.value) {
        audioElement.value.pause();
        isPlaying.value = false;
      }
    }
  }

  async function resume() {
    try {
      await fetch("/api/player/resume", { method: "POST" });
    } catch {
      if (audioElement.value) {
        await audioElement.value.play().catch(console.error);
        isPlaying.value = true;
      }
    }
  }

  async function stop() {
    try {
      await fetch("/api/player/stop", { method: "POST" });
    } catch {
      if (audioElement.value) {
        audioElement.value.pause();
        audioElement.value.currentTime = 0;
        isPlaying.value = false;
      }
    }
  }

  async function nextSong() {
    try {
      await fetch("/api/player/next", { method: "POST" });
    } catch {
      await nextSongLocally();
    }
  }

  async function nextSongLocally() {
    if (!queue.value.length) return;
    let nextIdx = currentQueueIndex.value + 1;
    if (shuffle.value) nextIdx = Math.floor(Math.random() * queue.value.length);
    else if (nextIdx >= queue.value.length) {
      if (repeatMode.value === "all") nextIdx = 0;
      else return;
    }
    currentQueueIndex.value = nextIdx;
    await playSongLocally(queue.value[nextIdx]);
  }

  async function previousSong() {
    try {
      await fetch("/api/player/previous", { method: "POST" });
    } catch {
      await previousSongLocally();
    }
  }

  async function previousSongLocally() {
    if (!queue.value.length) return;
    let idx = currentQueueIndex.value - 1;
    if (idx < 0) {
      if (repeatMode.value === "all") idx = queue.value.length - 1;
      else return;
    }
    currentQueueIndex.value = idx;
    await playSongLocally(queue.value[idx]);
  }

  async function handleSongEnded() {
    if (repeatMode.value === "one") {
      if (currentSong.value) await playSong(currentSong.value);
    } else {
      await nextSong();
    }
  }

  async function seek(t: number) {
    try {
      await fetch("/api/player/seek", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ time: t }),
      });
    } catch {
      if (audioElement.value) {
        audioElement.value.currentTime = t;
        currentTime.value = t;
      }
    }
  }

  async function setVolume(v: number) {
    volume.value = v;
    if (audioElement.value) audioElement.value.volume = v;
    try {
      await fetch("/api/player/volume", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ volume: v }),
      });
    } catch {}
  }

  async function setSpeed(v: number) {
    speed.value = v;
    if (audioElement.value) audioElement.value.playbackRate = v;
    try {
      await fetch("/api/player/speed", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ speed: v }),
      });
    } catch {}
  }

  async function rateSong(songId: number, rating: number) {
    try {
      await fetch("/api/player/rate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ song_id: songId, rating }),
      });
      if (currentSong.value && currentSong.value.id === songId)
        currentSong.value.rating = rating;
      const inQueue = queue.value.find((s) => s.id === songId);
      if (inQueue) inQueue.rating = rating;
    } catch (err) {
      console.warn("rate failed", err);
    }
  }

  function setRepeatMode(m: "none" | "all" | "one") {
    repeatMode.value = m;
  }
  function setShuffle(s: boolean) {
    shuffle.value = s;
  }
  function setQueue(songs: Song[], startIndex = 0) {
    queue.value = songs;
    currentQueueIndex.value = startIndex;
  }

  return {
    // state
    currentSong,
    isPlaying,
    currentTime,
    duration,
    volume,
    speed,
    repeatMode,
    shuffle,
    queue,
    currentQueueIndex,
    audioElement,
    isVisualizerConnected,
    // computed
    progressPercentage,
    hasNext,
    hasPrevious,
    // actions
    setAudioElement,
    fetchPlayerState,
    playSong,
    playSongLocally,
    pause,
    resume,
    stop,
    nextSong,
    previousSong,
    nextSongLocally,
    previousSongLocally,
    handleSongEnded,
    seek,
    setVolume,
    setSpeed,
    setRepeatMode,
    setShuffle,
    setQueue,
    rateSong,
  };
});
