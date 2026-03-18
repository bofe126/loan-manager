<template>
  <div class="app-container">
    <nav class="navbar navbar-dark">
      <div class="nav-container">
        <div class="brand-section">
          <router-link class="navbar-brand" to="/">
            <img src="/appicon.png" class="brand-icon" alt="icon" />
            贷款助手
          </router-link>
        </div>
        <div class="nav-links">
          <router-link class="nav-link" to="/">首页</router-link>
          <router-link class="nav-link" to="/loans">列表</router-link>
          <router-link class="nav-link" to="/add">新增</router-link>
        </div>
        <div class="drag-region"></div>
        <div class="window-controls">
          <button class="window-button minimize" @click="minimizeWindow" title="最小化">
            <svg width="12" height="12" viewBox="0 0 12 12">
              <rect x="0" y="5" width="12" height="2" fill="currentColor"/>
            </svg>
          </button>
          <button class="window-button maximize" @click="toggleMaximize" title="最大化">
            <svg width="12" height="12" viewBox="0 0 12 12">
              <rect x="0" y="0" width="12" height="12" fill="none" stroke="currentColor" stroke-width="1.5"/>
            </svg>
          </button>
          <button class="window-button close" @click="closeWindow" title="关闭">
            <svg width="12" height="12" viewBox="0 0 12 12">
              <path d="M1 1 L11 11 M11 1 L1 11" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
          </button>
        </div>
      </div>
    </nav>
    <div class="main-content">
      <router-view />
    </div>
    <footer v-if="$route.path === '/'" class="app-footer">
      <span class="footer-text">powered by jeffrey huang</span>
      <div class="github-link" @click="openGitHub">
        <svg height="16" width="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
        </svg>
        GitHub
      </div>
    </footer>
  </div>
</template>

<script lang="ts" setup>
import { WindowMinimise, WindowToggleMaximise, Quit, BrowserOpenURL } from '../wailsjs/runtime/runtime';

const minimizeWindow = () => {
  WindowMinimise();
};

const toggleMaximize = () => {
  WindowToggleMaximise();
};

const closeWindow = () => {
  Quit();
};

const openGitHub = () => {
  BrowserOpenURL('https://github.com/bofe126/loan-manager');
};
</script>

<style>
/* 全局样式已在 style.css 中定义 */
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

.app-footer {
  padding: 0.5rem 1rem;
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-text {
  font-size: 0.625rem;
  color: rgba(255, 255, 255, 0.3);
  font-family: var(--font-mono);
  letter-spacing: 0.05em;
  text-transform: lowercase;
}

.app-footer .github-link {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.25rem 0.5rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.625rem;
  font-weight: 500;
  transition: all 0.2s;
  cursor: pointer;
}

.app-footer .github-link:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  color: rgba(255, 255, 255, 0.9);
}

.navbar {
  flex-shrink: 0;
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  box-shadow: 0 4px 24px 0 rgba(0, 0, 0, 0.4), inset 0 -1px 0 0 rgba(255, 255, 255, 0.08);
  padding: 0.5rem 1rem !important;
  --wails-draggable: drag;
}

.nav-container {
  display: flex !important;
  flex-direction: row !important;
  justify-content: space-between !important;
  align-items: center !important;
  width: 100% !important;
  max-width: 100% !important;
  gap: 1rem !important;
}

.brand-section {
  flex-shrink: 0;
  --wails-draggable: no-drag;
}

.navbar-brand {
  font-family: var(--font-heading) !important;
  font-size: 1.125rem !important;
  font-weight: 700 !important;
  letter-spacing: 0.1em !important;
  background: linear-gradient(135deg, #3399ff 0%, #0080ff 50%, #0066cc 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  position: relative;
  margin: 0 !important;
  padding: 0 !important;
  flex-shrink: 0 !important;
  display: flex !important;
  align-items: center !important;
  gap: 0.5rem !important;
  -webkit-app-region: no-drag;
  app-region: no-drag;
}

.brand-icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.navbar-brand::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #0080ff, transparent);
  opacity: 0;
  transition: opacity var(--transition-base);
}

.navbar-brand:hover::after {
  opacity: 1;
}

.nav-links {
  display: flex !important;
  flex-direction: row !important;
  gap: 0.25rem !important;
  margin: 0 !important;
  padding: 0 !important;
  list-style: none !important;
  flex-shrink: 0 !important;
  --wails-draggable: no-drag;
}

.drag-region {
  flex: 1;
  min-width: 100px;
  --wails-draggable: drag;
}

.nav-link {
  font-weight: 500 !important;
  font-size: 0.8125rem !important;
  letter-spacing: 0.05em !important;
  text-transform: uppercase !important;
  position: relative !important;
  padding: 0.4rem 0.6rem !important;
  white-space: nowrap !important;
  display: inline-block !important;
  margin: 0 !important;
}

.router-link-active {
  color: var(--primary-light) !important;
  font-weight: 600 !important;
}

.router-link-active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 60%;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--primary), transparent);
  box-shadow: 0 0 8px rgba(0, 128, 255, 0.6);
}

.window-controls {
  display: flex;
  gap: 8px;
  margin-left: auto;
  --wails-draggable: no-drag;
}

.window-button {
  width: 32px;
  height: 24px;
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
  border-radius: 4px;
}

.window-button:hover {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.9);
}

.window-button.close:hover {
  background: #e81123;
  color: white;
}

.window-button:active {
  transform: scale(0.95);
}

@media (max-width: 768px) {
  .nav-container {
    gap: 0.5rem !important;
  }

  .navbar-brand {
    font-size: 0.875rem !important;
  }

  .brand-icon {
    width: 20px;
    height: 20px;
  }

  .nav-link {
    font-size: 0.75rem !important;
    padding: 0.3rem 0.4rem !important;
  }

  .window-button {
    width: 28px;
    height: 22px;
  }
}

@media (max-width: 600px) {
  .navbar {
    padding: 0.375rem 0.75rem !important;
  }

  .nav-container {
    gap: 0.25rem !important;
  }

  .navbar-brand {
    font-size: 0.75rem !important;
    gap: 0.25rem !important;
  }

  .brand-icon {
    width: 18px;
    height: 18px;
  }

  .nav-link {
    font-size: 0.7rem !important;
    padding: 0.25rem 0.3rem !important;
  }

  .window-button {
    width: 24px;
    height: 20px;
  }

  .window-button svg {
    width: 10px;
    height: 10px;
  }
}
</style>
