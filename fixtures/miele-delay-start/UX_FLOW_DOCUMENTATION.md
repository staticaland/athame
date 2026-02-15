# UX Flow Documentation with Cognitive Load Annotations

## Purpose

This document maps all user flows through the Miele Delay Start Calculator, annotated with cognitive load markers, interaction costs, and UX observations.

## Legend

### Interaction Types
- `[CLICK]` - Mouse click or tap
- `[TYPE]` - Keyboard input
- `[VIEW]` - Passive information consumption
- `[SCAN]` - Visual search
- `[SELECT]` - Dropdown or option selection

### Cognitive Load Markers
- `🧠 LOW` - Minimal mental effort (recognition, simple choice)
- `🧠 MEDIUM` - Moderate mental effort (recall, calculation, decision)
- `🧠 HIGH` - Significant mental effort (learning, problem-solving)

### Mental Operations
- `💭 RECOGNIZE` - Identify familiar pattern
- `💭 RECALL` - Retrieve information from memory
- `💭 DECIDE` - Make a choice between options
- `💭 LEARN` - Understand new concept/pattern
- `💭 VERIFY` - Check if result matches expectation
- `💭 SHIFT` - Change mental model or context

### Interaction Cost
- `⚡ 0-2` - Trivial (highly efficient)
- `⚡ 3-5` - Low (efficient)
- `⚡ 6-10` - Medium (acceptable)
- `⚡ 11+` - High (consider optimization)

---

## Flow 1: First-Time User with Default Settings

**Scenario:** User opens the app for the first time, happy with defaults

**Total Interaction Cost:** `⚡ 0` clicks
**Total Time:** ~3-5 seconds
**Cognitive Load:** `🧠 LOW`

### Step-by-Step Flow

```
┌─────────────────────────────────────────────────────────────────┐
│ Step 1: Page Load                                               │
├─────────────────────────────────────────────────────────────────┤
│ Action: [VIEW] Page loads with pre-populated values             │
│ System: Auto-fills finish time (smart default based on time)   │
│ System: Auto-selects "ECO 40-60 (3:39)" program                │
│ System: Auto-calculates delay start                             │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 RECOGNIZE - Scan interface layout         │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~2-3s (includes page load + visual scan)                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 2: Read Result                                             │
├─────────────────────────────────────────────────────────────────┤
│ Action: [VIEW] Result card displays automatically               │
│ Display: "Set your machine to Delay start: X h"                │
│ Display: "Starts at HH:MM, runs X h Y min"                     │
│ Display: "Finish at HH:MM"                                     │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 VERIFY - Check if result makes sense     │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~1-2s (read and comprehend)                              │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 3: Optional - View Timeline                                │
├─────────────────────────────────────────────────────────────────┤
│ Action: [CLICK] "Show timeline" expansion                       │
│ Display: Visual timeline with wait/running phases               │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 LEARN - Understand detailed breakdown     │
│ Interaction Cost: ⚡ 1 (optional)                               │
│ Time: ~1s click + 2-3s read                                    │
└─────────────────────────────────────────────────────────────────┘
```

**Key UX Observations:**
- ✅ **Zero-interaction path exists** - Critical for efficiency
- ✅ **Smart defaults reduce decision fatigue** - Finish time based on current time
- ✅ **Progressive disclosure** - Timeline hidden by default
- ✅ **Immediate feedback** - No waiting for calculation

---

## Flow 2: Adjust Finish Time (Preset Program)

**Scenario:** User wants laundry done at a different time

**Total Interaction Cost:** `⚡ 2-4` clicks
**Total Time:** ~4-8 seconds
**Cognitive Load:** `🧠 LOW`

### Step-by-Step Flow

```
┌─────────────────────────────────────────────────────────────────┐
│ Step 1: Decide to Change Finish Time                            │
├─────────────────────────────────────────────────────────────────┤
│ Action: [SCAN] Locate finish time field                         │
│ Label: "Have your laundry done around"                         │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 DECIDE - What time do I want?            │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~1s (visual search + decision)                           │
│                                                                  │
│ Information Scent: STRONG - Clear label, prominent placement    │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 2A: Manual Time Entry                                      │
├─────────────────────────────────────────────────────────────────┤
│ Action: [CLICK] Finish time input field                         │
│ Action: [TYPE] Enter desired time (HH:MM format)               │
│ System: Auto-calculates new delay on input                     │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 RECALL - What time do I need?            │
│ Interaction Cost: ⚡ 2 (1 click + 1 type action)               │
│ Time: ~2-3s                                                     │
└─────────────────────────────────────────────────────────────────┘

                    OR

┌─────────────────────────────────────────────────────────────────┐
│ Step 2B: Quick Adjust Buttons                                   │
├─────────────────────────────────────────────────────────────────┤
│ Action: [CLICK] "+30m" or "+1h" button (1-3 times)            │
│ System: Increments finish time automatically                   │
│ System: Auto-calculates new delay on each click                │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 RECOGNIZE - Simple increment pattern     │
│ Interaction Cost: ⚡ 1-3 (one click per adjustment)            │
│ Time: ~1-3s total                                              │
│                                                                  │
│ UX Note: Faster for small adjustments, slower for large ones   │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 3: Result Updates Automatically                            │
├─────────────────────────────────────────────────────────────────┤
│ Action: [VIEW] New result displays instantly                    │
│ System: No explicit "Calculate" button needed                  │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 VERIFY - Check new delay setting         │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~1s                                                       │
└─────────────────────────────────────────────────────────────────┘
```

**Key UX Observations:**
- ✅ **Dual input methods** - Manual and button-based adjustments
- ✅ **Real-time calculation** - Removes explicit calculation step
- ✅ **Low cognitive load** - Simple time selection task
- ⚠️ **Button efficiency varies** - +30m/+1h faster for small changes, slower for large

---

## Flow 3: Change Program (Preset Selection)

**Scenario:** User wants to use a different wash program

**Total Interaction Cost:** `⚡ 2` clicks
**Total Time:** ~3-5 seconds
**Cognitive Load:** `🧠 LOW-MEDIUM`

### Step-by-Step Flow

```
┌─────────────────────────────────────────────────────────────────┐
│ Step 1: Recognize Need to Change Program                        │
├─────────────────────────────────────────────────────────────────┤
│ Action: [SCAN] Locate program selection                         │
│ Location: "Program" tab (already active by default)            │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 RECOGNIZE - Tab is already selected       │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~0.5s                                                     │
│                                                                  │
│ Information Scent: STRONG - "Select program" label visible     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 2: Open Program Dropdown                                   │
├─────────────────────────────────────────────────────────────────┤
│ Action: [CLICK] Program dropdown                                │
│ Display: Shows 4 preset options:                               │
│   - ECO 40-60 (3:39)                                           │
│   - Bomull (2:39)                                              │
│   - Express (0:20)                                             │
│   - Ull (1:09)                                                 │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 RECOGNIZE - Standard dropdown pattern     │
│ Interaction Cost: ⚡ 1                                          │
│ Time: ~1s                                                       │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 3: Select Program                                          │
├─────────────────────────────────────────────────────────────────┤
│ Action: [SELECT] Choose desired program                         │
│ System: Auto-fills duration based on selection                 │
│ System: Auto-calculates new delay                              │
│                                                                  │
│ Cognitive Load: 🧠 LOW-MEDIUM                                   │
│ Mental Operation: 💭 DECIDE - Which program do I need?        │
│ Interaction Cost: ⚡ 1                                          │
│ Time: ~1-2s (includes decision time)                           │
│                                                                  │
│ Decision Complexity: LOW - Only 4 choices (Hick's Law ≈ 0.7s)  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 4: Result Updates Automatically                            │
├─────────────────────────────────────────────────────────────────┤
│ Action: [VIEW] New delay start calculation appears              │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 VERIFY - Check updated result            │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~1s                                                       │
└─────────────────────────────────────────────────────────────────┘
```

**Key UX Observations:**
- ✅ **Default tab state** - No tab switching required for presets
- ✅ **Limited choices** - 4 options prevent decision paralysis
- ✅ **Duration in labels** - Helps users make informed choice
- ✅ **No mode confusion** - Clear you're in "preset" mode

---

## Flow 4: Custom Program Duration (The "Mental Leap" Flow)

**Scenario:** User needs to enter a custom program duration not in presets

**Total Interaction Cost:** `⚡ 8-10` clicks/keystrokes
**Total Time:** ~9-13 seconds
**Cognitive Load:** `🧠 MEDIUM`

### Step-by-Step Flow

```
┌─────────────────────────────────────────────────────────────────┐
│ Step 1: Recognize Need for Custom Duration                      │
├─────────────────────────────────────────────────────────────────┤
│ Action: [SCAN] Look at preset programs                          │
│ Realization: "My program isn't in the list"                    │
│                                                                  │
│ Cognitive Load: 🧠 MEDIUM                                       │
│ Mental Operation: 💭 DECIDE - None of these match my need     │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~2-3s (scan all 4 options + realize none match)         │
│                                                                  │
│ UX Critical Point: User must self-discover custom option       │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 2: Discover "Custom time" Tab                              │
├─────────────────────────────────────────────────────────────────┤
│ Action: [SCAN] Look for alternative input method                │
│ Discovery: Notice "Custom time" tab next to "Program"          │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 RECOGNIZE - Tab indicates alternative     │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~0.5-1s                                                   │
│                                                                  │
│ Information Scent: STRONG - "Custom time" is clear label       │
│ Discoverability: GOOD - Tab is visible, labeled clearly        │
│                                                                  │
│ 🎯 KEY INTERACTION POINT - The "Mental Leap"                    │
│    User must recognize tab switching as solution               │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 3: Switch to Custom Time Mode                              │
├─────────────────────────────────────────────────────────────────┤
│ Action: [CLICK] "Custom time" tab                               │
│ System: Tab visual state changes (blue underline, highlight)   │
│ System: Content area switches from dropdown to time input      │
│                                                                  │
│ Cognitive Load: 🧠 MEDIUM                                       │
│ Mental Operation: 💭 SHIFT - Change from preset to custom mode │
│ Interaction Cost: ⚡ 1                                          │
│ Time: ~1s (click) + ~0.5s (perceive UI change)                │
│                                                                  │
│ 💡 MENTAL MODEL SHIFT OCCURS HERE                               │
│    User transitions from "select" to "input" paradigm          │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 4: Understand New Input Method                             │
├─────────────────────────────────────────────────────────────────┤
│ Action: [VIEW] Observe new interface                            │
│ Display: "Program length" label with time input (HH:MM)        │
│ Note: Previous dropdown is now hidden                          │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 RECOGNIZE - Standard time input field    │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~0.5s                                                     │
│                                                                  │
│ Learnability: HIGH - Time input is familiar pattern            │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 5: Recall Program Duration                                 │
├─────────────────────────────────────────────────────────────────┤
│ Action: [RECALL] User must remember/know their program length   │
│ Potential: User may need to check washing machine display      │
│                                                                  │
│ Cognitive Load: 🧠 MEDIUM-HIGH                                  │
│ Mental Operation: 💭 RECALL - Retrieve duration from memory   │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~1-5s (depends on memory/need to check machine)         │
│                                                                  │
│ ⚠️ POTENTIAL FRICTION POINT                                     │
│    User may not know exact duration off-hand                   │
│    This is intrinsic task complexity, not design flaw          │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 6: Enter Custom Duration                                   │
├─────────────────────────────────────────────────────────────────┤
│ Action: [CLICK] Duration time input field                       │
│ Action: [TYPE] Enter duration in HH:MM format                  │
│ System: Auto-calculates delay on input                         │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 INPUT - Transcribe known value           │
│ Interaction Cost: ⚡ 5 (1 click + ~4 keystrokes for HH:MM)     │
│ Time: ~2-3s                                                     │
│                                                                  │
│ Input Format: Standard HH:MM (familiar to users)               │
│ Error Prevention: Browser validates time format                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 7: Result Updates Automatically                            │
├─────────────────────────────────────────────────────────────────┤
│ Action: [VIEW] Delay calculation appears with custom duration   │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 VERIFY - Check result makes sense        │
│ Interaction Cost: ⚡ 0                                          │
│ Time: ~1s                                                       │
└─────────────────────────────────────────────────────────────────┘
```

**Key UX Observations:**

**Strengths:**
- ✅ **Clear tab labeling** - "Custom time" is unambiguous
- ✅ **Standard UI pattern** - Tab switching is familiar
- ✅ **Good visual feedback** - Tab state change is obvious
- ✅ **No dead ends** - User can easily switch back to preset mode

**The "Mental Leap" Breakdown:**
1. **Discovery phase** (~2-3s): Scanning presets, realizing none match
2. **Recognition phase** (~0.5-1s): Noticing the "Custom time" tab
3. **Mode shift** (~1.5s): Clicking tab and perceiving UI change
4. **Adaptation** (~0.5s): Understanding new input method

**Total "Mental Leap" cost: ~4.5-6s and +1 mental model shift**

**Friction Points:**
- ⚠️ **Recall requirement** - User must know exact program duration
- ⚠️ **Self-discovery needed** - No explicit prompt to use custom tab
- ✅ **Mitigated by design** - Clear labeling reduces discovery time

**Is this acceptable?**
- ✅ **Yes** - Custom duration is likely a minority use case
- ✅ **Yes** - 4-6s cognitive overhead is low in absolute terms
- ✅ **Yes** - Tab pattern prevents interface clutter for majority users
- ✅ **Yes** - Information scent is strong enough for self-discovery

---

## Flow 5: Power User (Multiple Adjustments)

**Scenario:** User tweaks both program and finish time multiple times

**Total Interaction Cost:** `⚡ 6-12` clicks (variable)
**Total Time:** ~10-20 seconds
**Cognitive Load:** `🧠 LOW-MEDIUM`

### Step-by-Step Flow

```
┌─────────────────────────────────────────────────────────────────┐
│ Step 1: Experiment with Different Programs                      │
├─────────────────────────────────────────────────────────────────┤
│ Action: [SELECT] Try "Express (0:20)"                           │
│ System: Updates delay instantly                                │
│ Action: [SELECT] Try "ECO 40-60 (3:39)"                        │
│ System: Updates delay instantly                                │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 DECIDE - Which program is best?          │
│ Interaction Cost: ⚡ 4 (2 clicks × 2 programs)                 │
│ Time: ~4-6s                                                     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 2: Fine-Tune Finish Time                                   │
├─────────────────────────────────────────────────────────────────┤
│ Action: [CLICK] "+30m" button twice                            │
│ System: Updates delay each time                                │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 OPTIMIZE - Find ideal finish time        │
│ Interaction Cost: ⚡ 2                                          │
│ Time: ~2-3s                                                     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Step 3: Expand Timeline for Verification                        │
├─────────────────────────────────────────────────────────────────┤
│ Action: [CLICK] "Show timeline"                                │
│ Display: Detailed timeline visualization                       │
│                                                                  │
│ Cognitive Load: 🧠 LOW                                          │
│ Mental Operation: 💭 VERIFY - Confirm schedule works          │
│ Interaction Cost: ⚡ 1                                          │
│ Time: ~1s + 2-3s reading                                       │
└─────────────────────────────────────────────────────────────────┘
```

**Key UX Observations:**
- ✅ **Exploration is frictionless** - Instant updates encourage experimentation
- ✅ **No commitment cost** - Changes apply immediately, easy to revert
- ✅ **Progressive disclosure supports learning** - Timeline adds depth for interested users
- ✅ **State is visible** - Current selections always clear

---

## Comparative Flow Analysis

### Interaction Cost Comparison Table

| Flow | Clicks | Time | Cognitive Load | Mental Shifts | Use Case Frequency |
|------|--------|------|----------------|---------------|-------------------|
| **Default (no changes)** | 0 | 3-5s | 🧠 LOW (2/10) | 0 | High (40-50%) |
| **Adjust finish time only** | 2-4 | 4-8s | 🧠 LOW (2/10) | 0 | High (30-40%) |
| **Change preset program** | 2 | 3-5s | 🧠 LOW-MED (3/10) | 0 | Medium (20-30%) |
| **Custom duration** | 8-10 | 9-13s | 🧠 MEDIUM (4/10) | 1 | Low (10-20%) |
| **Power user tweaking** | 6-12 | 10-20s | 🧠 LOW-MED (3/10) | 0 | Low (5-10%) |

### Mental Model Requirements by Flow

```
Flow 1 (Default):
  Mental Models Required: [Time, Basic calculation]
  Complexity: ★☆☆☆☆

Flow 2 (Adjust finish):
  Mental Models Required: [Time, Basic calculation, Time input]
  Complexity: ★☆☆☆☆

Flow 3 (Change program):
  Mental Models Required: [Time, Wash programs, Dropdown selection]
  Complexity: ★★☆☆☆

Flow 4 (Custom duration):
  Mental Models Required: [Time, Tabs, Mode switching, Duration recall]
  Complexity: ★★★☆☆  ← The "mental leap"

Flow 5 (Power user):
  Mental Models Required: [All of the above + experimentation mindset]
  Complexity: ★★☆☆☆
```

---

## Cognitive Load Distribution

### Load Types by Flow

**Flow 1 (Default):**
- Intrinsic Load: ★☆☆☆☆ (Task is simple)
- Extraneous Load: ★☆☆☆☆ (Interface adds no complexity)
- Germane Load: ★☆☆☆☆ (Minimal learning needed)
- **Total: 2/10**

**Flow 2 (Adjust finish):**
- Intrinsic Load: ★☆☆☆☆ (Task is simple)
- Extraneous Load: ★☆☆☆☆ (Standard time input)
- Germane Load: ★☆☆☆☆ (Quick adjustment buttons are intuitive)
- **Total: 2/10**

**Flow 3 (Change program):**
- Intrinsic Load: ★★☆☆☆ (Requires knowing wash programs)
- Extraneous Load: ★☆☆☆☆ (Standard dropdown)
- Germane Load: ★☆☆☆☆ (Minimal)
- **Total: 3/10**

**Flow 4 (Custom duration):**
- Intrinsic Load: ★★☆☆☆ (Must know/recall duration)
- Extraneous Load: ★★☆☆☆ (Tab switching adds cognitive step)
- Germane Load: ★★☆☆☆ (Learning tab pattern + mode concept)
- **Total: 4/10**

**Flow 5 (Power user):**
- Intrinsic Load: ★★☆☆☆ (Optimization decisions)
- Extraneous Load: ★☆☆☆☆ (Exploration is frictionless)
- Germane Load: ★★☆☆☆ (Learning timeline feature)
- **Total: 3/10**

---

## Error Prevention & Recovery

### Potential Error Scenarios

#### Error 1: Finish Time in the Past
```
Scenario: User enters finish time earlier than current time
System Response: Assumes "tomorrow" and adjusts calculation
Cognitive Load: 🧠 LOW
Recovery Cost: ⚡ 0 (automatic)
UX Quality: ✅ Excellent - intelligent default behavior
```

#### Error 2: Impossible Schedule
```
Scenario: User wants finish time too soon for program duration
System Response: Shows error "you'd have needed to start earlier"
Cognitive Load: 🧠 LOW
Recovery Cost: ⚡ 2-4 (adjust finish time)
UX Quality: ✅ Good - clear error message with solution hint
```

#### Error 3: Wrong Tab
```
Scenario: User clicks "Custom time" by accident
System Response: Shows time input instead of dropdown
Cognitive Load: 🧠 LOW
Recovery Cost: ⚡ 1 (click back to "Program" tab)
UX Quality: ✅ Excellent - easy reversal, no data loss
```

#### Error 4: Invalid Time Format
```
Scenario: User enters malformed time
System Response: Browser validates input (native time picker)
Cognitive Load: 🧠 LOW
Recovery Cost: ⚡ 1-2 (correct input)
UX Quality: ✅ Good - browser-level validation
```

**Error Prevention Score: 9/10**

---

## Accessibility Flow Annotations

### Keyboard Navigation Flow

```
Tab Order:
1. Finish time input [FOCUS]
2. "+30m" button [FOCUS]
3. "+1h" button [FOCUS]
4. "Program" tab button [FOCUS] [ARROW KEYS for tab switching]
5. "Custom time" tab button [FOCUS] [ARROW KEYS for tab switching]
6. Program dropdown [FOCUS] [ARROW KEYS for selection]
   OR
   Custom duration input [FOCUS]
7. "Show timeline" details element [FOCUS] [ENTER to expand]

Total Tab Stops: 7 (efficient)
Keyboard Accessibility: ✅ Full navigation possible
```

### Screen Reader Flow (Custom Duration Path)

```
1. [ANNOUNCE] "Have your laundry done around, time input"
2. [ANNOUNCE] "Plus 30 minutes, button"
3. [ANNOUNCE] "Plus 1 hour, button"
4. [ANNOUNCE] "Program, tab, selected"
5. [USER ACTION] Arrow right to next tab
6. [ANNOUNCE] "Custom time, tab"
7. [USER ACTION] Enter to activate tab
8. [SHOULD ANNOUNCE] "Program length, time input"
   ⚠️ Needs testing: Tab change should announce content update
9. [USER INPUT] Enter custom duration
10. [SHOULD ANNOUNCE] "Result card, Set your machine to Delay start X hours"
    ⚠️ Needs aria-live for dynamic updates

Accessibility Gaps:
- ⚠️ Tab switching may not announce content change to screen readers
- ⚠️ Auto-calculation results need aria-live region
- ⚠️ Timeline expansion needs aria-expanded state
```

---

## Journey Map: First-Time User (Custom Duration)

### Emotional & Cognitive Journey

```
Phase 1: Arrival (0-3s)
├─ Emotion: Neutral → Slightly curious
├─ Thought: "Let me figure out when to start this machine"
├─ Action: Scan interface
└─ Cognitive Load: 🧠 LOW

Phase 2: Initial Understanding (3-6s)
├─ Emotion: Curious → Slightly confused
├─ Thought: "These preset programs don't match mine"
├─ Action: Look for alternatives
└─ Cognitive Load: 🧠 MEDIUM
    └─ 💡 DISCOVERY MOMENT: Notices "Custom time" tab

Phase 3: Mode Switch (6-8s)
├─ Emotion: Slightly confused → Confident
├─ Thought: "Ah, I can enter my own time here"
├─ Action: Click "Custom time" tab
└─ Cognitive Load: 🧠 MEDIUM
    └─ 💭 MENTAL SHIFT: From selection to input paradigm

Phase 4: Input (8-11s)
├─ Emotion: Confident
├─ Thought: "I know my program is 2 hours 15 minutes"
├─ Action: Enter "02:15"
└─ Cognitive Load: 🧠 LOW

Phase 5: Validation (11-13s)
├─ Emotion: Confident → Satisfied
├─ Thought: "That looks right!"
├─ Action: Read result
└─ Cognitive Load: 🧠 LOW

Total Journey: 13 seconds, 1 moment of confusion (quickly resolved)
Outcome: ✅ Success
Satisfaction: 😊 High
```

---

## Design Pattern Assessment

### Tab Switching Pattern Evaluation

**Context:** Custom duration requires switching from "Program" to "Custom time" tab

**Pros:**
- ✅ Familiar UI pattern (high learnability)
- ✅ Clear separation of concerns (preset vs. custom)
- ✅ Prevents interface clutter (hides complexity)
- ✅ Reversible action (easy to switch back)
- ✅ Strong information scent (clear labels)

**Cons:**
- ⚠️ Adds one interaction (+1 click)
- ⚠️ Requires discovery (not immediately obvious)
- ⚠️ Introduces mental model shift (mode change)

**Alternatives Considered:**

1. **Always show both preset and custom inputs**
   - Pro: No tab switching needed
   - Con: Interface clutter, unclear which to use
   - Verdict: ❌ Worse UX for majority users

2. **"Other" option in dropdown**
   - Pro: No tab needed, stays in selection paradigm
   - Con: Hidden until dropdown opened, awkward "other" selection step
   - Verdict: ≈ Similar complexity, less clear

3. **Smart detection (show custom when no preset matches)**
   - Pro: Adaptive interface
   - Con: Unpredictable behavior, complex implementation
   - Verdict: ❌ Over-engineering

**Final Assessment:** ✅ Tab pattern is optimal for this use case

---

## Summary: The "Mental Leap" Quantified

### What is the "mental leap" in this interface?

The mental leap occurs in **Flow 4: Custom Program Duration**

**Quantified Cost:**
- **Interaction cost:** +1 click (tab switch)
- **Time cost:** +4-6 seconds (discovery + mode shift + adaptation)
- **Cognitive cost:** +2 points (from 2/10 to 4/10)
- **Mental model shifts:** +1 (preset selection → manual input)

**Breakdown:**
1. **Discovery phase** (2-3s): Realizing presets don't match need
2. **Recognition phase** (0.5-1s): Noticing "Custom time" tab
3. **Execution phase** (1s): Clicking tab
4. **Adaptation phase** (0.5-1s): Understanding new input method
5. **Mode shift** (ongoing): Switching from "select" to "input" mental model

**Is this acceptable?**

✅ **YES**, because:
- Absolute cost is low (~6s and 1 click)
- Information scent is strong (tab is clearly labeled)
- Pattern is familiar (tabs are well-understood)
- Use case is infrequent (~10-20% of users)
- Alternative designs would worsen UX for majority users
- No better solution exists without trade-offs

**Design Principle Validated:**
> "Optimize for the common case, make the uncommon case possible."

The interface optimizes for preset selection (common) while making custom input accessible (uncommon) without cluttering the interface.

---

## Recommendations

### Keep As-Is (High Confidence)
1. ✅ Tab-based mode switching for custom duration
2. ✅ Auto-calculation on input (no calculate button)
3. ✅ Smart default finish time
4. ✅ Progressive disclosure of timeline
5. ✅ Quick adjustment buttons (+30m/+1h)

### Consider for Future (Low Priority)
1. 💡 Add subtle hint on Program tab: "Need a different duration? →"
   - Would reduce discovery time by ~1-2s
   - Trade-off: Adds visual noise

2. 💡 Add aria-live regions for screen reader support
   - Critical for accessibility
   - No UX trade-offs

3. 💡 Add keyboard shortcut to switch tabs (Alt+1/Alt+2)
   - Power user efficiency gain
   - Zero cost to other users

### Do Not Change
1. ❌ Don't merge preset and custom into single interface
2. ❌ Don't add explicit "Calculate" button
3. ❌ Don't hide the custom time option further
4. ❌ Don't auto-switch to custom tab when typing in finish time

---

**Document Version:** 1.0
**Last Updated:** 2025-11-23
**Analysis Basis:** Miele Delay Start Calculator (fixtures/miele-delay-start)
