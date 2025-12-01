import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { Song } from "../types";

export const usePlayerStore = defineStore("player", () => {
  const currentSong = ref<Song | null>(null);
  const isPlaying = ref(false);
  const currentTime = ref(0);
  const duration = ref(0);
  const volume = ref(0.7);
  const speed = ref(1.0);
  const queue = ref<Song[]>([]);
  const currentIndex = ref(-1);

  // Add audio element
  const audioElement = ref<HTMLAudioElement | null>(null);
  const audioContext = ref<AudioContext | null>(null);
  const audioSource = ref<MediaElementAudioSourceNode | null>(null);

  const progressPercentage = computed(() => {
    if (duration.value === 0) return 0;
    return (currentTime.value / duration.value) * 100;
  });

  // Initialize audio element
  const initAudio = () => {
    if (!audioElement.value) {
      audioElement.value = new Audio();
      audioElement.value.preload = "metadata";

      // Set up event listeners
      audioElement.value.addEventListener("loadedmetadata", () => {
        duration.value = audioElement.value?.duration || 0;
      });

      audioElement.value.addEventListener("timeupdate", () => {
        currentTime.value = audioElement.value?.currentTime || 0;
      });

      audioElement.value.addEventListener("ended", () => {
        isPlaying.value = false;
        nextSong();
      });

      audioElement.value.addEventListener("canplaythrough", () => {
        console.log("🎵 Audio ready to play");
      });

      audioElement.value.addEventListener("error", (e) => {
        console.error("❌ Audio error:", e);
        isPlaying.value = false;
      });
    }
  };

  const playSong = async (song: Song) => {
    try {
      initAudio();

      if (!audioElement.value) {
        console.error("❌ Audio element not initialized");
        return;
      }

      // Stop current playback
      if (isPlaying.value) {
        audioElement.value.pause();
      }

      // Construct the audio URL - the backend serves files from /music/[filename]
      // Extract filename from filePath (remove 'music/' prefix if present)
      let filename = song.filePath;
      if (filename.startsWith("music/")) {
        filename = filename.substring(6); // Remove 'music/' prefix
      }

      const audioUrl = `http://localhost:1323/music/${encodeURIComponent(
        filename
      )}`;
      console.log("🎵 Loading audio from:", audioUrl);

      // Set audio source
      audioElement.value.src = audioUrl;
      audioElement.value.volume = volume.value;
      audioElement.value.playbackRate = speed.value;

      // Wait for metadata to load
      await new Promise((resolve, reject) => {
        if (!audioElement.value) return reject("No audio element");

        const onLoaded = () => {
          audioElement.value?.removeEventListener("loadedmetadata", onLoaded);
          resolve(true);
        };

        const onError = () => {
          audioElement.value?.removeEventListener("error", onError);
          reject("Failed to load audio");
        };

        audioElement.value.addEventListener("loadedmetadata", onLoaded);
        audioElement.value.addEventListener("error", onError);

        // Timeout fallback
        setTimeout(() => {
          audioElement.value?.removeEventListener("loadedmetadata", onLoaded);
          audioElement.value?.removeEventListener("error", onError);
          resolve(true); // Continue even if metadata takes time
        }, 2000);
      });

      // Play the audio
      await audioElement.value.play();

      // Update store state
      currentSong.value = song;
      isPlaying.value = true;
      duration.value = audioElement.value.duration || song.duration || 0;

      console.log("🎵 Now playing:", song.title, "Duration:", duration.value);

      // Update play count on backend
      try {
        const response = await fetch("http://localhost:1323/api/player/play", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ song_id: song.id }),
        });

        if (!response.ok) {
          console.warn("⚠️ Could not update play count on backend");
        }
      } catch (error) {
        console.warn("⚠️ Could not update play count:", error);
      }
    } catch (error) {
      console.error("❌ Failed to play song:", error);
      isPlaying.value = false;

      // Fallback: Try direct file access
      console.log("🔄 Attempting fallback playback method...");
      fallbackPlay(song);
    }
  };

  // Fallback method for audio playback
  const fallbackPlay = (song: Song) => {
    try {
      let filename = song.filePath;
      if (filename.startsWith("music/")) {
        filename = filename.substring(6);
      }

      const audioUrl = `http://localhost:1323/music/${encodeURIComponent(
        filename
      )}`;

      // Create a new audio element for fallback
      const fallbackAudio = new Audio(audioUrl);
      fallbackAudio.volume = volume.value;
      fallbackAudio.playbackRate = speed.value;

      fallbackAudio
        .play()
        .then(() => {
          console.log("🎵 Fallback playback successful");
          currentSong.value = song;
          isPlaying.value = true;

          // Update duration when metadata loads
          fallbackAudio.addEventListener("loadedmetadata", () => {
            duration.value = fallbackAudio.duration;
          });

          // Update current time
          fallbackAudio.addEventListener("timeupdate", () => {
            currentTime.value = fallbackAudio.currentTime;
          });

          fallbackAudio.addEventListener("ended", () => {
            isPlaying.value = false;
            nextSong();
          });

          // Replace the main audio element with fallback
          audioElement.value = fallbackAudio;
        })
        .catch((fallbackError) => {
          console.error("❌ Fallback playback failed:", fallbackError);
          alert(
            `Could not play audio: ${fallbackError}\n\nPlease check:\n1. Backend is running\n2. File exists: ${filename}\n3. CORS is configured`
          );
        });
    } catch (error) {
      console.error("❌ Fallback method also failed:", error);
    }
  };

  const pause = async () => {
    try {
      if (audioElement.value && isPlaying.value) {
        audioElement.value.pause();
        isPlaying.value = false;
      }

      // Optional: Update backend state
      await fetch("http://localhost:1323/api/player/pause", { method: "POST" });
    } catch (error) {
      console.error("Failed to pause:", error);
    }
  };

  const stop = async () => {
    try {
      if (audioElement.value) {
        audioElement.value.pause();
        audioElement.value.currentTime = 0;
        isPlaying.value = false;
        currentTime.value = 0;
      }

      await fetch("http://localhost:1323/api/player/stop", { method: "POST" });
    } catch (error) {
      console.error("Failed to stop:", error);
    }
  };

  const seek = async (time: number) => {
    try {
      if (audioElement.value) {
        audioElement.value.currentTime = time;
        currentTime.value = time;
      }

      await fetch("http://localhost:1323/api/player/seek", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ time }),
      });
    } catch (error) {
      console.error("Failed to seek:", error);
    }
  };

  const setVolume = async (newVolume: number) => {
    try {
      volume.value = newVolume;
      if (audioElement.value) {
        audioElement.value.volume = newVolume;
      }

      await fetch("http://localhost:1323/api/player/volume", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ volume: newVolume }),
      });
    } catch (error) {
      console.error("Failed to set volume:", error);
    }
  };

  const setSpeed = async (newSpeed: number) => {
    try {
      speed.value = newSpeed;
      if (audioElement.value) {
        audioElement.value.playbackRate = newSpeed;
      }

      await fetch("http://localhost:1323/api/player/speed", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ speed: newSpeed }),
      });
    } catch (error) {
      console.error("Failed to set speed:", error);
    }
  };

  const rateSong = async (songId: number, rating: number) => {
    try {
      const response = await fetch("http://localhost:1323/api/player/rate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ song_id: songId, rating }),
      });

      if (response.ok) {
        console.log("⭐ Rated song:", songId, "as", rating);
        if (currentSong.value && currentSong.value.id === songId) {
          currentSong.value.rating = rating;
        }
      }
    } catch (error) {
      console.error("Failed to rate song:", error);
    }
  };

  const nextSong = async () => {
    try {
      if (queue.value.length > 0) {
        currentIndex.value = (currentIndex.value + 1) % queue.value.length;
        const nextSong = queue.value[currentIndex.value];
        if (nextSong) {
          await playSong(nextSong);
        }
      } else {
        // If no queue, try to play next song from library
        console.log(
          "No queue set up, skipping to next song not implemented yet"
        );
      }

      await fetch("http://localhost:1323/api/player/next", { method: "POST" });
    } catch (error) {
      console.error("Failed to skip to next song:", error);
    }
  };

  const previousSong = async () => {
    try {
      if (queue.value.length > 0) {
        currentIndex.value =
          currentIndex.value > 0
            ? currentIndex.value - 1
            : queue.value.length - 1;
        const prevSong = queue.value[currentIndex.value];
        if (prevSong) {
          await playSong(prevSong);
        }
      } else {
        console.log("No queue set up, previous song not implemented yet");
      }

      await fetch("http://localhost:1323/api/player/previous", {
        method: "POST",
      });
    } catch (error) {
      console.error("Failed to go to previous song:", error);
    }
  };

  // Set queue for next/previous functionality
  const setQueue = (songs: Song[], startIndex: number = 0) => {
    queue.value = songs;
    currentIndex.value = startIndex;
  };

  // Clean up audio element
  const cleanup = () => {
    if (audioElement.value) {
      audioElement.value.pause();
      audioElement.value.src = "";
      audioElement.value = null;
    }
  };

  return {
    currentSong,
    isPlaying,
    currentTime,
    duration,
    volume,
    speed,
    queue,
    currentIndex,
    progressPercentage,
    audioElement,
    playSong,
    pause,
    stop,
    seek,
    setVolume,
    setSpeed,
    rateSong,
    nextSong,
    previousSong,
    setQueue,
    cleanup,
    initAudio,
  };
});
