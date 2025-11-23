import { describe, it, expect } from 'vitest';
import {
  pad2,
  parseTimeHM,
  parseDuration,
  formatHM,
  formatDuration,
  quantizeDelay
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

describe('quantizeDelay', () => {
  it('should enforce minimum delay of 30 minutes', () => {
    expect(quantizeDelay(0)).toBe(30);
    expect(quantizeDelay(15)).toBe(30);
    expect(quantizeDelay(29)).toBe(30);
  });

  it('should use 30-minute steps below 10 hours', () => {
    expect(quantizeDelay(30)).toBe(30);
    expect(quantizeDelay(44)).toBe(30); // rounds down
    expect(quantizeDelay(45)).toBe(60); // rounds up (1.5 rounds to 2)
    expect(quantizeDelay(60)).toBe(60);
    expect(quantizeDelay(90)).toBe(90);
    expect(quantizeDelay(120)).toBe(120);
    expect(quantizeDelay(570)).toBe(570); // 9.5 hours
  });

  it('should use 1-hour steps from 10 to 24 hours', () => {
    expect(quantizeDelay(600)).toBe(600); // 10 hours
    expect(quantizeDelay(630)).toBe(660); // 10.5h rounds to 11h
    expect(quantizeDelay(660)).toBe(660); // 11 hours
    expect(quantizeDelay(720)).toBe(720); // 12 hours
    expect(quantizeDelay(1440)).toBe(1440); // 24 hours
  });

  it('should enforce maximum delay of 24 hours', () => {
    expect(quantizeDelay(1441)).toBe(1440);
    expect(quantizeDelay(2000)).toBe(1440);
    expect(quantizeDelay(10000)).toBe(1440);
  });
});
