<script lang="ts">
  import type { WeatherDescriptor } from '$lib/types';

  let { scene, isDay }: { scene: WeatherDescriptor['scene']; isDay: boolean } = $props();

  const rainDrops = Array.from({ length: 34 });
  const snowFlakes = Array.from({ length: 28 });
  const stars = Array.from({ length: 24 });
</script>

<div class:night={!isDay} class="weather-scene {scene}" aria-hidden="true">
  <div class="sky-glow"></div>
  <div class="sun-or-moon"></div>

  <div class="stars">
    {#each stars as _, index (index)}
      <i style={`--i: ${index}`}></i>
    {/each}
  </div>

  <div class="cloud cloud-one"></div>
  <div class="cloud cloud-two"></div>
  <div class="cloud cloud-three"></div>

  <div class="rain-field">
    {#each rainDrops as _, index (index)}
      <i style={`--i: ${index}`}></i>
    {/each}
  </div>

  <div class="snow-field">
    {#each snowFlakes as _, index (index)}
      <i style={`--i: ${index}`}></i>
    {/each}
  </div>

  <div class="storm-flash"></div>
  <div class="atmosphere-grain"></div>
</div>

<style>
  .weather-scene {
    --sky-top: #829eb3;
    --sky-mid: #b8c8d2;
    --sky-low: #e5ddd0;
    position: fixed;
    inset: 0;
    z-index: -2;
    overflow: hidden;
    background:
      radial-gradient(circle at 68% 14%, rgb(255 244 220 / 58%), transparent 27%),
      linear-gradient(160deg, var(--sky-top), var(--sky-mid) 47%, var(--sky-low));
    transition: background 1.2s ease;
  }

  .weather-scene.clear {
    --sky-top: #6c94b4;
    --sky-mid: #b6cfdb;
    --sky-low: #f2dfbf;
  }

  .weather-scene.clouds {
    --sky-top: #657681;
    --sky-mid: #9ba7aa;
    --sky-low: #d2c9bc;
  }

  .weather-scene.rain,
  .weather-scene.storm {
    --sky-top: #273943;
    --sky-mid: #53656c;
    --sky-low: #8d9694;
  }

  .weather-scene.snow {
    --sky-top: #8296a1;
    --sky-mid: #c1ccd0;
    --sky-low: #ece8df;
  }

  .weather-scene.fog {
    --sky-top: #8e9898;
    --sky-mid: #bfc3bf;
    --sky-low: #ddd8ce;
  }

  .weather-scene.night {
    --sky-top: #07111d;
    --sky-mid: #172739;
    --sky-low: #46515d;
  }

  .sky-glow {
    position: absolute;
    width: 70vw;
    height: 70vw;
    top: -35vw;
    right: -10vw;
    border-radius: 50%;
    background: rgb(255 234 194 / 26%);
    filter: blur(30px);
  }

  .sun-or-moon {
    position: absolute;
    top: 12%;
    right: 11%;
    width: clamp(72px, 9vw, 144px);
    aspect-ratio: 1;
    border-radius: 50%;
    background: #fff1cb;
    box-shadow: 0 0 80px 20px rgb(255 226 163 / 38%);
    opacity: 0.76;
  }

  .night .sun-or-moon {
    background: #d9e2e8;
    box-shadow: 0 0 46px 9px rgb(202 218 232 / 18%);
    opacity: 0.58;
  }

  .cloud {
    position: absolute;
    width: 50vw;
    height: 17vw;
    min-height: 130px;
    border-radius: 50%;
    background: rgb(225 231 232 / 34%);
    filter: blur(20px);
    opacity: 0;
  }

  .cloud-one {
    top: 18%;
    left: -18%;
    animation: cloud-drift 34s linear infinite;
  }

  .cloud-two {
    top: 40%;
    left: 44%;
    transform: scale(0.8);
    animation: cloud-drift 43s linear -18s infinite;
  }

  .cloud-three {
    top: 64%;
    left: 2%;
    transform: scale(1.15);
    animation: cloud-drift 52s linear -31s infinite;
  }

  .clouds .cloud,
  .rain .cloud,
  .snow .cloud,
  .storm .cloud,
  .fog .cloud {
    opacity: 0.74;
  }

  .rain .cloud,
  .storm .cloud {
    background: rgb(42 56 65 / 64%);
  }

  .fog .cloud {
    opacity: 0.95;
    filter: blur(42px);
  }

  .rain-field,
  .snow-field,
  .stars {
    position: absolute;
    inset: 0;
    opacity: 0;
  }

  .rain .rain-field,
  .storm .rain-field,
  .snow .snow-field,
  .night .stars {
    opacity: 1;
  }

  .rain-field i {
    position: absolute;
    top: -16%;
    left: calc((var(--i) + 1) * 3%);
    width: 1px;
    height: clamp(80px, 9vw, 150px);
    background: linear-gradient(transparent, rgb(222 239 246 / 58%));
    transform: rotate(13deg);
    animation: rainfall 1.1s linear infinite;
  }

  .rain-field i:nth-child(3n) {
    animation-duration: 0.84s;
  }
  .rain-field i:nth-child(4n) {
    animation-delay: -0.62s;
    opacity: 0.45;
  }
  .rain-field i:nth-child(5n) {
    animation-duration: 1.35s;
  }

  .snow-field i {
    position: absolute;
    top: -4%;
    left: calc((var(--i) + 1) * 3.5%);
    width: 5px;
    height: 5px;
    border-radius: 50%;
    background: rgb(255 255 255 / 78%);
    box-shadow: 0 0 8px rgb(255 255 255 / 46%);
    animation: snowfall 8s linear infinite;
  }

  .snow-field i:nth-child(odd) {
    animation-duration: 11s;
    opacity: 0.52;
  }
  .snow-field i:nth-child(3n) {
    animation-delay: -5s;
    transform: scale(1.6);
  }
  .snow-field i:nth-child(4n) {
    animation-delay: -8s;
  }
  .stars i {
    position: absolute;
    left: calc((var(--i) + 1) * 4%);
    top: calc(7% + var(--i) * 2.2%);
    width: 2px;
    height: 2px;
    border-radius: 50%;
    background: white;
    box-shadow: 0 0 7px white;
    animation: twinkle 4s ease-in-out infinite alternate;
  }

  .stars i:nth-child(3n) {
    top: 22%;
    animation-delay: -2s;
  }
  .stars i:nth-child(4n) {
    top: 37%;
    opacity: 0.45;
  }
  .stars i:nth-child(5n) {
    top: 51%;
    animation-delay: -1s;
  }

  .storm-flash {
    position: absolute;
    inset: 0;
    background: #e8f3ff;
    opacity: 0;
  }

  .storm .storm-flash {
    animation: lightning 9s linear infinite;
  }

  .atmosphere-grain {
    position: absolute;
    inset: 0;
    opacity: 0.12;
    background-image:
      repeating-radial-gradient(
        circle at 20% 30%,
        transparent 0 1px,
        rgb(255 255 255 / 12%) 1.5px 2px
      ),
      repeating-radial-gradient(circle at 80% 70%, transparent 0 1px, rgb(0 0 0 / 10%) 1.5px 2px);
    background-size:
      7px 9px,
      11px 13px;
    mix-blend-mode: soft-light;
  }

  @keyframes cloud-drift {
    from {
      translate: -10vw 0;
    }
    to {
      translate: 85vw 0;
    }
  }

  @keyframes rainfall {
    from {
      translate: -8vw -20vh;
    }
    to {
      translate: 8vw 130vh;
    }
  }

  @keyframes snowfall {
    from {
      translate: -2vw -8vh;
      rotate: 0deg;
    }
    50% {
      translate: 3vw 52vh;
    }
    to {
      translate: -1vw 112vh;
      rotate: 240deg;
    }
  }

  @keyframes twinkle {
    from {
      opacity: 0.25;
      scale: 0.75;
    }
    to {
      opacity: 0.85;
      scale: 1.35;
    }
  }

  @keyframes lightning {
    0%,
    88%,
    92%,
    100% {
      opacity: 0;
    }
    89% {
      opacity: 0.32;
    }
    90% {
      opacity: 0.04;
    }
    91% {
      opacity: 0.22;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .weather-scene *,
    .weather-scene *::before,
    .weather-scene *::after {
      animation: none !important;
    }
  }
</style>
