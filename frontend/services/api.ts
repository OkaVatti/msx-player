const API_BASE = "http://localhost:1323/api";

export const apiService = {
  async getSongs(): Promise<any[]> {
    const response = await fetch(`${API_BASE}/songs`);
    if (!response.ok) throw new Error("Failed to fetch songs");
    return response.json();
  },

  async uploadSong(formData: FormData): Promise<any> {
    const response = await fetch(`${API_BASE}/songs`, {
      method: "POST",
      body: formData,
    });
    if (!response.ok) throw new Error("Upload failed");
    return response.json();
  },

  async deleteSong(songId: number): Promise<void> {
    const response = await fetch(`${API_BASE}/songs/${songId}`, {
      method: "DELETE",
    });
    if (!response.ok) throw new Error("Delete failed");
  },

  async rateSong(songId: number, rating: number): Promise<void> {
    const response = await fetch(`${API_BASE}/songs/${songId}/rate`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ rating }),
    });
    if (!response.ok) throw new Error("Rating failed");
  },
};
