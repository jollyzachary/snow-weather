<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Activity,
    Compass,
    Droplets,
    Eye,
    Gauge,
    RefreshCw,
    Sunrise,
    Sunset,
    Wind
  } from '@lucide/svelte';
  import LocationSearch from '$lib/components/LocationSearch.svelte';
  import WeatherGlyph from '$lib/components/WeatherGlyph.svelte';
  import WeatherScene from '$lib/components/WeatherScene.svelte';
  import { fetchWeather } from '$lib/api';
  import type { LocationResult, Units, WeatherBundle } from '$lib/types';
  import {
    aqiLabel,
    degreesToCompass,
    describeWeather,
    formatDay,
    formatHour,
    formatTime
  } from '$lib/weather';

  const DEFAULT_LOCATION: LocationResult = {
    id: 5259502,
    name: 'Lake Mills',
    latitude: 43.0814,
    longitude: -88.9118,
    country: 'United States',
    country_code: 'US',
    admin1: 'Wisconsin',
    timezone: 'America/Chicago'
  };

  let location = $state<LocationResult>(DEFAULT_LOCATION);
  let units = $state<Units>('imperial');
  let bundle = $state<WeatherBundle | null>(null);
  let loading = $state(true);
  let errorMessage = $state('');
  let requestController: AbortController | undefined;

  let weather = $derived(bundle?.weather);
  let current = $derived(weather?.current);
  let descriptor = $derived(describeWeather(current?.weather_code ?? 0));
  let isDay = $derived((current?.is_day ?? 1) === 1);
  let temperatureUnit = $derived(units === 'imperial' ? '°F' : '°C');
  let speedUnit = $derived(units === 'imperial' ? 'mph' : 'km/h');
  let visibilityUnit = $derived(units === 'imperial' ? 'mi' : 'km');

  let hourlyItems = $derived.by(() => {
    if (!weather || !current) return [];
    const start = Math.max(
      0,
      weather.hourly.time.findIndex((time) => time >= current.time)
    );
    return weather.hourly.time.slice(start, start + 12).map((time, offset) => {
      const index = start + offset;
      return {
        time,
        temperature: weather.hourly.temperature_2m[index],
        precipitation: weather.hourly.precipitation_probability[index],
        code: weather.hourly.weather_code[index]
      };
    });
  });

  let dailyItems = $derived.by(() => {
    if (!weather) return [];
    return weather.daily.time.slice(0, 7).map((date, index) => ({
      date,
      code: weather.daily.weather_code[index],
      high: weather.daily.temperature_2m_max[index],
      low: weather.daily.temperature_2m_min[index],
      precipitation: weather.daily.precipitation_probability_max[index]
    }));
  });

  let currentVisibility = $derived.by(() => {
    if (!weather || !current) return undefined;
    const index = Math.max(0, weather.hourly.time.findIndex((time) => time >= current.time));
    const visibilityValue = weather.hourly.visibility[index];
    if (visibilityValue == null) return undefined;
    return units === 'imperial' ? visibilityValue / 5280 : visibilityValue / 1000;
  });

  let currentAqi = $derived(bundle?.air_quality.current?.us_aqi);

  async function loadWeather(nextLocation = location, nextUnits = units) {
    requestController?.abort();
    requestController = new AbortController();
    loading = true;
    errorMessage = '';
    try {
      bundle = await fetchWeather(nextLocation, nextUnits, requestController.signal);
    } catch (error) {
      if (!(error instanceof DOMException && error.name === 'AbortError')) {
        errorMessage = error instanceof Error ? error.message : 'Weather service unavailable';
      }
    } finally {
      loading = false;
    }
  }

  function selectLocation(nextLocation: LocationResult) {
    location = nextLocation;
    localStorage.setItem('snow-location', JSON.stringify(nextLocation));
    void loadWeather(nextLocation, units);
  }

  function setUnits(nextUnits: Units) {
    if (nextUnits === units) return;
    units = nextUnits;
    localStorage.setItem('snow-units', nextUnits);
    void loadWeather(location, nextUnits);
  }

  onMount(() => {
    const savedUnits = localStorage.getItem('snow-units');
    if (savedUnits === 'imperial' || savedUnits === 'metric') units = savedUnits;

    const savedLocation = localStorage.getItem('snow-location');
    if (savedLocation) {
      try {
        location = JSON.parse(savedLocation) as LocationResult;
      } catch {
        localStorage.removeItem('snow-location');
      }
    }
    void loadWeather(location, units);
  });
</script>

<svelte:head>
  <title>Snow — Atmospheric Weather</title>
  <meta
    name="description"
    content="A cinematic weather instrument for current conditions, hourly forecasts, and air quality."
  />
</svelte:head>

<WeatherScene scene={descriptor.scene} {isDay} />

<div class:night={!isDay} class="app-shell">
  <header class="topbar">
    <a class="brand" href="/" aria-label="Snow weather home">
      <span class="brand-mark">S</span>
      <span>
        <strong>SNOW</strong>
        <small>Atmospheric weather</small>
      </span>
    </a>

    <LocationSearch {location} onselect={selectLocation} />

    <div class="unit-switch" aria-label="Temperature units" role="group">
      <button aria-pressed={units === 'imperial'} class:active={units === 'imperial'} onclick={() => setUnits('imperial')} type="button">°F</button>
      <button aria-pressed={units === 'metric'} class:active={units === 'metric'} onclick={() => setUnits('metric')} type="button">°C</button>
    </div>
  </header>

  <main>
    {#if errorMessage && !bundle}
      <section class="error-panel">
        <p>Atmospheric feed interrupted</p>
        <span>{errorMessage}</span>
        <button onclick={() => loadWeather()} type="button">
          <RefreshCw size={15} /> Retry
        </button>
      </section>
    {:else if weather && current}
      <section class="hero" aria-live="polite">
        {#if errorMessage}
          <div class="stale-notice" role="status">
            Live update unavailable. Showing the latest observation.
          </div>
        {/if}
        <div class="location-heading">
          <p class="eyebrow">{weather.timezone_abbreviation} · Updated {formatTime(current.time)}</p>
          <h1>{location.name}</h1>
          <p>{[location.admin1, location.country].filter(Boolean).join(', ')}</p>
        </div>

        <div class="current-reading">
          <div class="condition-glyph">
            <WeatherGlyph code={current.weather_code} size={42} strokeWidth={1.15} />
          </div>
          <div class="temperature">
            <span>{Math.round(current.temperature_2m)}</span><sup>°</sup>
          </div>
          <div class="condition-copy">
            <strong>{descriptor.label}</strong>
            <span>Feels like {Math.round(current.apparent_temperature)}{temperatureUnit}</span>
          </div>
        </div>

        <div class="hero-rule"></div>

        <div class="hero-stats">
          <div>
            <Droplets size={18} strokeWidth={1.5} />
            <span>Humidity</span>
            <strong>{current.relative_humidity_2m}%</strong>
          </div>
          <div>
            <Wind size={18} strokeWidth={1.5} />
            <span>Wind</span>
            <strong>{Math.round(current.wind_speed_10m)} {speedUnit}</strong>
          </div>
          <div>
            <Compass size={18} strokeWidth={1.5} />
            <span>Direction</span>
            <strong>{degreesToCompass(current.wind_direction_10m)}</strong>
          </div>
          <div>
            <Gauge size={18} strokeWidth={1.5} />
            <span>Pressure</span>
            <strong>{Math.round(current.surface_pressure)} hPa</strong>
          </div>
        </div>
      </section>

      <section class="forecast-panel">
        <div class="section-heading">
          <div>
            <p class="eyebrow">The next twelve hours</p>
            <h2>Today’s atmosphere</h2>
          </div>
          {#if bundle?.cached}
            <span class="cache-note">Cached observation</span>
          {/if}
        </div>

        <div class="hourly-strip">
          {#each hourlyItems as hour, index}
            <article class:now={index === 0}>
              <span>{index === 0 ? 'Now' : formatHour(hour.time)}</span>
              <WeatherGlyph code={hour.code} size={24} strokeWidth={1.35} />
              <strong>{Math.round(hour.temperature)}°</strong>
              <small>{hour.precipitation}%</small>
            </article>
          {/each}
        </div>
      </section>

      <section class="detail-grid">
        <article class="week-panel panel">
          <div class="section-heading compact">
            <div>
              <p class="eyebrow">Seven-day outlook</p>
              <h2>Week ahead</h2>
            </div>
          </div>

          <div class="week-list">
            {#each dailyItems as day, index}
              <div class="day-row">
                <strong>{index === 0 ? 'Today' : formatDay(day.date)}</strong>
                <span class="day-condition">
                  <WeatherGlyph code={day.code} size={21} strokeWidth={1.4} />
                  {describeWeather(day.code).shortLabel}
                </span>
                <span class="rain-chance">{day.precipitation}%</span>
                <span class="temperature-range">
                  <b>{Math.round(day.high)}°</b>
                  <i></i>
                  <em>{Math.round(day.low)}°</em>
                </span>
              </div>
            {/each}
          </div>
        </article>

        <div class="metric-column">
          <article class="metric-card panel air-card">
            <div class="metric-icon"><Activity size={20} strokeWidth={1.4} /></div>
            <p>Air quality</p>
            <strong>{currentAqi === undefined ? '—' : Math.round(currentAqi)}</strong>
            <span>{aqiLabel(currentAqi)} · US AQI</span>
            <div class="aqi-track"><i style:width={`${Math.min(100, (currentAqi ?? 0) / 3)}%`}></i></div>
          </article>

          <article class="metric-card panel">
            <div class="metric-icon"><Eye size={20} strokeWidth={1.4} /></div>
            <p>Visibility</p>
            <strong>{currentVisibility?.toFixed(1) ?? '—'}</strong>
            <span>{visibilityUnit} · Surface range</span>
          </article>

          <article class="sun-card panel">
            <div>
              <span><Sunrise size={19} strokeWidth={1.4} /> Sunrise</span>
              <strong>{formatTime(weather.daily.sunrise[0])}</strong>
            </div>
            <i class="sun-path"><b></b></i>
            <div>
              <span><Sunset size={19} strokeWidth={1.4} /> Sunset</span>
              <strong>{formatTime(weather.daily.sunset[0])}</strong>
            </div>
          </article>
        </div>
      </section>

      <footer>
        <span>Forecast data by Open-Meteo</span>
        <span>{weather.latitude.toFixed(3)}° / {weather.longitude.toFixed(3)}°</span>
      </footer>
    {:else}
      <section class="loading-state">
        <span></span>
        <p>Reading the atmosphere</p>
      </section>
    {/if}
  </main>

  {#if loading && bundle}
    <div class="refresh-indicator"><RefreshCw size={14} /> Updating</div>
  {/if}
</div>
