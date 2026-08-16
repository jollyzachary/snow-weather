<script lang="ts">
  import { Cloud, CloudFog, CloudLightning, CloudRain, CloudSnow, Sun } from '@lucide/svelte';
  import { describeWeather } from '$lib/weather';

  let {
    code,
    size = 24,
    strokeWidth = 1.5
  }: { code: number; size?: number; strokeWidth?: number } = $props();

  let descriptor = $derived(describeWeather(code));
</script>

{#if descriptor.scene === 'clear'}
  <Sun {size} {strokeWidth} />
{:else if descriptor.scene === 'rain'}
  <CloudRain {size} {strokeWidth} />
{:else if descriptor.scene === 'snow'}
  <CloudSnow {size} {strokeWidth} />
{:else if descriptor.scene === 'storm'}
  <CloudLightning {size} {strokeWidth} />
{:else if descriptor.scene === 'fog'}
  <CloudFog {size} {strokeWidth} />
{:else}
  <Cloud {size} {strokeWidth} />
{/if}
