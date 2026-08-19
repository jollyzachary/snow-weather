<script lang="ts">
  import { LoaderCircle, MapPin, Search, X } from '@lucide/svelte';
  import { searchLocations } from '$lib/api';
  import type { LocationResult } from '$lib/types';

  let {
    location,
    onselect,
  }: {
    location: LocationResult;
    onselect: (location: LocationResult) => void;
  } = $props();

  let query = $state('');
  let results = $state<LocationResult[]>([]);
  let open = $state(false);
  let loading = $state(false);
  let activeIndex = $state(0);
  let controller: AbortController | undefined;
  let debounce: ReturnType<typeof setTimeout> | undefined;

  function onInput(event: Event) {
    query = (event.currentTarget as HTMLInputElement).value;
    activeIndex = 0;
    if (debounce) clearTimeout(debounce);
    controller?.abort();

    if (query.trim().length < 2) {
      loading = false;
      results = [];
      open = false;
      return;
    }

    loading = true;
    debounce = setTimeout(async () => {
      controller = new AbortController();
      try {
        results = await searchLocations(query.trim(), controller.signal);
        open = true;
      } catch (error) {
        if (!(error instanceof DOMException && error.name === 'AbortError')) {
          results = [];
        }
      } finally {
        loading = false;
      }
    }, 260);
  }

  function choose(result: LocationResult) {
    onselect(result);
    query = '';
    results = [];
    open = false;
  }

  function clearSearch() {
    controller?.abort();
    if (debounce) clearTimeout(debounce);
    query = '';
    results = [];
    open = false;
    loading = false;
  }

  function onKeydown(event: KeyboardEvent) {
    if (!open || results.length === 0) return;
    if (event.key === 'ArrowDown') {
      event.preventDefault();
      activeIndex = Math.min(activeIndex + 1, results.length - 1);
    } else if (event.key === 'ArrowUp') {
      event.preventDefault();
      activeIndex = Math.max(activeIndex - 1, 0);
    } else if (event.key === 'Enter') {
      event.preventDefault();
      choose(results[activeIndex]);
    } else if (event.key === 'Escape') {
      open = false;
    }
  }

  function onFocusOut(event: FocusEvent) {
    const currentTarget = event.currentTarget as HTMLElement;
    const nextTarget = event.relatedTarget;
    if (!(nextTarget instanceof Node) || !currentTarget.contains(nextTarget)) {
      open = false;
    }
  }
</script>

<div class="search-wrap" onfocusout={onFocusOut} role="search">
  <div class:expanded={open} class="search-field">
    <Search size={18} strokeWidth={1.6} />
    <input
      aria-autocomplete="list"
      aria-controls="location-results"
      aria-expanded={open}
      aria-label="Search locations"
      autocomplete="off"
      onfocus={() => (open = results.length > 0)}
      oninput={onInput}
      onkeydown={onKeydown}
      placeholder={`${location.name}, ${location.admin1 ?? location.country}`}
      role="combobox"
      value={query}
    />
    {#if loading}
      <LoaderCircle class="spinner" size={17} strokeWidth={1.6} />
    {:else if query}
      <button aria-label="Clear search" class="clear" onclick={clearSearch} type="button">
        <X size={16} />
      </button>
    {/if}
  </div>

  {#if open}
    <div class="results" id="location-results" role="listbox">
      {#if results.length}
        {#each results as result, index (result.id)}
          <button
            aria-selected={index === activeIndex}
            class:active={index === activeIndex}
            onclick={() => choose(result)}
            onmouseenter={() => (activeIndex = index)}
            role="option"
            type="button"
          >
            <MapPin size={16} strokeWidth={1.5} />
            <span>
              <strong>{result.name}</strong>
              <small>{[result.admin1, result.country].filter(Boolean).join(', ')}</small>
            </span>
          </button>
        {/each}
      {:else if !loading}
        <p>No matching places</p>
      {/if}
    </div>
  {/if}
</div>

<style>
  .search-wrap {
    position: relative;
    width: min(100%, 420px);
    z-index: 20;
  }

  .search-field {
    height: 48px;
    display: grid;
    grid-template-columns: auto 1fr auto;
    gap: 12px;
    align-items: center;
    padding: 0 16px;
    border: 1px solid rgb(255 255 255 / 24%);
    border-radius: 999px;
    color: var(--ink);
    background: rgb(255 255 255 / 13%);
    box-shadow: inset 0 1px 0 rgb(255 255 255 / 15%);
    backdrop-filter: blur(18px) saturate(110%);
    transition:
      border-color 180ms ease,
      background 180ms ease;
  }

  .search-field:focus-within,
  .search-field.expanded {
    border-color: rgb(255 255 255 / 52%);
    background: rgb(255 255 255 / 19%);
  }

  input {
    min-width: 0;
    border: 0;
    outline: 0;
    color: inherit;
    background: transparent;
    font: 500 0.82rem/1 var(--font-sans);
    letter-spacing: 0.02em;
  }

  input::placeholder {
    color: rgb(246 244 238 / 74%);
  }

  .clear {
    display: grid;
    place-items: center;
    padding: 2px;
    border: 0;
    color: inherit;
    background: transparent;
    cursor: pointer;
  }

  :global(.spinner) {
    animation: spin 0.8s linear infinite;
  }

  .results {
    position: absolute;
    top: calc(100% + 9px);
    left: 0;
    right: 0;
    padding: 7px;
    border: 1px solid rgb(255 255 255 / 23%);
    border-radius: 20px;
    background: rgb(24 34 42 / 88%);
    box-shadow: 0 22px 60px rgb(2 9 14 / 32%);
    backdrop-filter: blur(28px) saturate(120%);
  }

  .results button {
    width: 100%;
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 11px;
    align-items: center;
    padding: 11px 12px;
    border: 0;
    border-radius: 14px;
    color: var(--ink);
    text-align: left;
    background: transparent;
    cursor: pointer;
  }

  .results button:hover,
  .results button.active {
    background: rgb(255 255 255 / 10%);
  }

  .results span {
    display: grid;
    gap: 3px;
  }
  .results strong {
    font-size: 0.8rem;
    font-weight: 650;
  }
  .results small {
    color: var(--muted);
    font-size: 0.68rem;
  }
  .results p {
    margin: 14px;
    color: var(--muted);
    font-size: 0.75rem;
  }

  @keyframes spin {
    to {
      rotate: 360deg;
    }
  }
</style>
