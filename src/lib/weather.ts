import type { WeatherDescriptor } from './types';

export function describeWeather(code: number): WeatherDescriptor {
  if (code === 0) return { label: 'Clear sky', shortLabel: 'Clear', scene: 'clear' };
  if (code <= 3) return { label: 'Passing clouds', shortLabel: 'Clouds', scene: 'clouds' };
  if (code === 45 || code === 48) {
    return { label: 'Low visibility', shortLabel: 'Fog', scene: 'fog' };
  }
  if (code >= 51 && code <= 67) {
    return { label: 'Steady rain', shortLabel: 'Rain', scene: 'rain' };
  }
  if (code >= 71 && code <= 77) {
    return { label: 'Falling snow', shortLabel: 'Snow', scene: 'snow' };
  }
  if (code >= 80 && code <= 82) {
    return { label: 'Rain showers', shortLabel: 'Showers', scene: 'rain' };
  }
  if (code >= 85 && code <= 86) {
    return { label: 'Snow showers', shortLabel: 'Snow', scene: 'snow' };
  }
  if (code >= 95) {
    return { label: 'Thunderstorms', shortLabel: 'Storm', scene: 'storm' };
  }
  return { label: 'Changing conditions', shortLabel: 'Variable', scene: 'clouds' };
}

export function degreesToCompass(degrees: number): string {
  const directions = ['N', 'NE', 'E', 'SE', 'S', 'SW', 'W', 'NW'];
  return directions[Math.round(degrees / 45) % directions.length];
}

export function aqiLabel(value?: number): string {
  if (value === undefined) return 'Unavailable';
  if (value <= 50) return 'Good';
  if (value <= 100) return 'Moderate';
  if (value <= 150) return 'Sensitive';
  if (value <= 200) return 'Unhealthy';
  if (value <= 300) return 'Very unhealthy';
  return 'Hazardous';
}

export function formatHour(value: string): string {
  const hour = Number(value.slice(11, 13));
  if (!Number.isFinite(hour)) return value;
  const displayHour = hour % 12 || 12;
  return `${displayHour} ${hour < 12 ? 'AM' : 'PM'}`;
}

export function formatDay(value: string): string {
  return new Intl.DateTimeFormat('en-US', {
    weekday: 'short',
  }).format(new Date(`${value}T12:00:00`));
}

export function formatTime(value: string): string {
  const hour = Number(value.slice(11, 13));
  const minute = value.slice(14, 16);
  if (!Number.isFinite(hour) || minute.length !== 2) return value;
  const displayHour = hour % 12 || 12;
  return `${displayHour}:${minute} ${hour < 12 ? 'AM' : 'PM'}`;
}
