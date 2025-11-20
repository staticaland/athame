export function pad2(n) {
  return n.toString().padStart(2, "0");
}

export function parseTimeHM(str) {
  // "HH:MM" -> minutes from midnight
  if (!str || !str.includes(":")) return null;
  const [hStr, mStr] = str.split(":");
  const h = Number(hStr);
  const m = Number(mStr);
  if (
    Number.isNaN(h) || Number.isNaN(m) ||
    h < 0 || h > 23 ||
    m < 0 || m > 59
  ) return null;
  return h * 60 + m;
}

export function parseDuration(str) {
  // "H:MM" or "HH:MM" -> minutes
  if (!str || !str.includes(":")) return null;
  const [hStr, mStr] = str.split(":");
  const h = Number(hStr);
  const m = Number(mStr);
  if (
    Number.isNaN(h) || Number.isNaN(m) ||
    h < 0 || m < 0 || m > 59
  ) return null;
  return h * 60 + m;
}

export function formatHM(minutes) {
  minutes = ((minutes % (24 * 60)) + (24 * 60)) % (24 * 60); // wrap 0–1439
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  return pad2(h) + ":" + pad2(m);
}

export function formatDuration(min) {
  const h = Math.floor(min / 60);
  const m = min % 60;
  const parts = [];
  if (h) parts.push(h + " h");
  if (m) parts.push(m + " min");
  if (!parts.length) return "0 min";
  return parts.join(" ");
}

export function quantizeForW1Delay(delayMinutes) {
  const MIN_DELAY = 30;         // 30 min
  const MAX_DELAY = 24 * 60;    // 24 h

  if (delayMinutes < MIN_DELAY) delayMinutes = MIN_DELAY;
  if (delayMinutes > MAX_DELAY) delayMinutes = MAX_DELAY;

  const hours = delayMinutes / 60;

  let step;
  if (hours < 10) {
    step = 30; // 30-min steps
  } else {
    step = 60; // 1-hour steps
  }

  const rounded = Math.round(delayMinutes / step) * step;
  const snapped = Math.min(MAX_DELAY, Math.max(MIN_DELAY, rounded));

  return snapped;
}

function calculate() {
  const errorEl = document.getElementById("error");
  const resEl = document.getElementById("results");
  errorEl.textContent = "";
  resEl.innerHTML = "";

  const currentStr = document.getElementById("currentTime").value;
  const durationStr = document.getElementById("duration").value.trim();
  const finishStr = document.getElementById("finishTime").value;

  const nowMin = parseTimeHM(currentStr);
  const durMin = parseDuration(durationStr);
  const finishMinRaw = parseTimeHM(finishStr);

  if (nowMin === null || durMin === null || finishMinRaw === null) {
    errorEl.textContent = "Check that all times are filled in as HH:MM.";
    return;
  }

  // Treat finish earlier than current as "tomorrow"
  let finishMin = finishMinRaw;
  if (finishMin <= nowMin) {
    finishMin += 24 * 60;
  }

  const startMin = finishMin - durMin;
  const delayMin = startMin - nowMin;

  if (delayMin < 0) {
    errorEl.textContent =
      "With that duration you'd have needed to start earlier. Try a later finish time.";
    return;
  }

  if (delayMin > 24 * 60) {
    errorEl.textContent =
      "Delay would be more than 24 hours, which is beyond what W1 usually supports.";
    return;
  }

  const delayW1 = quantizeForW1Delay(delayMin);
  const startExact = nowMin + delayMin;
  const startW1 = nowMin + delayW1;
  const finishW1 = startW1 + durMin;

  const exactDelayText = formatDuration(delayMin);
  const w1DelayText = formatDuration(delayW1);

  resEl.innerHTML = `
    <div class="mb-1">
      <span class="text-gray-600">Exact math:</span>
      <span class="font-semibold text-blue-600">${exactDelayText}</span> delay
      → start at <span class="font-semibold text-blue-600">${formatHM(startExact)}</span>
      → finish at ~${formatHM(startExact + durMin)}.
    </div>
    <div class="mb-1">
      <span class="text-gray-600">Closest W1 delay setting:</span>
      <span class="font-semibold text-blue-600">${w1DelayText}</span>
      → start at <span class="font-semibold text-blue-600">${formatHM(startW1)}</span>
      → finish at ~${formatHM(finishW1)}.
    </div>
    <div class="text-xs text-gray-600">
      (W1 typically uses 30-min steps up to 10 h, then 1-h steps up to 24 h.)
    </div>
  `;
}

function setNow() {
  const now = new Date();
  const h = now.getHours();
  const m = now.getMinutes();
  const t = pad2(h) + ":" + pad2(m);
  document.getElementById("currentTime").value = t;
}

// Only run DOM code in browser environment (not in tests)
if (typeof document !== 'undefined') {
  document.getElementById("calcBtn").addEventListener("click", calculate);
  document.getElementById("nowBtn").addEventListener("click", setNow);

  // Prefill current time & a typical duration
  window.addEventListener("load", () => {
    setNow();
    const durEl = document.getElementById("duration");
    if (!durEl.value) durEl.value = "03:39"; // your example
  });
}
