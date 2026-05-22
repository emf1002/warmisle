<!-- frontend/src/layouts/AuthLayout.vue -->
<template>
  <div class="auth-layout" :class="isDark ? 'auth-dark' : 'auth-light'">
    <!-- 暗色主题：星空背景 -->
    <template v-if="isDark">
      <div class="stars" ref="starsRef"></div>
      <div class="moon">
        <div class="moon-crater moon-crater-1"></div>
        <div class="moon-crater moon-crater-2"></div>
        <div class="moon-crater moon-crater-3"></div>
      </div>
      <div class="cloud cloud-1"></div>
      <div class="cloud cloud-2"></div>
      <div class="landscape">
        <svg class="landscape-svg" viewBox="0 0 1440 200" preserveAspectRatio="none">
          <path d="M0 200 L0 160 Q100 120 200 140 Q280 90 350 100 Q400 60 440 80 L460 75 L462 30 L468 30 L470 75 Q480 78 500 90 Q560 60 640 80 Q720 100 800 90 Q880 70 960 85 Q1040 100 1100 75 Q1150 50 1200 80 Q1300 110 1400 95 L1440 100 L1440 200Z" fill="#0B1120"/>
          <rect x="445" y="48" width="30" height="27" rx="2" fill="#162040"/>
          <path d="M440 50 L460 30 L480 50" fill="none" stroke="#1A2340" stroke-width="2"/>
          <rect x="452" y="55" width="8" height="8" rx="1" fill="rgba(229,168,75,0.4)"/>
          <rect x="464" y="55" width="8" height="8" rx="1" fill="rgba(229,168,75,0.3)"/>
          <circle cx="420" cy="72" r="12" fill="#0F1729"/>
          <circle cx="495" cy="68" r="10" fill="#0F1729"/>
          <circle cx="410" cy="78" r="8" fill="#0D1425"/>
        </svg>
      </div>
      <div class="lighthouse-glow"></div>
      <div class="water">
        <div class="water-line" style="bottom: 55px; --dur: 5s;"></div>
        <div class="water-line" style="bottom: 40px; --dur: 6s;"></div>
        <div class="water-line" style="bottom: 25px; --dur: 7s;"></div>
        <div class="water-line" style="bottom: 12px; --dur: 5.5s;"></div>
      </div>
    </template>

    <!-- 亮色主题：有机色块背景 -->
    <template v-else>
      <div class="bg-shapes">
        <div class="blob blob-1"></div>
        <div class="blob blob-2"></div>
        <div class="blob blob-3"></div>
        <div class="blob blob-4"></div>
      </div>
      <div class="bg-grid"></div>
      <div class="float-deco float-circle fc-1"></div>
      <div class="float-deco float-circle fc-2"></div>
      <div class="float-deco fc-3"></div>
      <div class="float-deco float-circle fc-4"></div>
      <div class="float-deco fc-5"></div>
    </template>

    <div class="auth-container" :class="{ 'auth-container-split': showBrand && !isMobile }">
      <!-- 亮色主题大屏：左侧品牌介绍 -->
      <div v-if="showBrand && !isMobile" class="brand-section">
        <div class="brand-badge">
          <span class="brand-badge-dot"></span>
          家庭数字中心
        </div>
        <div class="brand-heading">暖 屿</div>
        <div class="brand-heading-en">Warm Isle</div>
        <div class="brand-desc">
          全家人的记账、待办、愿望、交流。<br>
          一个温暖的地方，搞定所有家庭琐事。
        </div>
        <div class="brand-features">
          <div class="feature-item"><span class="feature-icon">💰</span><span>记账</span></div>
          <div class="feature-item"><span class="feature-icon">✅</span><span>待办</span></div>
          <div class="feature-item"><span class="feature-icon">💫</span><span>愿望</span></div>
          <div class="feature-item"><span class="feature-icon">💬</span><span>论坛</span></div>
        </div>
      </div>

      <!-- 右侧：卡片区 -->
      <div class="auth-card-wrapper">
        <div v-if="!showBrand || isMobile" class="auth-header-compact">
          <span class="auth-logo-icon">🏠</span>
          <h1 class="auth-title-compact">暖屿 · WarmIsle</h1>
        </div>
        <div class="auth-card">
          <slot />
        </div>
        <div class="auth-footer">
          <p>&copy; 2026 暖屿 · WarmIsle</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const isDark = computed(() => themeStore.theme === 'dark')

// 亮色主题在大屏显示品牌介绍
const showBrand = ref(window.innerWidth >= 860)
const isMobile = ref(window.innerWidth < 520)

function onResize() {
  showBrand.value = window.innerWidth >= 860
  isMobile.value = window.innerWidth < 520
}

const starsRef = ref<HTMLElement | null>(null)

function generateStars() {
  if (!starsRef.value) return
  for (let i = 0; i < 80; i++) {
    const star = document.createElement('div')
    star.className = 'star'
    star.style.left = Math.random() * 100 + '%'
    star.style.top = Math.random() * 70 + '%'
    star.style.setProperty('--dur', (2 + Math.random() * 4) + 's')
    star.style.setProperty('--max-opacity', String(0.3 + Math.random() * 0.7))
    star.style.animationDelay = Math.random() * 5 + 's'
    const size = (1 + Math.random() * 2) + 'px'
    star.style.width = size
    star.style.height = size
    starsRef.value.appendChild(star)
  }
}

onMounted(() => {
  window.addEventListener('resize', onResize)
  // 暗色主题生成星星
  nextTick(() => {
    if (isDark.value && starsRef.value) {
      generateStars()
    }
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
})
</script>

<style scoped>
/* ==================== 基础布局 ==================== */
.auth-layout {
  position: relative;
  width: 100vw;
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  padding: 16px;
}

.auth-light {
  background: var(--color-bg-layout);
}

.auth-dark {
  background: linear-gradient(180deg, #0B1120 0%, #0F1729 40%, #162040 100%);
}

/* ==================== 亮色主题：有机色块背景 ==================== */
.bg-shapes {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0;
  animation: blobIn 2s ease-out forwards;
}

.blob-1 { width: 500px; height: 500px; top: -15%; left: -10%; background: radial-gradient(circle, #FDDAC6 0%, transparent 70%); animation-delay: 0.2s; }
.blob-2 { width: 400px; height: 400px; top: 10%; right: -8%; background: radial-gradient(circle, #B8D4E8 0%, transparent 70%); animation-delay: 0.4s; }
.blob-3 { width: 350px; height: 350px; bottom: -10%; left: 20%; background: radial-gradient(circle, #A8D8CA 0%, transparent 70%); animation-delay: 0.6s; }
.blob-4 { width: 300px; height: 300px; bottom: 5%; right: 10%; background: radial-gradient(circle, #C5B3D9 0%, transparent 70%); animation-delay: 0.8s; }

@keyframes blobIn {
  from { opacity: 0; transform: scale(0.7); }
  to { opacity: 0.5; transform: scale(1); }
}

.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(61, 53, 48, 0.02) 1px, transparent 1px),
    linear-gradient(90deg, rgba(61, 53, 48, 0.02) 1px, transparent 1px);
  background-size: 48px 48px;
  pointer-events: none;
}

.float-deco {
  position: absolute;
  pointer-events: none;
  opacity: 0;
  animation: floatIn 1s ease-out forwards;
}

.float-circle { border-radius: 50%; border: 2px solid; }

.fc-1 { width: 20px; height: 20px; top: 15%; left: 12%; border-color: #F09888; animation-delay: 1.2s; animation-name: floatIn, floatBob1; animation-duration: 1s, 6s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 1.2s, 2.2s; }
.fc-2 { width: 14px; height: 14px; top: 25%; right: 18%; border-color: #6BBAA7; animation-delay: 1.4s; animation-name: floatIn, floatBob2; animation-duration: 1s, 7s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 1.4s, 2.4s; }
.fc-3 { width: 10px; height: 10px; bottom: 20%; left: 8%; background: #F7D08A; border-radius: 50%; animation-delay: 1.6s; animation-name: floatIn, floatBob3; animation-duration: 1s, 8s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 1.6s, 2.6s; }
.fc-4 { width: 16px; height: 16px; bottom: 30%; right: 10%; border-color: #C5B3D9; border-radius: 4px; transform: rotate(45deg); animation-delay: 1.8s; animation-name: floatIn, floatBob1; animation-duration: 1s, 9s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 1.8s, 2.8s; }
.fc-5 { width: 8px; height: 8px; top: 60%; left: 25%; background: #F09888; border-radius: 50%; animation-delay: 2s; animation-name: floatIn, floatBob2; animation-duration: 1s, 5s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 2s, 3s; }

@keyframes floatIn { from { opacity: 0; } to { opacity: 0.6; } }
@keyframes floatBob1 { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-12px); } }
@keyframes floatBob2 { 0%, 100% { transform: translateY(0) rotate(0deg); } 50% { transform: translateY(-8px) rotate(5deg); } }
@keyframes floatBob3 { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-15px); } }

/* ==================== 暗色主题：星空背景 ==================== */
.stars { position: absolute; inset: 0; pointer-events: none; }

.star {
  position: absolute;
  background: #FFF;
  border-radius: 50%;
  animation: twinkle var(--dur) ease-in-out infinite;
  opacity: 0;
}

@keyframes twinkle {
  0%, 100% { opacity: 0.1; transform: scale(0.8); }
  50% { opacity: var(--max-opacity); transform: scale(1.2); }
}

.moon {
  position: absolute;
  top: 8%;
  right: 12%;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 35%, #FFF8E7 0%, #F0D890 40%, #D4A850 100%);
  box-shadow: 0 0 40px rgba(240, 216, 144, 0.3), 0 0 80px rgba(240, 216, 144, 0.15);
  opacity: 0;
  animation: moonRise 2s ease-out 0.5s forwards;
}

.moon-crater { position: absolute; border-radius: 50%; background: rgba(180, 150, 80, 0.15); }
.moon-crater-1 { width: 15px; height: 15px; top: 20px; left: 25px; }
.moon-crater-2 { width: 10px; height: 10px; top: 40px; left: 15px; }
.moon-crater-3 { width: 8px; height: 8px; top: 30px; left: 45px; }

@keyframes moonRise { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }

.cloud { position: absolute; background: rgba(255, 255, 255, 0.03); border-radius: 100px; filter: blur(20px); }
.cloud-1 { width: 300px; height: 60px; top: 18%; left: -50px; animation: cloudDrift 60s linear infinite; }
.cloud-2 { width: 200px; height: 40px; top: 30%; right: -30px; animation: cloudDrift 80s linear infinite reverse; }

@keyframes cloudDrift { from { transform: translateX(-100%); } to { transform: translateX(calc(100vw + 100%)); } }

.landscape { position: absolute; bottom: 0; left: 0; width: 100%; height: 200px; pointer-events: none; }
.landscape-svg { position: absolute; bottom: 0; width: 100%; height: 100%; }

.lighthouse-glow {
  position: absolute;
  bottom: 130px;
  left: 50%;
  transform: translateX(-50%);
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(229, 168, 75, 0.15) 0%, transparent 70%);
  animation: glowPulse 4s ease-in-out infinite;
  pointer-events: none;
}

@keyframes glowPulse {
  0%, 100% { opacity: 0.6; transform: translateX(-50%) scale(1); }
  50% { opacity: 1; transform: translateX(-50%) scale(1.1); }
}

.water { position: absolute; bottom: 0; left: 0; width: 100%; height: 80px; background: linear-gradient(180deg, rgba(15, 23, 41, 0) 0%, rgba(20, 30, 55, 0.8) 100%); pointer-events: none; }

.water-line {
  position: absolute;
  width: 100%;
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, rgba(229, 168, 75, 0.15) 30%, rgba(229, 168, 75, 0.25) 50%, rgba(229, 168, 75, 0.15) 70%, transparent 100%);
  animation: waterShimmer var(--dur) ease-in-out infinite;
}

@keyframes waterShimmer {
  0%, 100% { opacity: 0.3; transform: scaleX(0.8); }
  50% { opacity: 0.8; transform: scaleX(1); }
}

/* ==================== 容器 ==================== */
.auth-container {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 420px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.auth-container-split {
  max-width: 860px;
  flex-direction: row;
  gap: 60px;
  align-items: center;
}

/* ==================== 亮色大屏品牌区 ==================== */
.brand-section {
  flex: 1;
  opacity: 0;
  animation: slideRight 1s cubic-bezier(0.22, 1, 0.36, 1) 0.4s forwards;
}

@keyframes slideRight {
  from { opacity: 0; transform: translateX(-30px); }
  to { opacity: 1; transform: translateX(0); }
}

.brand-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  background: rgba(232, 116, 97, 0.08);
  border: 1px solid rgba(232, 116, 97, 0.15);
  border-radius: 100px;
  font-size: 12px;
  font-weight: 500;
  color: var(--color-brand);
  margin-bottom: 28px;
  letter-spacing: 0.03em;
}

.brand-badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-brand);
  animation: dotPulse 2s ease-in-out infinite;
}

@keyframes dotPulse {
  0%, 100% { opacity: 0.4; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.3); }
}

.brand-heading {
  font-size: 48px;
  font-weight: 700;
  color: var(--color-text-primary);
  line-height: 1.2;
  margin-bottom: 8px;
  letter-spacing: 0.04em;
}

.brand-heading-en {
  font-size: 14px;
  font-weight: 300;
  color: var(--color-text-disabled);
  letter-spacing: 0.3em;
  text-transform: uppercase;
  margin-bottom: 24px;
}

.brand-desc {
  font-size: 16px;
  color: var(--color-text-secondary);
  line-height: 1.8;
  max-width: 340px;
  margin-bottom: 40px;
}

.brand-features {
  display: flex;
  gap: 24px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.feature-icon {
  font-size: 18px;
}

/* ==================== 卡片区 ==================== */
.auth-card-wrapper {
  flex: 0 0 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 400px;
}

.auth-header-compact {
  text-align: center;
  margin-bottom: 24px;
}

.auth-logo-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 8px;
}

.auth-title-compact {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: 2px;
}

[data-theme="dark"] .auth-title-compact {
  color: #E8E4DE;
}

.auth-card {
  width: 100%;
  background: var(--auth-card-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--auth-card-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--auth-card-shadow);
  padding: 36px 32px;
  opacity: 0;
  animation: cardAppear 0.9s cubic-bezier(0.22, 1, 0.36, 1) 0.6s forwards;
  position: relative;
}

@keyframes cardAppear {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 暗色主题卡片顶部高光线 */
[data-theme="dark"] .auth-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 60px;
  height: 3px;
  background: linear-gradient(90deg, transparent, var(--color-brand), transparent);
  border-radius: 0 0 3px 3px;
}

.auth-footer {
  margin-top: 24px;
}

.auth-footer p {
  font-size: 12px;
  color: var(--color-text-disabled);
  margin: 0;
}

[data-theme="dark"] .auth-footer p {
  color: rgba(255, 255, 255, 0.3);
}

/* ==================== 响应式 ==================== */
@media (max-width: 860px) {
  .auth-container-split {
    flex-direction: column;
    gap: 24px;
  }
  .brand-section { text-align: center; }
  .brand-desc { margin-left: auto; margin-right: auto; }
  .brand-features { justify-content: center; }
  .auth-card-wrapper { flex: none; }
}

@media (max-width: 520px) {
  .auth-card { padding: 28px 24px; }
  .moon { width: 50px; height: 50px; top: 5%; right: 8%; }
}
</style>
