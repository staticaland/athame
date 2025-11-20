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
  const resultCardEl = document.getElementById("resultCard");
  const delayValueEl = document.getElementById("delayValue");
  const delayDetailsEl = document.getElementById("delayDetails");
  const finishAroundEl = document.getElementById("finishAround");

  errorEl.textContent = "";
  resultCardEl.classList.add("hidden");

  // Get current time automatically
  const now = new Date();
  const nowMin = now.getHours() * 60 + now.getMinutes();

  const durationStr = document.getElementById("duration").value;
  const finishStr = document.getElementById("finishTime").value;

  const durMin = parseTimeHM(durationStr);
  const finishMinRaw = parseTimeHM(finishStr);

  if (durMin === null || finishMinRaw === null) {
    errorEl.textContent = "Check that program length and finish time are filled in correctly.";
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

  // Convert delay minutes to hours for display
  const delayHours = Math.round(quantizedDelay / 60);

  // Show timeline visualization
  renderTimeline(nowMin, startTime, durMin, finishTime);

  // Show result card
  resultCardEl.classList.remove("hidden");
  finishAroundEl.textContent = formatHM(finishTime);
  delayValueEl.textContent = `Delay start: ${delayHours} h`;
  delayDetailsEl.textContent = `Starts at ${formatHM(startTime)}, runs ${formatDuration(durMin)}`;
}

function setSmartFinishTime() {
  const now = new Date();
  const h = now.getHours();

  let finishTime;
  if (h >= 20 || h < 6) {
    // Late evening or early morning → finish at 9:30am
    finishTime = "09:30";
  } else if (h >= 6 && h < 14) {
    // Morning/early afternoon → finish at 5pm
    finishTime = "17:00";
  } else {
    // Afternoon/evening → finish at 10pm
    finishTime = "22:00";
  }

  document.getElementById("finishTime").value = finishTime;
}

const quotes = [
  { text: "The first electric washing machine was invented in 1908 by Alva J. Fisher. Before that, laundry day could take an entire day of hard labor.", author: "Historical fact" },
  { text: "Ancient Romans used urine as a cleaning agent for laundry. It contains ammonia, which is still used in modern cleaning products.", author: "Historical fact" },
  { text: "In medieval Europe, laundry was done only 3-4 times a year. People owned very few clothes and wore them until they were visibly dirty.", author: "Historical fact" },
  { text: "The delay start function was first introduced in the 1980s, revolutionizing how people could schedule their laundry around electricity rates.", author: "Technology history" },
  { text: "Miele was founded in 1899 by Carl Miele and Reinhard Zinkann in Herzebrock, Germany. Their motto 'Immer besser' (Forever better) still guides the company.", author: "Company history" },
  { text: "Before washing machines, a 'washing day' typically involved boiling water in large copper tubs, scrubbing with lye soap, and wringing by hand.", author: "Historical fact" },
  { text: "The average washing machine uses about 50 liters of water per cycle. Modern eco-programs can reduce this by up to 40%.", author: "Environmental fact" },
  { text: "In 18th century England, wealthy families sent their laundry to professional laundresses who worked in communal wash houses.", author: "Historical fact" },
  { text: "The rotating drum design, standard in modern machines, was patented in 1858 by Hamilton Smith.", author: "Patent history" },
  { text: "Blue dye was traditionally added to white laundry to make it appear whiter by counteracting yellow tints. This practice dates back to ancient times.", author: "Historical fact" },
  { text: "The first laundromats appeared in the 1930s in the US, making washing machines accessible to people who couldn't afford their own.", author: "Social history" },
  { text: "Soap was a luxury item until the 19th century. Most people used wood ash, sand, or clay to clean their clothes.", author: "Historical fact" },
  { text: "Monday was traditionally wash day because Sunday leftovers could be eaten while doing laundry, freeing up the rest of the week for other tasks.", author: "Cultural tradition" },
  { text: "The first fully automatic washing machine was introduced by Bendix in 1937. It could wash, rinse, and spin dry without user intervention.", author: "Technology history" },
  { text: "Energy-efficient washing machines can save an average household over 7,000 liters of water per year compared to older models.", author: "Environmental fact" }
];

function displayRandomQuote() {
  const randomIndex = Math.floor(Math.random() * quotes.length);
  const quote = quotes[randomIndex];
  document.getElementById("quote").textContent = `"${quote.text}"`;
  document.getElementById("quoteAuthor").textContent = `— ${quote.author}`;
}

const cheers = [
  "You're absolutely crushing this whole adulting thing!",
  "You deserve fresh, perfectly timed laundry. And that's exactly what you're getting!",
  "You're the kind of person who plans ahead. That's powerful!",
  "You're making tomorrow-you so incredibly happy right now!",
  "You're not just doing laundry, you're mastering time itself!",
  "You're going to wake up tomorrow feeling like a genius. Because you are one!",
  "You're turning a mundane task into an art form. Respect!",
  "You're basically a time wizard, and your fresh laundry is proof!",
  "You're the hero your future self didn't know they needed!",
  "You're making life easier for yourself. That's self-love right there!",
  "You're going to feel so satisfied when that machine beeps at the perfect time!",
  "You're living in 3024 while everyone else is stuck in manual mode!",
  "You're taking control of your schedule like the boss you are!",
  "You're about to experience the pure joy of perfectly timed clean clothes!",
  "You're doing something kind for yourself today. That matters!"
];

function displayRandomCheer() {
  const randomIndex = Math.floor(Math.random() * cheers.length);
  document.getElementById("cheer").textContent = cheers[randomIndex];
}

// Only run DOM code in browser environment (not in tests)
if (typeof document !== 'undefined') {
  // Tab switching logic
  const tabProgram = document.getElementById("tabProgram");
  const tabTime = document.getElementById("tabTime");
  const programTabContent = document.getElementById("programTabContent");
  const timeTabContent = document.getElementById("timeTabContent");

  function switchToTab(tabName) {
    if (tabName === "program") {
      // Update tab buttons
      tabProgram.classList.add("text-blue-600", "border-b-2", "border-blue-600", "bg-blue-50");
      tabProgram.classList.remove("text-gray-500", "hover:text-gray-700");
      tabTime.classList.remove("text-blue-600", "border-b-2", "border-blue-600", "bg-blue-50");
      tabTime.classList.add("text-gray-500", "hover:text-gray-700");

      // Show/hide content
      programTabContent.classList.remove("hidden");
      timeTabContent.classList.add("hidden");
    } else {
      // Update tab buttons
      tabTime.classList.add("text-blue-600", "border-b-2", "border-blue-600", "bg-blue-50");
      tabTime.classList.remove("text-gray-500", "hover:text-gray-700");
      tabProgram.classList.remove("text-blue-600", "border-b-2", "border-blue-600", "bg-blue-50");
      tabProgram.classList.add("text-gray-500", "hover:text-gray-700");

      // Show/hide content
      timeTabContent.classList.remove("hidden");
      programTabContent.classList.add("hidden");
    }
  }

  tabProgram.addEventListener("click", () => switchToTab("program"));
  tabTime.addEventListener("click", () => switchToTab("time"));

  // Auto-calculate on input changes
  const durationInput = document.getElementById("duration");
  const finishTimeInput = document.getElementById("finishTime");
  const programSelect = document.getElementById("programSelect");

  durationInput.addEventListener("input", calculate);
  finishTimeInput.addEventListener("input", calculate);

  // Quick adjustment buttons
  document.getElementById("add30m").addEventListener("click", () => {
    const currentValue = finishTimeInput.value;
    if (currentValue) {
      const currentMin = parseTimeHM(currentValue);
      if (currentMin !== null) {
        const newMin = currentMin + 30;
        finishTimeInput.value = formatHM(newMin);
        calculate();
      }
    }
  });

  document.getElementById("add1h").addEventListener("click", () => {
    const currentValue = finishTimeInput.value;
    if (currentValue) {
      const currentMin = parseTimeHM(currentValue);
      if (currentMin !== null) {
        const newMin = currentMin + 60;
        finishTimeInput.value = formatHM(newMin);
        calculate();
      }
    }
  });

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
    displayRandomQuote(); // Show random quote
    displayRandomCheer(); // Show random cheer
    calculate(); // Initial calculation
  });
}
