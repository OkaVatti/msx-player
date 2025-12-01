// stores/player.ts - Extended version (patched)
import { defineStore } from "pinia";
import { ref, computed, onMounted, onUnmounted } from "vue";
import type { Song, PlayerState } from "../types";

export const usePlayerStore = defineStore("player", () => {
  // State
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
  const ws = ref<WebSocket | null>(null);
  const isConnected = ref(false);

  // Computed
  const progressPercentage = computed(() => {
    if (!duration.value) return 0;
    return (currentTime.value / duration.value) * 100;
  });

  // WebSocket Connection
  const connectWebSocket = () => {
    try {
      const socket = new WebSocket("ws://localhost:1323/ws");

      socket.onopen = () => {
        console.log("WebSocket connected");
        isConnected.value = true;
      };

      socket.onmessage = (event) => {
        try {
          const state: PlayerState = JSON.parse(event.data);
          updateFromServerState(state);
        } catch (err) {
          console.error("Error parsing WebSocket message:", err);
        }
      };

      socket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      socket.onclose = () => {
        console.log("WebSocket disconnected");
        isConnected.value = false;
        // Attempt to reconnect after 3 seconds
        setTimeout(connectWebSocket, 3000);
      };

      ws.value = socket;
    } catch (err) {
      console.error("Failed to create WebSocket:", err);
    }
  };

  // Update local state from server state
  const updateFromServerState = (state: PlayerState) => {
    isPlaying.value = state.is_playing;
    currentSong.value = state.current_song;
    currentTime.value = state.current_time;
    duration.value = state.duration;
    volume.value = state.volume;
    speed.value = state.speed;

    // Update audio element if it exists
    if (audioElement.value) {
      audioElement.value.volume = state.volume;
      audioElement.value.playbackRate = state.speed;

      if (
        state.current_song &&
        audioElement.value.src !==
          `http://localhost:1323${state.current_song.file_path}`
      ) {
        audioElement.value.src = `http://localhost:1323${state.current_song.file_path}`;
        if (state.is_playing) {
          audioElement.value.play().catch(console.error);
        }
      }
    }
  };

  // Actions
  const setAudioElement = (element: HTMLAudioElement) => {
    audioElement.value = element;

    // Set up event listeners
    element.addEventListener("timeupdate", () => {
      currentTime.value = element.currentTime;
    });

    element.addEventListener("loadedmetadata", () => {
      duration.value = element.duration;
    });

    element.addEventListener("ended", handleSongEnded);
  };

  const fetchPlayerState = async () => {
    try {
      const response = await fetch("http://localhost:1323/api/player/state");
      if (!response.ok) throw new Error("Failed to fetch player state");

      const state: PlayerState = await response.json();
      updateFromServerState(state);
    } catch (err) {
      console.error("Error fetching player state:", err);
    }
  };

  const playSong = async (song: Song) => {
    try {
      const response = await fetch("http://localhost:1323/api/player/play", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ song_id: song.id }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to play song");
      }

      // Server will update state via WebSocket
    } catch (error) {
      console.error("Error playing song:", error);
      // Fallback to local playback if WebSocket not available
      if (!isConnected.value) {
        await playSongLocally(song);
      }
    }
  };

  const playSongLocally = async (song: Song) => {
    // Create audio element if it doesn't exist
    if (!audioElement.value) {
      audioElement.value = new Audio();
      setAudioElement(audioElement.value);
    }

    // Set source and metadata
    audioElement.value.src = `http://localhost:1323${song.file_path}`;

    // Update play count via API (best-effort)
    try {
      await fetch(`http://localhost:1323/api/player/rate`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ song_id: song.id, rating: song.rating }),
      });
    } catch (err) {
      // non-fatal
      console.warn("Failed to update play count/rating:", err);
    }

    // Play the song
    audioElement.value.volume = volume.value;
    audioElement.value.playbackRate = speed.value;
    await audioElement.value.play();

    currentSong.value = song;
    isPlaying.value = true;
    currentTime.value = 0;

    // Add to queue if not already
    if (!queue.value.some((s) => s.id === song.id)) {
      queue.value.push(song);
      currentQueueIndex.value = queue.value.length - 1;
    } else {
      currentQueueIndex.value = queue.value.findIndex((s) => s.id === song.id);
    }
  };

  const pause = async () => {
    try {
      const response = await fetch("http://localhost:1323/api/player/pause", {
        method: "POST",
      });

      if (!response.ok) throw new Error("Failed to pause");
    } catch (err) {
      console.error("Error pausing:", err);
      if (audioElement.value) {
        audioElement.value.pause();
        isPlaying.value = false;
      }
    }
  };

  const resume = async () => {
    // Resume either via server or locally as a fallback
    try {
      const response = await fetch("http://localhost:1323/api/player/resume", {
        method: "POST",
      });

      if (!response.ok) throw new Error("Failed to resume");
    } catch (err) {
      console.error("Error resuming via server:", err);
      if (audioElement.value) {
        await audioElement.value.play().catch(console.error);
        isPlaying.value = true;
      }
    }
  };

  const stop = async () => {
    try {
      const response = await fetch("http://localhost:1323/api/player/stop", {
        method: "POST",
      });

      if (!response.ok) throw new Error("Failed to stop");
    } catch (err) {
      console.error("Error stopping:", err);
      if (audioElement.value) {
        audioElement.value.pause();
        audioElement.value.currentTime = 0;
        isPlaying.value = false;
        currentTime.value = 0;
      }
    }
  };

  const nextSong = async () => {
    try {
      const response = await fetch("http://localhost:1323/api/player/next", {
        method: "POST",
      });

      if (!response.ok) throw new Error("Failed to skip to next song");
    } catch (err) {
      console.error("Error skipping to next song:", err);
      if (!isConnected.value) {
        await nextSongLocally();
      }
    }
  };

  const nextSongLocally = async () => {
    if (queue.value.length === 0) return;

    let nextIndex = currentQueueIndex.value + 1;

    if (shuffle.value) {
      nextIndex = Math.floor(Math.random() * queue.value.length);
    } else if (nextIndex >= queue.value.length) {
      if (repeatMode.value === "all") {
        nextIndex = 0;
      } else {
        return; // Stop if no repeat
      }
    }

    currentQueueIndex.value = nextIndex;
    await playSongLocally(queue.value[nextIndex]);
  };

  const previousSong = async () => {
    try {
      const response = await fetch(
        "http://localhost:1323/api/player/previous",
        {
          method: "POST",
        }
      );

      if (!response.ok) throw new Error("Failed to go to previous song");
    } catch (err) {
      console.error("Error going to previous song:", err);
      if (!isConnected.value) {
        await previousSongLocally();
      }
    }
  };

  const previousSongLocally = async () => {
    if (queue.value.length === 0) return;

    let prevIndex = currentQueueIndex.value - 1;

    if (prevIndex < 0) {
      if (repeatMode.value === "all") {
        prevIndex = queue.value.length - 1;
      } else {
        return; // Stop if no repeat
      }
    }

    currentQueueIndex.value = prevIndex;
    await playSongLocally(queue.value[prevIndex]);
  };

  const handleSongEnded = async () => {
    if (repeatMode.value === "one") {
      // Play same song again via server or locally
      if (currentSong.value) {
        await playSong(currentSong.value);
      }
    } else {
      await nextSong();
    }
  };

  const seek = async (time: number) => {
    try {
      const response = await fetch("http://localhost:1323/api/player/seek", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ time }),
      });

      if (!response.ok) throw new Error("Failed to seek");
    } catch (err) {
      console.error("Error seeking:", err);
      if (audioElement.value) {
        audioElement.value.currentTime = time;
        currentTime.value = time;
      }
    }
  };

  const setVolume = async (newVolume: number) => {
    try {
      const response = await fetch("http://localhost:1323/api/player/volume", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ volume: newVolume }),
      });

      if (!response.ok) throw new Error("Failed to set volume");
    } catch (err) {
      console.error("Error setting volume:", err);
      volume.value = newVolume;
      if (audioElement.value) {
        audioElement.value.volume = newVolume;
      }
    }
  };

  const setSpeed = async (newSpeed: number) => {
    try {
      const response = await fetch("http://localhost:1323/api/player/speed", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ speed: newSpeed }),
      });

      if (!response.ok) throw new Error("Failed to set speed");
    } catch (err) {
      console.error("Error setting speed:", err);
      speed.value = newSpeed;
      if (audioElement.value) {
        audioElement.value.playbackRate = newSpeed;
      }
    }
  };

  const setRepeatMode = (mode: "none" | "all" | "one") => {
    repeatMode.value = mode;
  };

  const setShuffle = (value: boolean) => {
    shuffle.value = value;
  };

  const setQueue = (songs: Song[], startIndex: number = 0) => {
    queue.value = songs;
    currentQueueIndex.value = startIndex;
  };

  const shuffleQueue = async () => {
    try {
      const response = await fetch("http://localhost:1323/api/player/shuffle", {
        method: "POST",
      });

      if (!response.ok) throw new Error("Failed to shuffle queue");
    } catch (err) {
      console.error("Error shuffling queue:", err);
      const shuffled = [...queue.value];
      for (let i = shuffled.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
      }
      queue.value = shuffled;
    }
  };

  const rateSong = async (songId: number, rating: number) => {
    try {
      const response = await fetch("http://localhost:1323/api/player/rate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ song_id: songId, rating }),
      });

      if (!response.ok) throw new Error("Failed to rate song");

      // Update local state
      if (currentSong.value?.id === songId) {
        currentSong.value.rating = rating;
      }

      // Update in queue
      const songInQueue = queue.value.find((s) => s.id === songId);
      if (songInQueue) {
        songInQueue.rating = rating;
      }
    } catch (error) {
      console.error("Error rating song:", error);
    }
  };

  const hasNext = computed(() => {
    if (!queue.value || queue.value.length === 0) return false;
    if (shuffle.value) return queue.value.length > 1;
    // if repeat all, always has next (wraps)
    if (repeatMode.value === "all") return queue.value.length > 0;
    return currentQueueIndex.value < queue.value.length - 1;
  });

  const hasPrevious = computed(() => {
    if (!queue.value || queue.value.length === 0) return false;
    // if shuffle, previous is available when queue length > 1
    if (shuffle.value) return queue.value.length > 1;
    if (repeatMode.value === "all") return queue.value.length > 0;
    return currentQueueIndex.value > 0;
  });

  // Lifecycle
  onMounted(() => {
    connectWebSocket();
    fetchPlayerState();
  });

  onUnmounted(() => {
    if (ws.value) {
      ws.value.close();
    }
  });

  return {
    // State
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
    isConnected,

    // Computed
    progressPercentage,
    hasNext,
    hasPrevious,

    // Actions
    setAudioElement,
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
    shuffleQueue,
    rateSong,
    // WebSocket / lifecycle helpers if you want them accessible
    connectWebSocket,
    fetchPlayerState,
  };
});
