# UX Interaction Cost Analysis: Miele Delay Start Calculator

## Executive Summary

This document provides a quantitative analysis of the Miele Delay Start Calculator's user experience using established UX evaluation methodologies. The analysis identifies interaction costs, cognitive loads, and optimization opportunities.

## Methodology

This analysis employs multiple UX evaluation frameworks:

- **Interaction Cost Analysis** - Counting discrete user actions
- **Keystroke-Level Model (KLM)** - Time-based task prediction
- **Cognitive Load Assessment** - Mental effort required
- **Information Scent Analysis** - Navigation clarity
- **Hick's Law** - Decision-making complexity

## User Goals

Primary user goal: **Configure washing machine to finish at a desired time**

### Sub-goals:
1. Select program duration (preset or custom)
2. Specify desired finish time
3. Determine delay start setting

## Task Analysis

### Task 1: Quick Calculation with Preset Program (Default Path)

**User Actions:**
1. (Optional) Adjust finish time if needed
2. (Optional) Select different preset program
3. View result

**Interaction Cost Breakdown:**

| Step | Action Type | Cost | Time (KLM) | Notes |
|------|-------------|------|------------|-------|
| Load page | System | 0 | 2.0s | Auto-populated with defaults |
| Understand interface | Visual search | 0 | 1.0s | Scan layout |
| (Optional) Click finish time field | Point + Click | 2 | 1.3s | P=1.1s, K=0.2s |
| (Optional) Adjust time | Keystroke | 1-4 | 0.2-0.8s | Depends on adjustments |
| (Optional) Use +30m/+1h buttons | Point + Click | 2 | 1.3s | Alternative to manual entry |
| (Optional) Select program dropdown | Point + Click | 2 | 1.3s | P=1.1s, K=0.2s |
| (Optional) Choose program | Point + Click | 2 | 1.3s | P=1.1s, K=0.2s |
| Result auto-displays | System | 0 | 0.2s | Automatic calculation |

**Total for minimal path (using defaults):**
- **Interaction Cost: 0 clicks**
- **Time: ~3.2s** (page load + comprehension + auto-result)
- **Mental operations: 1** (understand result)

**Total for customized path (change finish time + program):**
- **Interaction Cost: 4-6 clicks**
- **Time: ~8-10s**
- **Mental operations: 3** (decide finish time, choose program, verify result)

### Task 2: Calculation with Custom Program Duration

**User Actions:**
1. Click "Custom time" tab
2. Enter custom duration
3. (Optional) Adjust finish time
4. View result

**Interaction Cost Breakdown:**

| Step | Action Type | Cost | Time (KLM) | Notes |
|------|-------------|------|------------|-------|
| Load page | System | 0 | 2.0s | Auto-populated with defaults |
| Recognize need for custom time | Mental | 0 | 1.35s | M=1.35s (mental preparation) |
| Locate "Custom time" tab | Visual search | 0 | 0.5s | Good information scent |
| Click "Custom time" tab | Point + Click | 2 | 1.3s | P=1.1s, K=0.2s |
| Observe UI change | Perception | 0 | 0.3s | Visual feedback |
| Click duration field | Point + Click | 2 | 1.3s | P=1.1s, K=0.2s |
| Enter custom duration | Keystroke | 4 | 0.8s | HH:MM format |
| (Optional) Adjust finish time | Point + Click + Keys | 2-6 | 1.5-2.1s | If needed |
| Result auto-displays | System | 0 | 0.2s | Automatic calculation |

**Total:**
- **Interaction Cost: 8-14 clicks/keystrokes**
- **Time: ~9-11s** (without finish time adjustment) or **~10.5-13s** (with adjustment)
- **Mental operations: 2** (recognize custom need, verify result)
- **Mental model shifts: 1** (from preset to custom paradigm)

## Cognitive Load Analysis

### Task 1: Preset Program (Default Flow)

**Intrinsic Cognitive Load: LOW**
- Task is straightforward: pick finish time, pick program
- Mental model matches real-world task
- Minimal domain knowledge required

**Extraneous Cognitive Load: VERY LOW**
- Defaults pre-populate sensible values
- Auto-calculation removes manual steps
- Clear visual hierarchy
- Immediate feedback

**Germane Cognitive Load: MINIMAL**
- Users learn the pattern quickly
- Reinforces understanding of delay start concept
- Timeline visualization aids learning

**Total Cognitive Load: LOW (2/10)**

### Task 2: Custom Duration

**Intrinsic Cognitive Load: LOW-MEDIUM**
- Requires knowing exact program duration
- Same core task complexity

**Extraneous Cognitive Load: LOW-MEDIUM**
- **+1 mental leap:** Tab switching introduces mode change
- Users must discover the "Custom time" tab
- Tab affordance is clear (good information scent)
- Time input format is standard

**Germane Cognitive Load: MINIMAL**
- One-time learning of tab location

**Total Cognitive Load: LOW-MEDIUM (4/10)**

**Key Finding:** The tab switch adds approximately **+2 points** of cognitive load due to:
1. Recognizing the need for a different mode
2. Discovering the tab mechanism
3. Understanding the mode change

## Information Scent Analysis

### Navigation Elements

| Element | Information Scent | Findability | Notes |
|---------|------------------|-------------|-------|
| "Program" tab | Strong | Excellent | Clear, visible, default state |
| "Custom time" tab | Strong | Excellent | Clear alternative label |
| Program dropdown | Strong | Excellent | Standard UI pattern |
| Finish time input | Strong | Excellent | Prominent placement, clear label |
| +30m/+1h buttons | Medium | Good | Secondary affordance, discoverable |
| Timeline details | Medium | Good | Hidden by default (progressive disclosure) |

**Overall Information Scent: STRONG**

The interface provides clear signifiers for all major actions. The tab distinction between "Program" and "Custom time" creates a clear mental model.

## Decision Complexity (Hick's Law)

### Initial Decision Points

**Primary decision: Which mode to use?**
- **Choices: 2** (Program dropdown vs. Custom time)
- **Hick's Law prediction:** T = 0.3 × log₂(2 + 1) ≈ 0.48s
- **Actual time likely higher** due to semantic processing of options

**Secondary decision (Program mode): Which program?**
- **Choices: 4** preset programs
- **Hick's Law prediction:** T = 0.3 × log₂(4 + 1) ≈ 0.70s

**Total decision time: ~1.2s** (just for decision-making)

### Optimization Note

The current design minimizes decision paralysis by:
1. Pre-selecting a sensible default program
2. Clearly separating the two modes
3. Limiting preset choices to 4 common programs

## Comparative Analysis

### Preset vs. Custom Duration

| Metric | Preset Program | Custom Duration | Difference |
|--------|----------------|-----------------|------------|
| **Clicks** | 0-4 | 8-14 | +8-10 clicks |
| **Time** | 3-10s | 9-13s | +3-6s |
| **Mental operations** | 1-3 | 2 | Similar |
| **Mental model shifts** | 0 | 1 | +1 shift |
| **Cognitive load** | 2/10 | 4/10 | +2 points |
| **Error potential** | Very Low | Low | Slightly higher |

**Key Finding:** Custom duration path has approximately:
- **200-300% more interactions**
- **50-100% more time**
- **50% more cognitive load**

This is acceptable because:
1. Custom duration is a less common use case
2. The interaction cost is still low in absolute terms
3. Tab switching is a well-understood pattern

## Interface Strengths

### 1. **Excellent Defaults**
- Auto-populates with sensible values
- Immediate calculation on page load
- Reduces interaction cost to near-zero for happy path

### 2. **Automatic Calculation**
- No "Calculate" button needed
- Real-time updates reduce mental effort
- Eliminates one click per interaction

### 3. **Progressive Disclosure**
- Timeline hidden by default
- Reduces visual complexity
- Advanced users can expand for details

### 4. **Quick Adjustment Affordances**
- +30m/+1h buttons reduce typing
- Directly manipulates the time value
- Saves ~2s per adjustment vs. manual entry

### 5. **Clear Mode Separation**
- Tabs create distinct mental models
- Prevents interface clutter
- Good information architecture

## Optimization Opportunities

### 1. **Reduce Custom Duration Discovery Time** (Priority: Low)

**Current state:** Users must recognize tab, click to switch
**Potential improvement:** Add subtle hint on Program tab like "Need a custom duration?"

**Impact:**
- Reduces discovery time by ~1-2s for first-time users
- Lowers mental model shift cost
- Trade-off: Adds visual clutter

**Recommendation:** Not worth implementing - current design is sufficiently clear

### 2. **Smart Program Suggestions** (Priority: Low)

**Potential improvement:** If user types a custom duration that matches a preset, suggest the preset

**Impact:**
- Reduces custom path usage by ~10-20%
- Saves 4-6 clicks for applicable cases

**Recommendation:** Consider for future iteration if usage data shows high custom duration usage

### 3. **Duration Input Optimization** (Priority: Very Low)

**Current state:** Standard time picker (HH:MM)
**Potential improvement:** Number stepper or duration-specific picker

**Impact:**
- Might save 1-2s on mobile devices
- Could improve error prevention

**Recommendation:** Current implementation is standard and works well

## Accessibility Considerations

### Keyboard Navigation
- Tab interface is keyboard-accessible
- Time inputs use native controls (good accessibility)
- No keyboard traps identified

### Screen Reader Experience
- Tab switching should announce mode change
- Auto-calculation should announce result updates
- Timeline expansion needs ARIA labels

**Recommendation:** Add `aria-live` regions for dynamic content updates

## Summary Metrics

### Overall Interface Performance

| Metric | Rating | Score |
|--------|--------|-------|
| **Interaction Efficiency** | Excellent | 9/10 |
| **Cognitive Load** | Very Low | 9/10 |
| **Learnability** | Excellent | 9/10 |
| **Error Prevention** | Excellent | 9/10 |
| **Information Scent** | Strong | 9/10 |
| **Accessibility** | Good | 7/10 |

**Overall UX Score: 8.7/10**

### Key Findings

1. **Default path is optimized:** Near-zero interaction cost with sensible defaults
2. **Custom path is acceptable:** While requiring more interactions, the cost is still low
3. **Tab pattern is appropriate:** Clear separation without excessive complexity
4. **Auto-calculation is a major win:** Eliminates explicit calculation step
5. **Minor accessibility improvements needed:** Add ARIA labels for dynamic content

## Conclusion

The Miele Delay Start Calculator demonstrates excellent UX design with:

- **Minimal interaction costs** for the primary use case
- **Low cognitive load** across all flows
- **Clear information architecture** with strong scent
- **Appropriate use of progressive disclosure**

The "mental leap" required for custom duration (tab switching) is:
- **Quantifiable:** +1 click, +1-2s, +2 cognitive load points
- **Justified:** Keeps the default interface clean and focused
- **Well-executed:** Clear labeling and standard pattern

**Recommendation:** The current design is well-optimized. No significant changes needed.

---

**Analysis Methodology Reference:**

- Card, S. K., Moran, T. P., & Newell, A. (1983). *The Psychology of Human-Computer Interaction*. Lawrence Erlbaum Associates.
- Nielsen, J. (1993). *Usability Engineering*. Academic Press.
- Hick, W. E. (1952). "On the rate of gain of information". *Quarterly Journal of Experimental Psychology*, 4(1), 11-26.
- Sweller, J. (1988). "Cognitive load during problem solving: Effects on learning". *Cognitive Science*, 12(2), 257-285.
