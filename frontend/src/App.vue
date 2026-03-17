<template>
  <div class="app-container">
    <nav class="navbar navbar-dark">
      <div class="nav-container">
        <router-link class="navbar-brand" to="/">LOAN</router-link>
        <div class="nav-links">
          <router-link class="nav-link" to="/">首页</router-link>
          <router-link class="nav-link" to="/loans">列表</router-link>
          <router-link class="nav-link" to="/add">新增</router-link>
        </div>
      </div>
    </nav>
    <router-view />
    <footer class="app-footer">
      <span class="footer-text">powered by jeffrey huang</span>
    </footer>
  </div>
</template>

<script lang="ts" setup>
</script>

<style>
/* 全局样式已在 style.css 中定义 */
.app-container {
  position: relative;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-footer {
  position: fixed;
  bottom: 0;
  left: 0;
  padding: 0.5rem 1rem;
  z-index: 1;
  pointer-events: none;
}

.footer-text {
  font-size: 0.625rem;
  color: rgba(255, 255, 255, 0.3);
  font-family: var(--font-mono);
  letter-spacing: 0.05em;
  text-transform: lowercase;
}

.navbar {
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  box-shadow: 0 4px 24px 0 rgba(0, 0, 0, 0.4), inset 0 -1px 0 0 rgba(255, 255, 255, 0.08);
  padding: 0.5rem 1rem !important;
}

.nav-container {
  display: flex !important;
  flex-direction: row !important;
  justify-content: space-between !important;
  align-items: center !important;
  width: 100% !important;
  max-width: 100% !important;
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
</style>
