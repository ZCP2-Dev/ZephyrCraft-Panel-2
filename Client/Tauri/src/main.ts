import { createApp } from "vue";
import App from "./App.vue";
import router from './router';

// 创建全局终端消息总线
class TerminalBus {
  private listeners: Map<string, Function[]> = new Map();

  on(event: string, callback: Function) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event)!.push(callback);
  }

  off(event: string, callback: Function) {
    const callbacks = this.listeners.get(event);
    if (callbacks) {
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    }
  }

  emit(event: string, data: any) {
    const callbacks = this.listeners.get(event);
    if (callbacks) {
      callbacks.forEach(callback => {
        try {
          callback(data);
        } catch (error) {
          console.error('TerminalBus callback error:', error);
        }
      });
    }
  }

  clear() {
    this.listeners.clear();
  }
}

// 初始化全局终端消息总线
(window as any).__TERMINAL_BUS__ = new TerminalBus();

// 初始化全局系统监控消息总线
(window as any).__SYSTEM_BUS__ = new TerminalBus();

createApp(App).use(router).mount("#app");
