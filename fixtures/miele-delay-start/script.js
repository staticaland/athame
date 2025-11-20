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
  // "H:MM" or "HH:MM" or "H" (defaults to H:00) -> minutes
  if (!str) return null;

  if (!str.includes(":")) {
    // Just hours, default minutes to 00
    const h = Number(str);
    if (Number.isNaN(h) || h < 0) return null;
    return h * 60;
  }

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

export function quantizeDelay(delayMinutes) {
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

function renderTimeline(nowMin, startMin, durMin, finishMin) {
  const timelineEl = document.getElementById("timeline");
  const timelineCardEl = document.getElementById("timelineCard");

  timelineCardEl.classList.remove("hidden");

  const delayMin = startMin - nowMin;
  const totalMin = finishMin - nowMin;

  // Calculate percentages for visualization
  const waitPercent = (delayMin / totalMin) * 100;
  const runPercent = (durMin / totalMin) * 100;

  timelineEl.innerHTML = `
    <div class="flex items-center gap-2 text-xs sm:text-sm">
      <div class="w-16 sm:w-20 text-gray-600 font-medium">Now</div>
      <div class="flex-1 flex items-center">
        <div class="w-3 h-3 rounded-full bg-gray-400"></div>
        <div class="text-gray-500 ml-2">${formatHM(nowMin)}</div>
      </div>
    </div>

    <div class="flex items-center gap-2">
      <div class="w-16 sm:w-20"></div>
      <div class="flex-1">
        <div class="h-1 bg-gray-200 rounded"></div>
      </div>
    </div>

    <div class="flex items-center gap-2 text-xs sm:text-sm">
      <div class="w-16 sm:w-20 text-gray-600">Wait</div>
      <div class="flex-1 relative">
        <div class="h-8 bg-yellow-100 rounded-lg border-2 border-yellow-300 flex items-center px-2">
          <span class="text-xs font-medium text-yellow-800">${formatDuration(delayMin)}</span>
        </div>
      </div>
    </div>

    <div class="flex items-center gap-2 text-xs sm:text-sm">
      <div class="w-16 sm:w-20 text-gray-600 font-medium">Start</div>
      <div class="flex-1 flex items-center">
        <div class="w-3 h-3 rounded-full bg-blue-500"></div>
        <div class="text-blue-600 font-medium ml-2">${formatHM(startMin)}</div>
      </div>
    </div>

    <div class="flex items-center gap-2">
      <div class="w-16 sm:w-20"></div>
      <div class="flex-1">
        <div class="h-1 bg-gray-200 rounded"></div>
      </div>
    </div>

    <div class="flex items-center gap-2 text-xs sm:text-sm">
      <div class="w-16 sm:w-20 text-gray-600">Running</div>
      <div class="flex-1 relative">
        <div class="h-8 bg-blue-100 rounded-lg border-2 border-blue-300 flex items-center px-2">
          <span class="text-xs font-medium text-blue-800">${formatDuration(durMin)}</span>
        </div>
      </div>
    </div>

    <div class="flex items-center gap-2 text-xs sm:text-sm">
      <div class="w-16 sm:w-20 text-gray-600 font-medium">Finish</div>
      <div class="flex-1 flex items-center">
        <div class="w-3 h-3 rounded-full bg-green-500"></div>
        <div class="text-green-600 font-medium ml-2">${formatHM(finishMin)}</div>
      </div>
    </div>
  `;
}

function calculate() {
  const errorEl = document.getElementById("error");
  const delayCardEl = document.getElementById("delayCard");
  const delayValueEl = document.getElementById("delayValue");
  const delayDetailsEl = document.getElementById("delayDetails");
  const timelineCardEl = document.getElementById("timelineCard");

  errorEl.textContent = "";
  delayCardEl.classList.add("hidden");
  timelineCardEl.classList.add("hidden");

  // Get current time automatically
  const now = new Date();
  const nowMin = now.getHours() * 60 + now.getMinutes();

  const durationStr = document.getElementById("duration").value;
  const finishStr = document.getElementById("finishTime").value;

  const durMin = parseTimeHM(durationStr);
  const finishMinRaw = parseTimeHM(finishStr);

  if (durMin === null || finishMinRaw === null) {
    errorEl.textContent = "Check that programme length and finish time are filled in correctly.";
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
      "Delay would be more than 24 hours, which is beyond what is usually supported.";
    return;
  }

  const quantizedDelay = quantizeDelay(delayMin);
  const startExact = nowMin + delayMin;
  const startTime = nowMin + quantizedDelay;
  const finishTime = startTime + durMin;

  const exactDelayText = formatDuration(delayMin);
  const delayText = formatDuration(quantizedDelay);

  // Show timeline visualization
  renderTimeline(nowMin, startTime, durMin, finishTime);

  // Show prominent delay card
  delayCardEl.classList.remove("hidden");
  delayValueEl.textContent = delayText;
  delayDetailsEl.innerHTML = `Finish at ~<strong>${formatHM(finishTime)}</strong>`;
}

function setSmartFinishTime() {
  const now = new Date();
  const h = now.getHours();

  let finishTime;
  if (h >= 20 || h < 6) {
    // Late evening or early morning → finish at 9am
    finishTime = "09:00";
  } else if (h >= 6 && h < 14) {
    // Morning/early afternoon → finish at 5pm
    finishTime = "17:00";
  } else {
    // Afternoon/evening → finish at 10pm
    finishTime = "22:00";
  }

  document.getElementById("finishTime").value = finishTime;
}

// Only run DOM code in browser environment (not in tests)
if (typeof document !== 'undefined') {
  // Auto-calculate on input changes
  const durationInput = document.getElementById("duration");
  const finishTimeInput = document.getElementById("finishTime");
  const programSelect = document.getElementById("programSelect");

  durationInput.addEventListener("input", calculate);
  finishTimeInput.addEventListener("input", calculate);


  // Add event listener to program dropdown
  programSelect.addEventListener("change", (e) => {
    const duration = e.target.value;
    if (duration) {
      // Convert to HH:MM format for time input
      const durMin = parseDuration(duration);
      if (durMin !== null) {
        durationInput.value = formatHM(durMin);
      }
      calculate(); // Trigger calculation
    }
  });

  // Prefill smart finish time & a typical duration
  window.addEventListener("load", () => {
    setSmartFinishTime();
    const durEl = document.getElementById("duration");
    if (!durEl.value) {
      const defaultDur = parseDuration("1:26");
      durEl.value = formatHM(defaultDur); // Bomull default in HH:MM format
    }
    programSelect.value = "1:26"; // Default to Bomull
    calculate(); // Initial calculation
  });
}
