<!-- App.vue (main file) -->
<template>
  <div class="min-h-screen bg-black text-white font-mono overflow-hidden">
    <!-- Top Navigation Bar -->
    <div class="h-12 bg-black border-b border-[#837dbd] flex items-center justify-between px-4">
      <div class="flex items-center gap-6 text-sm">
        <div class="text-[#d3ceff] font-bold tracking-wide">MSX PLAYER</div>
        <div class="flex gap-4">
          <button 
            @click="showFileMenu = !showFileMenu"
            class="text-[#837dbd] hover:text-[#d3ceff] transition-colors relative"
          >
            FILE
            <!-- File Menu Dropdown -->
            <div 
              v-if="showFileMenu"
              class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48"
            >
              <button 
                @click="uploadSong"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                UPLOAD SONG...
              </button>
              <button 
                @click="scanLibrary"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                SCAN LIBRARY
              </button>
              <div class="border-t border-[#5a548d] my-1"></div>
              <button 
                @click="exportLibrary"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                EXPORT LIBRARY
              </button>
              <button 
                @click="importLibrary"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                IMPORT LIBRARY
              </button>
            </div>
          </button>
          <button 
            @click="showEditMenu = !showEditMenu"
            class="text-[#837dbd] hover:text-[#d3ceff] transition-colors relative"
          >
            EDIT
            <!-- Edit Menu Dropdown -->
            <div 
              v-if="showEditMenu"
              class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48"
            >
              <button 
                @click="findDuplicates"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                FIND DUPLICATES
              </button>
              <button 
                @click="cleanLibrary"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                CLEAN LIBRARY
              </button>
              <button 
                @click="updateMetadata"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                UPDATE METADATA
              </button>
            </div>
          </button>
          <button 
            @click="showViewMenu = !showViewMenu"
            class="text-[#837dbd] hover:text-[#d3ceff] transition-colors relative"
          >
            VIEW
            <!-- View Menu Dropdown -->
            <div 
              v-if="showViewMenu"
              class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48"
            >
              <button 
                @click="toggleDarkMode"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                {{ darkMode ? 'LIGHT MODE' : 'DARK MODE' }}
              </button>
              <button 
                @click="toggleCompactView"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                {{ compactView ? 'EXPAND VIEW' : 'COMPACT VIEW' }}
              </button>
              <div class="border-t border-[#5a548d] my-1"></div>
              <button 
                @click="changeView('library')"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                LIBRARY VIEW
              </button>
              <button 
                @click="changeView('playlists')"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                PLAYLISTS VIEW
              </button>
              <button 
                @click="changeView('visualizer')"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                VISUALIZER VIEW
              </button>
            </div>
          </button>
          <button 
            @click="showControlsMenu = !showControlsMenu"
            class="text-[#837dbd] hover:text-[#d3ceff] transition-colors relative"
          >
            CONTROLS
            <!-- Controls Menu Dropdown -->
            <div 
              v-if="showControlsMenu"
              class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48"
            >
              <button 
                @click="playerStore.shuffleQueue()"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                SHUFFLE QUEUE
              </button>
              <button 
                @click="clearQueue"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                CLEAR QUEUE
              </button>
              <div class="border-t border-[#5a548d] my-1"></div>
              <div class="px-3 py-1 text-xs text-[#5a548d]">REPEAT MODE</div>
              <div class="flex gap-1 px-3 py-1">
                <button 
                  @click="playerStore.setRepeatMode('none')"
                  :class="[
                    'px-2 py-1 text-xs border',
                    playerStore.repeatMode === 'none' 
                      ? 'border-[#d3ceff] text-[#d3ceff]' 
                      : 'border-[#5a548d] text-[#5a548d] hover:border-[#837dbd] hover:text-[#837dbd]'
                  ]"
                >
                  OFF
                </button>
                <button 
                  @click="playerStore.setRepeatMode('all')"
                  :class="[
                    'px-2 py-1 text-xs border',
                    playerStore.repeatMode === 'all' 
                      ? 'border-[#d3ceff] text-[#d3ceff]' 
                      : 'border-[#5a548d] text-[#5a548d] hover:border-[#837dbd] hover:text-[#837dbd]'
                  ]"
                >
                  ALL
                </button>
                <button 
                  @click="playerStore.setRepeatMode('one')"
                  :class="[
                    'px-2 py-1 text-xs border',
                    playerStore.repeatMode === 'one' 
                      ? 'border-[#d3ceff] text-[#d3ceff]' 
                      : 'border-[#5a548d] text-[#5a548d] hover:border-[#837dbd] hover:text-[#837dbd]'
                  ]"
                >
                  ONE
                </button>
              </div>
            </div>
          </button>
          <button 
            @click="showHelpMenu = !showHelpMenu"
            class="text-[#837dbd] hover:text-[#d3ceff] transition-colors relative"
          >
            HELP
            <!-- Help Menu Dropdown -->
            <div 
              v-if="showHelpMenu"
              class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48"
            >
              <button 
                @click="showKeyboardShortcuts"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                KEYBOARD SHORTCUTS
              </button>
              <button 
                @click="showAbout"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                ABOUT
              </button>
              <div class="border-t border-[#5a548d] my-1"></div>
              <button 
                @click="checkForUpdates"
                class="block w-full text-left px-3 py-1 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10"
              >
                CHECK FOR UPDATES
              </button>
            </div>
          </button>
        </div>
      </div>
      
      <div class="flex items-center gap-4 text-xs">
        <div 
          :class="[
            'flex items-center gap-2',
            playerStore.isConnected ? 'text-green-400' : 'text-red-400'
          ]"
        >
          <div 
            :class="[
              'w-2 h-2 rounded-full animate-pulse',
              playerStore.isConnected ? 'bg-green-400' : 'bg-red-400'
            ]"
          ></div>
          <span>{{ playerStore.isConnected ? '[CONNECTED]' : '[DISCONNECTED]' }}</span>
        </div>
        <div class="text-[#837dbd]">
          {{ libraryStore.songs.length }} SONGS
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex h-[calc(100vh-12rem)]">
      <LeftSidebar 
        :current-view="currentView" 
        :playlists="libraryStore.playlists" 
        @change-view="changeView"
        @create-playlist="createPlaylist"
        @select-playlist="selectPlaylist" 
      />
      
      <!-- Main Content -->
      <div class="flex-1 overflow-hidden">
        <!-- Library View -->
        <div v-if="currentView === 'library'" class="h-full">
          <div class="p-4 border-b border-[#837dbd] bg-black">
            <div class="flex justify-between items-center">
              <h2 class="text-lg text-[#d3ceff] font-bold">LIBRARY</h2>
              <div class="flex items-center gap-3">
                <input
                  v-model="libraryStore.searchQuery"
                  type="text"
                  placeholder="Search songs..."
                  class="px-3 py-1 bg-black border border-[#837dbd] text-[#d3ceff] text-sm focus:outline-none focus:border-[#d3ceff] w-48"
                  @input="debouncedSearch"
                />
                <select
                  v-model="libraryStore.filters.sortBy"
                  class="px-2 py-1 bg-black border border-[#837dbd] text-[#837dbd] text-sm focus:outline-none focus:border-[#d3ceff]"
                  @change="libraryStore.fetchSongs()"
                >
                  <option value="title">Sort by Title</option>
                  <option value="artist">Sort by Artist</option>
                  <option value="album">Sort by Album</option>
                  <option value="rating">Sort by Rating</option>
                  <option value="play_count">Sort by Play Count</option>
                </select>
              </div>
            </div>
          </div>
          
          <div class="h-[calc(100%-4rem)] overflow-y-auto p-4">
            <!-- Loading State -->
            <div v-if="libraryStore.loading" class="flex items-center justify-center h-32">
              <div class="text-center">
                <div class="text-[#837dbd] text-xl font-mono animate-pulse">LOADING...</div>
              </div>
            </div>
            
            <!-- Error State -->
            <div v-else-if="libraryStore.error" class="flex items-center justify-center h-32">
              <div class="text-center">
                <div class="text-red-400 text-xl font-mono">ERROR</div>
                <div class="text-[#837dbd] text-sm mt-2">{{ libraryStore.error }}</div>
                <button
                  @click="libraryStore.fetchSongs()"
                  class="mt-3 px-4 py-2 border border-red-400 text-red-400 hover:bg-red-400 hover:text-black transition-all"
                >
                  RETRY
                </button>
              </div>
            </div>
            
            <!-- Songs List -->
            <div v-else class="space-y-2">
              <div
                v-for="song in libraryStore.filteredSongs"
                :key="song.id"
                class="flex items-center justify-between p-3 border border-[#5a548d] hover:border-[#837dbd] transition-all group"
                @dblclick="playSong(song)"
              >
                <div class="flex items-center gap-4 flex-1 min-w-0">
                  <div class="w-10 h-10 bg-[#5a548d] flex items-center justify-center">
                    <span class="text-[#d3ceff] text-lg">♪</span>
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="text-[#d3ceff] text-sm font-mono truncate">{{ song.title }}</div>
                    <div class="text-[#837dbd] text-xs font-mono truncate">{{ song.artist }} • {{ song.album }}</div>
                  </div>
                </div>
                
                <div class="flex items-center gap-6">
                  <!-- Rating -->
                  <div class="flex gap-1">
                    <button
                      v-for="star in 5"
                      :key="star"
                      @click.stop="rateSong(song, star * 2)"
                      class="text-xs"
                      :class="star * 2 <= song.rating ? 'text-yellow-400' : 'text-[#5a548d]'"
                    >
                      ★
                    </button>
                  </div>
                  
                  <!-- Duration -->
                  <div class="text-[#837dbd] text-xs font-mono w-16 text-right">
                    {{ formatDuration(song.duration) }}
                  </div>
                  
                  <!-- Actions -->
                  <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button
                      @click.stop="playSong(song)"
                      class="text-[#837dbd] hover:text-[#d3ceff] transition-colors text-sm"
                      title="Play"
                    >
                      ▶
                    </button>
                    <button
                      @click.stop="addToQueue(song)"
                      class="text-[#837dbd] hover:text-[#d3ceff] transition-colors text-sm"
                      title="Add to Queue"
                    >
                      +
                    </button>
                    <button
                      @click.stop="editSong(song)"
                      class="text-[#837dbd] hover:text-[#d3ceff] transition-colors text-sm"
                      title="Edit"
                    >
                      ✎
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- Empty State -->
              <div v-if="libraryStore.filteredSongs.length === 0" class="text-center py-8">
                <div class="text-[#837dbd] text-lg">No songs found</div>
                <div class="text-[#5a548d] text-sm mt-2">
                  {{ libraryStore.filters.search ? 'Try a different search term' : 'Upload some songs to get started' }}
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Playlists View -->
        <PlaylistsView v-else-if="currentView === 'playlists'" />
        
        <!-- Visualizer View -->
        <VisualizerView v-else-if="currentView === 'visualizer'" />
        
        <!-- Settings View -->
        <div v-else-if="currentView === 'settings'" class="h-full p-4">
          <div class="border-2 border-[#837dbd] p-4 bg-black">
            <h2 class="text-lg text-[#d3ceff] font-bold">SETTINGS</h2>
            <p class="text-[#837dbd] text-sm mt-1">Configure your MSX Player</p>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
            <div class="border border-[#837dbd] p-4 bg-black">
              <h3 class="text-[#d3ceff] text-sm font-bold mb-3">PLAYER SETTINGS</h3>
              <div class="space-y-3">
                <div>
                  <label class="text-[#837dbd] text-xs block mb-1">VOLUME</label>
                  <input
                    type="range"
                    v-model="playerStore.volume"
                    min="0"
                    max="1"
                    step="0.01"
                    @change="playerStore.setVolume(playerStore.volume)"
                    class="w-full"
                  />
                  <div class="text-[#5a548d] text-xs text-right mt-1">
                    {{ Math.round(playerStore.volume * 100) }}%
                  </div>
                </div>
                <div>
                  <label class="text-[#837dbd] text-xs block mb-1">PLAYBACK SPEED</label>
                  <input
                    type="range"
                    v-model="playerStore.speed"
                    min="0.5"
                    max="2"
                    step="0.1"
                    @change="playerStore.setSpeed(playerStore.speed)"
                    class="w-full"
                  />
                  <div class="text-[#5a548d] text-xs text-right mt-1">
                    {{ playerStore.speed.toFixed(1) }}x
                  </div>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-[#837dbd] text-xs">SHUFFLE</span>
                  <button
                    @click="playerStore.setShuffle(!playerStore.shuffle)"
                    :class="[
                      'w-8 h-4 rounded-full transition-all relative',
                      playerStore.shuffle ? 'bg-[#d3ceff]' : 'bg-[#5a548d]'
                    ]"
                  >
                    <div
                      :class="[
                        'absolute top-0.5 w-3 h-3 rounded-full transition-all',
                        playerStore.shuffle ? 'right-0.5 bg-black' : 'left-0.5 bg-white'
                      ]"
                    ></div>
                  </button>
                </div>
              </div>
            </div>
            
            <div class="border border-[#837dbd] p-4 bg-black">
              <h3 class="text-[#d3ceff] text-sm font-bold mb-3">LIBRARY SETTINGS</h3>
              <div class="space-y-3">
                <div>
                  <label class="text-[#837dbd] text-xs block mb-1">AUTO-SCAN FOR NEW FILES</label>
                  <button
                    @click="toggleAutoScan"
                    :class="[
                      'w-8 h-4 rounded-full transition-all relative',
                      autoScan ? 'bg-[#d3ceff]' : 'bg-[#5a548d]'
                    ]"
                  >
                    <div
                      :class="[
                        'absolute top-0.5 w-3 h-3 rounded-full transition-all',
                        autoScan ? 'right-0.5 bg-black' : 'left-0.5 bg-white'
                      ]"
                    ></div>
                  </button>
                </div>
                <div>
                  <label class="text-[#837dbd] text-xs block mb-1">DEFAULT SORT ORDER</label>
                  <select
                    v-model="libraryStore.filters.sortBy"
                    class="w-full px-2 py-1 bg-black border border-[#837dbd] text-[#837dbd] text-xs focus:outline-none focus:border-[#d3ceff]"
                  >
                    <option value="title">Title</option>
                    <option value="artist">Artist</option>
                    <option value="album">Album</option>
                    <option value="rating">Rating</option>
                  </select>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Right Sidebar -->
      <div class="w-64 border-l border-[#837dbd] bg-black overflow-y-auto">
        <div class="p-4 border-b border-[#837dbd]">
          <h3 class="text-[#d3ceff] text-sm font-bold">NOW PLAYING</h3>
        </div>
        
        <div v-if="playerStore.currentSong" class="p-4">
          <div class="text-center mb-4">
            <div class="w-32 h-32 mx-auto bg-[#5a548d] flex items-center justify-center mb-3">
              <span class="text-[#d3ceff] text-4xl">♪</span>
            </div>
            <div class="text-[#d3ceff] text-sm font-mono truncate">{{ playerStore.currentSong.title }}</div>
            <div class="text-[#837dbd] text-xs font-mono">{{ playerStore.currentSong.artist }}</div>
          </div>
          
          <!-- Progress -->
          <div class="mb-4">
            <div class="h-1 bg-[#5a548d] rounded-full overflow-hidden">
              <div 
                class="h-full bg-[#d3ceff] transition-all duration-300"
                :style="{ width: `${playerStore.progressPercentage}%` }"
              ></div>
            </div>
            <div class="flex justify-between text-[#5a548d] text-xs mt-1">
              <span>{{ formatTime(playerStore.currentTime) }}</span>
              <span>{{ formatTime(playerStore.duration) }}</span>
            </div>
          </div>
          
          <!-- Controls -->
          <div class="flex justify-center items-center gap-4 mb-4">
            <button
              @click="playerStore.previousSong()"
              class="text-[#837dbd] hover:text-[#d3ceff] text-lg"
              title="Previous"
            >
              ⏮
            </button>
            <button
              @click="playerStore.isPlaying ? playerStore.pause() : playerStore.playSong(playerStore.currentSong!)"
              class="text-[#d3ceff] hover:text-white text-2xl"
              title="Play/Pause"
            >
              {{ playerStore.isPlaying ? '⏸' : '▶' }}
            </button>
            <button
              @click="playerStore.nextSong()"
              class="text-[#837dbd] hover:text-[#d3ceff] text-lg"
              title="Next"
            >
              ⏭
            </button>
          </div>
        </div>
        
        <div v-else class="p-8 text-center">
          <div class="text-[#5a548d] text-lg mb-2">♪</div>
          <div class="text-[#837dbd] text-sm">No song playing</div>
          <div class="text-[#5a548d] text-xs mt-1">Select a song to begin</div>
        </div>
        
        <!-- Queue -->
        <div class="border-t border-[#5a548d] p-4">
          <h3 class="text-[#d3ceff] text-sm font-bold mb-3">QUEUE ({{ playerStore.queue.length }})</h3>
          <div class="space-y-2 max-h-48 overflow-y-auto">
            <div
              v-for="(song, index) in playerStore.queue"
              :key="song.id"
              :class="[
                'p-2 text-xs cursor-pointer transition-all',
                index === playerStore.currentQueueIndex
                  ? 'bg-[#837dbd] text-black'
                  : 'text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10'
              ]"
              @click="playerStore.playSong(song)"
            >
              <div class="font-mono truncate">{{ song.title }}</div>
              <div class="text-[#5a548d] truncate">{{ song.artist }}</div>
            </div>
            <div v-if="playerStore.queue.length === 0" class="text-center text-[#5a548d] text-xs py-4">
              Queue is empty
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Now Playing Bar -->
    <div class="h-12 border-t border-[#837dbd] bg-black flex items-center px-4">
      <div class="flex items-center gap-4 flex-1">
        <!-- Song Info -->
        <div v-if="playerStore.currentSong" class="flex items-center gap-3 min-w-0 flex-1">
          <div class="w-8 h-8 bg-[#5a548d] flex items-center justify-center">
            <span class="text-[#d3ceff]">♪</span>
          </div>
          <div class="min-w-0 flex-1">
            <div class="text-[#d3ceff] text-sm font-mono truncate">{{ playerStore.currentSong.title }}</div>
            <div class="text-[#837dbd] text-xs font-mono truncate">{{ playerStore.currentSong.artist }}</div>
          </div>
        </div>
        <div v-else class="text-[#5a548d] text-sm">No song selected</div>
        
        <!-- Mini Player Controls -->
        <div class="flex items-center gap-3">
          <button
            @click="playerStore.previousSong()"
            class="text-[#837dbd] hover:text-[#d3ceff]"
            title="Previous"
          >
            ⏮
          </button>
          <button
            @click="playerStore.isPlaying ? playerStore.pause() : playerStore.playSong(playerStore.currentSong!)"
            class="text-[#d3ceff] hover:text-white"
            title="Play/Pause"
          >
            {{ playerStore.isPlaying ? '⏸' : '▶' }}
          </button>
          <button
            @click="playerStore.nextSong()"
            class="text-[#837dbd] hover:text-[#d3ceff]"
            title="Next"
          >
            ⏭
          </button>
        </div>
      </div>
      
      <!-- Volume Control -->
      <div class="flex items-center gap-2 w-32">
        <span class="text-[#837dbd] text-xs">VOL</span>
        <input
          type="range"
          v-model="playerStore.volume"
          min="0"
          max="1"
          step="0.01"
          @change="playerStore.setVolume(playerStore.volume)"
          class="flex-1"
        />
      </div>
    </div>

    <!-- File Upload Input (Hidden) -->
    <input
      ref="fileInput"
      type="file"
      accept=".mp3"
      @change="handleFileUpload"
      class="hidden"
      multiple
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { usePlayerStore } from '../stores/player'
import { useLibraryStore } from '../stores/library'
import PlaylistsView from './views/PlaylistsView.vue'
import VisualizerView from './views/VisualizerView.vue'
import type { Song, Playlist } from '../types'

const playerStore = usePlayerStore()
const libraryStore = useLibraryStore()

// Refs
const currentView = ref('library')
const showFileMenu = ref(false)
const showEditMenu = ref(false)
const showViewMenu = ref(false)
const showControlsMenu = ref(false)
const showHelpMenu = ref(false)
const fileInput = ref<HTMLInputElement>()
const darkMode = ref(true)
const compactView = ref(false)
const autoScan = ref(true)
let searchTimeout: number | null = null

// Close menus when clicking outside
const closeMenus = () => {
  showFileMenu.value = false
  showEditMenu.value = false
  showViewMenu.value = false
  showControlsMenu.value = false
  showHelpMenu.value = false
}

const changeView = (view: string) => {
  currentView.value = view
  closeMenus()
  
  if (view === 'playlists') {
    libraryStore.fetchPlaylists()
  }
}

const createPlaylist = async () => {
  const name = prompt('Enter playlist name:')
  if (name) {
    try {
      await libraryStore.createPlaylist(name)
      await libraryStore.fetchPlaylists()
    } catch (err) {
      console.error('Error creating playlist:', err)
      alert('Failed to create playlist')
    }
  }
}

const selectPlaylist = (playlistId: number) => {
  currentView.value = 'playlists'
  // The PlaylistsView component will handle loading the playlist
}

const playSong = async (song: Song) => {
  try {
    await playerStore.playSong(song)
  } catch (err) {
    console.error('Error playing song:', err)
    alert('Failed to play song')
  }
}

const addToQueue = (song: Song) => {
  if (!playerStore.queue.some(s => s.id === song.id)) {
    playerStore.setQueue([...playerStore.queue, song])
  }
}

const clearQueue = () => {
  playerStore.setQueue([])
}

const rateSong = async (song: Song, rating: number) => {
  try {
    await playerStore.rateSong(song.id, rating)
    await libraryStore.updateSong(song.id, { rating })
  } catch (err) {
    console.error('Error rating song:', err)
  }
}

const editSong = async (song: Song) => {
  const title = prompt('Edit title:', song.title)
  const artist = prompt('Edit artist:', song.artist)
  const album = prompt('Edit album:', song.album)
  
  if (title !== null || artist !== null || album !== null) {
    try {
      await libraryStore.updateSong(song.id, {
        title: title || song.title,
        artist: artist || song.artist,
        album: album || song.album,
      })
    } catch (err) {
      console.error('Error updating song:', err)
      alert('Failed to update song')
    }
  }
}

const uploadSong = () => {
  if (fileInput.value) {
    fileInput.value.click()
  }
}

const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  
  if (files && files.length > 0) {
    for (let i = 0; i < files.length; i++) {
      const file = files[i]
      if (file.type === 'audio/mpeg') {
        try {
          await libraryStore.uploadSong(file)
        } catch (err) {
          console.error('Error uploading file:', err)
          alert(`Failed to upload ${file.name}: ${err instanceof Error ? err.message : 'Unknown error'}`)
        }
      } else {
        alert(`${file.name} is not an MP3 file`)
      }
    }
    
    // Clear the input
    target.value = ''
  }
}

const scanLibrary = async () => {
  try {
    // This would typically call a backend endpoint to scan the music folder
    alert('Library scan functionality would be implemented here')
    await libraryStore.fetchSongs()
  } catch (err) {
    console.error('Error scanning library:', err)
  }
}

const findDuplicates = () => {
  const duplicates = libraryStore.scanForDuplicates()
  if (duplicates.length > 0) {
    alert(`Found ${duplicates.length} duplicate songs. Check the console for details.`)
    console.log('Duplicates found:', duplicates)
  } else {
    alert('No duplicates found!')
  }
}

const toggleDarkMode = () => {
  darkMode.value = !darkMode.value
  if (darkMode.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

const toggleCompactView = () => {
  compactView.value = !compactView.value
  // This would adjust CSS classes or store preference
}

const toggleAutoScan = () => {
  autoScan.value = !autoScan.value
  // Save preference to localStorage
  localStorage.setItem('autoScan', autoScan.value.toString())
}

const showKeyboardShortcuts = () => {
  alert(`
Keyboard Shortcuts:
Space - Play/Pause
← → - Seek forward/backward
↑ ↓ - Volume up/down
N - Next song
P - Previous song
S - Shuffle
R - Repeat mode
F - Toggle fullscreen
ESC - Close dialogs
  `)
}

const showAbout = () => {
  alert(`
MSX Player v1.0.0
A retro-futuristic music player
Built with Vue.js, TypeScript, and Go
  `)
}

const checkForUpdates = () => {
  alert('Checking for updates...\n(This would check your package.json or a remote API)')
}

const exportLibrary = () => {
  alert('Export functionality would be implemented here')
}

const importLibrary = () => {
  alert('Import functionality would be implemented here')
}

const cleanLibrary = () => {
  alert('Clean library functionality would be implemented here')
}

const updateMetadata = () => {
  alert('Update metadata functionality would be implemented here')
}

const formatDuration = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const debouncedSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    libraryStore.fetchSongs()
  }, 500)
}

// Initialize
onMounted(() => {
  console.log('🚀 MSX Player initialized')
  
  // Load initial data
  libraryStore.fetchSongs()
  libraryStore.fetchPlaylists()
  
  // Load preferences
  const savedAutoScan = localStorage.getItem('autoScan')
  if (savedAutoScan !== null) {
    autoScan.value = savedAutoScan === 'true'
  }
  
  // Set up keyboard shortcuts
  document.addEventListener('keydown', handleKeydown)
  
  // Close menus when clicking outside
  document.addEventListener('click', closeMenus)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  document.removeEventListener('click', closeMenus)
})

const handleKeydown = (event: KeyboardEvent) => {
  // Prevent shortcuts when typing in inputs
  if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
    return
  }
  
  switch (event.key) {
    case ' ':
      event.preventDefault()
      if (playerStore.currentSong) {
        playerStore.isPlaying ? playerStore.pause() : playerStore.playSong(playerStore.currentSong)
      }
      break
    case 'ArrowRight':
      event.preventDefault()
      if (playerStore.currentSong) {
        const newTime = Math.min(playerStore.currentTime + 10, playerStore.duration)
        playerStore.seek(newTime)
      }
      break
    case 'ArrowLeft':
      event.preventDefault()
      if (playerStore.currentSong) {
        const newTime = Math.max(playerStore.currentTime - 10, 0)
        playerStore.seek(newTime)
      }
      break
    case 'ArrowUp':
      event.preventDefault()
      playerStore.setVolume(Math.min(playerStore.volume + 0.1, 1))
      break
    case 'ArrowDown':
      event.preventDefault()
      playerStore.setVolume(Math.max(playerStore.volume - 0.1, 0))
      break
    case 'n':
    case 'N':
      event.preventDefault()
      playerStore.nextSong()
      break
    case 'p':
    case 'P':
      event.preventDefault()
      playerStore.previousSong()
      break
    case 's':
    case 'S':
      event.preventDefault()
      playerStore.setShuffle(!playerStore.shuffle)
      break
    case 'r':
    case 'R':
      event.preventDefault()
      const modes = ['none', 'all', 'one'] as const
      const currentIndex = modes.indexOf(playerStore.repeatMode)
      const nextIndex = (currentIndex + 1) % modes.length
      playerStore.setRepeatMode(modes[nextIndex])
      break
    case 'f':
    case 'F':
      event.preventDefault()
      if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen()
      } else {
        document.exitFullscreen()
      }
      break
    case 'Escape':
      closeMenus()
      break
  }
}
</script>

<style>
/* Global styles */
.scanlines {
  background: linear-gradient(
    to bottom,
    transparent 50%,
    rgba(211, 206, 255, 0.05) 51%
  );
  background-size: 100% 4px;
}

::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #000;
  border: 1px solid #837dbd;
}

::-webkit-scrollbar-thumb {
  background: #837dbd;
}

::-webkit-scrollbar-thumb:hover {
  background: #d3ceff;
}

::selection {
  background: rgba(211, 206, 255, 0.3);
}

input[type="range"] {
  appearance: none;
  height: 2px;
  background: #5a548d;
}

input[type="range"]::-webkit-slider-thumb {
  appearance: none;
  width: 12px;
  height: 12px;
  background: #d3ceff;
  border-radius: 50%;
  cursor: pointer;
}

tr:hover {
  background-color: rgba(211, 206, 255, 0.05) !important;
}

/* Font for terminal look */
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;500;600&display=swap');

body {
  font-family: 'JetBrains Mono', monospace;
}

/* Ensure proper spacing for dropdowns */
button:has(+ div) {
  position: relative;
}

button + div {
  position: absolute;
}
</style>