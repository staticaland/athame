import { describe, it, expect } from 'vitest';
import {
  pad2,
  parseTimeHM,
  parseDuration,
  formatHM,
  formatDuration,
  quantizeForW1Delay,
  W1_PROGRAM_PRESETS
} from './script.js';

describe('pad2', () => {
  it('should pad single digit numbers with zero', () => {
    expect(pad2(0)).toBe('00');
    expect(pad2(5)).toBe('05');
    expect(pad2(9)).toBe('09');
  });

  it('should not pad double digit numbers', () => {
    expect(pad2(10)).toBe('10');
    expect(pad2(23)).toBe('23');
    expect(pad2(59)).toBe('59');
  });
});

describe('parseTimeHM', () => {
  it('should parse valid time strings to minutes from midnight', () => {
    expect(parseTimeHM('00:00')).toBe(0);
    expect(parseTimeHM('12:00')).toBe(720);
    expect(parseTimeHM('23:59')).toBe(1439);
    expect(parseTimeHM('09:30')).toBe(570);
  });

  it('should return null for invalid time strings', () => {
    expect(parseTimeHM('')).toBe(null);
    expect(parseTimeHM('invalid')).toBe(null);
    expect(parseTimeHM('24:00')).toBe(null);
    expect(parseTimeHM('12:60')).toBe(null);
    expect(parseTimeHM('-1:00')).toBe(null);
  });
});

describe('parseDuration', () => {
  it('should parse duration strings to minutes', () => {
    expect(parseDuration('0:00')).toBe(0);
    expect(parseDuration('1:30')).toBe(90);
    expect(parseDuration('3:39')).toBe(219);
    expect(parseDuration('10:00')).toBe(600);
  });

  it('should return null for invalid duration strings', () => {
    expect(parseDuration('')).toBe(null);
    expect(parseDuration('invalid')).toBe(null);
    expect(parseDuration('1:60')).toBe(null);
    expect(parseDuration('-1:00')).toBe(null);
  });
});

describe('formatHM', () => {
  it('should format minutes to HH:MM string', () => {
    expect(formatHM(0)).toBe('00:00');
    expect(formatHM(720)).toBe('12:00');
    expect(formatHM(1439)).toBe('23:59');
    expect(formatHM(570)).toBe('09:30');
  });

  it('should wrap around 24 hours', () => {
    expect(formatHM(1440)).toBe('00:00'); // 24 hours = midnight
    expect(formatHM(1500)).toBe('01:00'); // 25 hours = 1am
  });
});

describe('formatDuration', () => {
  it('should format duration in minutes to human readable string', () => {
    expect(formatDuration(0)).toBe('0 min');
    expect(formatDuration(30)).toBe('30 min');
    expect(formatDuration(60)).toBe('1 h');
    expect(formatDuration(90)).toBe('1 h 30 min');
    expect(formatDuration(219)).toBe('3 h 39 min');
    expect(formatDuration(1440)).toBe('24 h');
  });
});

describe('quantizeForW1Delay', () => {
  it('should enforce minimum delay of 30 minutes', () => {
    expect(quantizeForW1Delay(0)).toBe(30);
    expect(quantizeForW1Delay(15)).toBe(30);
    expect(quantizeForW1Delay(29)).toBe(30);
  });

  it('should use 30-minute steps below 10 hours', () => {
    expect(quantizeForW1Delay(30)).toBe(30);
    expect(quantizeForW1Delay(44)).toBe(30); // rounds down
    expect(quantizeForW1Delay(45)).toBe(60); // rounds up (1.5 rounds to 2)
    expect(quantizeForW1Delay(60)).toBe(60);
    expect(quantizeForW1Delay(90)).toBe(90);
    expect(quantizeForW1Delay(120)).toBe(120);
    expect(quantizeForW1Delay(570)).toBe(570); // 9.5 hours
  });

  it('should use 1-hour steps from 10 to 24 hours', () => {
    expect(quantizeForW1Delay(600)).toBe(600); // 10 hours
    expect(quantizeForW1Delay(630)).toBe(660); // 10.5h rounds to 11h
    expect(quantizeForW1Delay(660)).toBe(660); // 11 hours
    expect(quantizeForW1Delay(720)).toBe(720); // 12 hours
    expect(quantizeForW1Delay(1440)).toBe(1440); // 24 hours
  });

  it('should enforce maximum delay of 24 hours', () => {
    expect(quantizeForW1Delay(1441)).toBe(1440);
    expect(quantizeForW1Delay(2000)).toBe(1440);
    expect(quantizeForW1Delay(10000)).toBe(1440);
  });
});

describe('W1_PROGRAM_PRESETS', () => {
  it('should be an array with at least one preset', () => {
    expect(Array.isArray(W1_PROGRAM_PRESETS)).toBe(true);
    expect(W1_PROGRAM_PRESETS.length).toBeGreaterThan(0);
  });

  it('should have presets with name and duration properties', () => {
    W1_PROGRAM_PRESETS.forEach(preset => {
      expect(preset).toHaveProperty('name');
      expect(preset).toHaveProperty('duration');
      expect(typeof preset.name).toBe('string');
      expect(typeof preset.duration).toBe('string');
      expect(preset.name.length).toBeGreaterThan(0);
      expect(preset.duration.length).toBeGreaterThan(0);
    });
  });

  it('should have valid parseable durations', () => {
    W1_PROGRAM_PRESETS.forEach(preset => {
      const parsed = parseDuration(preset.duration);
      expect(parsed).not.toBe(null);
      expect(parsed).toBeGreaterThan(0);
    });
  });

  it('should include common W1 programs', () => {
    const presetNames = W1_PROGRAM_PRESETS.map(p => p.name.toLowerCase());
    expect(presetNames.some(name => name.includes('cotton'))).toBe(true);
    expect(presetNames.some(name => name.includes('express') || name.includes('quick'))).toBe(true);
  });
});
