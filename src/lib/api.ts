import type { LocationResult, Units, WeatherBundle } from './types';

interface LocationEnvelope {
  results?: LocationResult[];
}

async function readJson<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(body.error ?? `Request failed with status ${response.status}`);
  }
  return response.json() as Promise<T>;
}

export async function searchLocations(
  query: string,
  signal?: AbortSignal
): Promise<LocationResult[]> {
  const params = new URLSearchParams({ q: query });
  const response = await fetch(`/api/locations?${params}`, { signal });
  const payload = await readJson<LocationEnvelope>(response);
  return payload.results ?? [];
}

export async function fetchWeather(
  location: Pick<LocationResult, 'latitude' | 'longitude'>,
  units: Units,
  signal?: AbortSignal
): Promise<WeatherBundle> {
  const params = new URLSearchParams({
    lat: String(location.latitude),
    lon: String(location.longitude),
    units
  });
  const response = await fetch(`/api/weather?${params}`, { signal });
  return readJson<WeatherBundle>(response);
}
