export type Units = 'imperial' | 'metric';

export interface LocationResult {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
  country: string;
  country_code: string;
  admin1?: string;
  timezone: string;
}

export interface CurrentConditions {
  time: string;
  interval: number;
  temperature_2m: number;
  relative_humidity_2m: number;
  apparent_temperature: number;
  is_day: number;
  precipitation: number;
  rain: number;
  showers: number;
  snowfall: number;
  weather_code: number;
  cloud_cover: number;
  surface_pressure: number;
  wind_speed_10m: number;
  wind_direction_10m: number;
  wind_gusts_10m: number;
}

export interface HourlyForecast {
  time: string[];
  temperature_2m: number[];
  apparent_temperature: number[];
  precipitation_probability: number[];
  weather_code: number[];
  wind_speed_10m: number[];
  uv_index: number[];
  visibility: Array<number | null>;
}

export interface DailyForecast {
  time: string[];
  weather_code: number[];
  temperature_2m_max: number[];
  temperature_2m_min: number[];
  sunrise: string[];
  sunset: string[];
  uv_index_max: number[];
  precipitation_probability_max: number[];
  wind_speed_10m_max: number[];
}

export interface WeatherResponse {
  latitude: number;
  longitude: number;
  timezone: string;
  timezone_abbreviation: string;
  utc_offset_seconds: number;
  current: CurrentConditions;
  hourly: HourlyForecast;
  daily: DailyForecast;
  current_units: Record<string, string>;
  hourly_units: Record<string, string>;
  daily_units: Record<string, string>;
}

export interface AirQualityResponse {
  current?: {
    time: string;
    us_aqi?: number;
    pm2_5?: number;
  };
}

export interface WeatherBundle {
  weather: WeatherResponse;
  air_quality: AirQualityResponse;
  cached: boolean;
  generated_at: string;
}

export interface WeatherDescriptor {
  label: string;
  shortLabel: string;
  scene: 'clear' | 'clouds' | 'rain' | 'snow' | 'storm' | 'fog';
}
