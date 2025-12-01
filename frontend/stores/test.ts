// stores/test.ts
import { defineStore } from "pinia";

export const useTestStore = defineStore("test", {
  state: () => ({
    message: "Store is working!",
  }),
  actions: {
    updateMessage(newMessage: string) {
      this.message = newMessage;
    },
  },
});
