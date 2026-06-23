# media-single-modal: Smoothness, Memory & Zoom Refactor

## Context

`app/src/components/media/media-single-modal.vue` is the full-screen viewer for images, video and audio. Three bugs and one architectural question were reported:

1. **Audio → Image swap**: after swiping from an audio item to an image, the image is not initially centered/scaled — it sits at native 100% width until the user pinches.
2. **Inconsistent swipe**: swipe-to-navigate sometimes only fires on the second attempt.
3. **High memory footprint** of the modal.
4. **`zoompinch` package** (0.0.48, pre-1.0, only used here): should we replace it with our own implementation?

The goal is to make the modal feel snappy and predictable on iOS Safari/Chromium, and to reduce its memory footprint on long browsing sessions.

---

## Diagnosis

### Issue 1 — Audio → Image: image not centered/sized

The modal renders **two** image elements when `isZoomMode` is true:

- A hidden `<img ref="imgRef">` ([media-single-modal.vue:52-60](app/src/components/media/media-single-modal.vue#L52-L60)) whose only job is to expose `naturalWidth/Height` to the parent via `onMediaLoaded` ([media-single-modal.vue:227-231](app/src/components/media/media-single-modal.vue#L227-L231)).
- `<media-zoom>` ([media-single-modal.vue:78-88](app/src/components/media/media-single-modal.vue#L78-L88)) which renders the real image via `<ion-img>` inside `zoompinch`.

`media-zoom.computed initialScale = containerWidth / naturalWidth` ([media-zoom.vue:50-53](app/src/components/media/media-zoom.vue#L50-L53)), and `reset()` is only called once, from `onImgDidLoad` of the inner `ion-img` ([media-zoom.vue:61-64](app/src/components/media/media-zoom.vue#L61-L64)).

**Race**: when the user navigates audio → image, the parent's `imgRef` is freshly mounted (`:key="'img-' + media.id"`) — so `imgNaturalWidth/Height` still hold their default `1920/1080` until that hidden `<img>` fires `@load`. If `ion-img-did-load` inside `media-zoom` fires first (often it does — ion-img has its own internal load pipeline), `reset()` runs with `naturalWidth = 1920`, producing a wrong `initialScale`. When the hidden `<img>` later loads and updates `naturalWidth`, the `initialScale` computed updates — **but `transform.scale` is not re-applied**, so the visual stays wrong. A pinch triggers a new transform calculation and the image snaps to fit.

### Issue 2 — Swipe sometimes only works on 2nd try

`initSwipe()` fires `onPrev/onNext` inside `onMove` and then disables the gesture for **350 ms** via a `reEnableTimeout` ([media-single-modal.vue:352-383](app/src/components/media/media-single-modal.vue#L352-L383)). The slide-out animation itself is 300 ms ([media-single-modal.vue:325-332](app/src/components/media/media-single-modal.vue#L325-L332)).

If the user lifts and re-swipes within ~350 ms, the gesture is still disabled and the new swipe is silently dropped. There's no `onStart`/`onEnd` reset — the lock is purely time-based, which is fragile across animation timing, image-load timing, and user pace.

### Issue 3 — Memory pressure

Concrete sources:

- **Duplicate image**: hidden `<img>` + `<ion-img>` inside `media-zoom` both download/decode the same file. ~2× memory for every image.
- **Blob URL pattern** in `loadOriginal()` ([media-single-modal.vue:248-270](app/src/components/media/media-single-modal.vue#L248-L270)): the full original is read into a JS `Blob`, then `URL.createObjectURL` produces a URL that retains the bytes until revoked. The blob URL is *only* revoked on next navigation. Compared to `<img src=...>` (where the browser frees decoded bitmaps under pressure), this pins ~5–15 MB per image.
- **Preload**: `preloadMedia()` ([media-single-modal.vue:233-237](app/src/components/media/media-single-modal.vue#L233-L237)) creates `new Image()` for prev *and* next thumbnails on every change. Detached from DOM, never tracked — relies on GC.
- **`zoompinch` canvas** keeps its own bitmap layer on top of `<ion-img>`'s already-decoded image.

### Issue 4 — Replace `zoompinch`?

`zoompinch@0.0.48` is pre-1.0 and only consumed by [media-zoom.vue](app/src/components/media/media-zoom.vue). It is the source of one indirect problem (its `touch-action: none` may swallow pointer events the parent swipe gesture wants to see), but **it is not the root cause of issues 1–3**. A homegrown pinch/pan implementation is ~250–400 LOC with non-trivial edge cases (bounds, momentum, pinch focus point, double-tap-to-zoom, wheel zoom, RTL). Recommendation: **defer**. Fix 1–3 inside the current architecture; revisit only if zoompinch keeps causing gesture conflicts after Issue 2 is fixed.

---

## Plan

### Step 1 — Move natural-dimension tracking *into* `media-zoom` (fixes Issue 1, halves image memory)

In [media-zoom.vue](app/src/components/media/media-zoom.vue):

- Capture `naturalWidth/naturalHeight` from the actual loaded `<ion-img>` (event target's `naturalWidth`, or the inner `<img>` via `querySelector`) inside `onImgDidLoad`. Store as local refs.
- Drop the `naturalWidth` / `naturalHeight` props.
- Use a `ResizeObserver` on `.projection-wrapper` to keep an internal `containerWidth` ref up-to-date; drop the `containerWidth` prop.
- Add `watch([naturalWidth, containerWidth], () => reset())` so when either dimension settles, the transform is re-applied. This is the actual fix for Issue 1.

In [media-single-modal.vue](app/src/components/media/media-single-modal.vue):

- Remove the hidden `<img ref="imgRef">` block ([L52-L60](app/src/components/media/media-single-modal.vue#L52-L60)) and the `imgNaturalWidth/Height`, `imgRef`, `onMediaLoaded` refs.
- Remove the `:natural-width`, `:natural-height`, `:container-width` bindings on `<media-zoom>`.
- For `isMediaLoaded` (used to gate the spinner on video/audio), wire it from `<video @loadeddata>`/`<audio @loadeddata>` only — image loading state moves to `media-zoom` (it already shows the image when ready).

### Step 2 — Rewrite the swipe gesture (fixes Issue 2)

In [media-single-modal.vue:initSwipe](app/src/components/media/media-single-modal.vue#L352-L383):

- Replace the time-based `reEnableTimeout` with a per-gesture lock:
  ```ts
  let firedThisGesture = false;
  createGesture({
    onStart() { firedThisGesture = false; },
    onMove(ev) {
      if (firedThisGesture) return;
      if (!canNavigate.value) return;
      if (ev.deltaX > 60 && props.prevMediaId) { firedThisGesture = true; onPrev(); }
      else if (ev.deltaX < -60 && props.nextMediaId) { firedThisGesture = true; onNext(); }
    },
  });
  ```
- No `enable(false)` / `enable(true)` dance, no timeout. Each fresh finger-down resets the lock.
- Keep the 300 ms slide-out animation, but stop awaiting it before emitting `prev/next` — emit immediately so the parent can fetch the next file in parallel with the visual transition.

### Step 3 — Drop the blob-URL fetch pattern (fixes Issue 3, partially)

In [media-single-modal.vue:loadOriginal](app/src/components/media/media-single-modal.vue#L248-L270):

- Remove `fetch` → `blob()` → `URL.createObjectURL` pipeline.
- Set `fileUrl.value = getFileUrl(newMedia.id)` directly — the browser handles credentials (same-origin cookies on `/api/...`) and HTTP caching. Decoded bitmaps can be freed under memory pressure, unlike retained blobs.
- Replace `AbortController` cancellation with the natural cancellation of changing `<img src>` (browser cancels the previous request). Keep `cleanupOriginal()` only for the case where we *do* want to clear `fileUrl` to show the spinner — collapse into setting `fileUrl.value = undefined`.
- Still show the thumbnail first by setting `fileUrl` to the thumbnail URL, then swap to the full file URL after a `nextTick` — the browser will fetch the original while the thumb is visible (or just go straight to full file; profile both).

### Step 4 — Trim preloading (fixes Issue 3, partially)

In [media-single-modal.vue:preloadMedia](app/src/components/media/media-single-modal.vue#L233-L237):

- Preload only the next thumbnail (drop prev), and only when the network is fast (`navigator.connection?.effectiveType !== 'slow-2g'` if available — otherwise unconditional is fine).
- Use `<link rel="prefetch">` injection instead of `new Image()` so the browser handles lifecycle.

### Step 5 — Keep `zoompinch` for now (Issue 4)

Document the decision in a 1-line comment in `media-zoom.vue` referencing this plan. Re-evaluate if Issue 2 reappears post-fix.

---

## Files to be modified

- [app/src/components/media/media-single-modal.vue](app/src/components/media/media-single-modal.vue) — remove hidden img, simplify gesture, drop blob-URL pipeline, trim preload.
- [app/src/components/media/media-zoom.vue](app/src/components/media/media-zoom.vue) — own its dimensions, watch them, re-reset transform.

No new files. No changes to backend, stores, or `library-page.vue` (parent integration stays unchanged).

---

## Verification

Manual (the only way to verify these UX bugs):

1. **Issue 1**: Open an audio item in the modal, swipe to a portrait image, swipe to a landscape image, swipe back. Image should always be centered and fit-to-screen on first paint. Repeat with the order video → image and image → audio → image.
2. **Issue 2**: Rapid alternating swipes left-right-left-right. Every swipe should register. Test with `prevMediaId === null` (first item) and `nextMediaId === null` (last item) — no swipe should fire.
3. **Issue 3**: Chrome DevTools → Memory → take heap snapshot before opening the modal, swipe through 30 images, close modal, take a second snapshot. Detached image bytes should drop near zero. Also watch `performance.memory.usedJSHeapSize` during the swipe sequence — should plateau, not grow linearly.
4. **Zoom still works**: pinch-zoom on an image, then verify swipe-to-next is correctly blocked (via `canNavigate`), reset transform on double-tap (existing behavior), and that `mediaZoomRef.reset()` works on media change.

Automated:

- `cd app && npm run lint && npm run build` after changes.
- No unit tests cover this component today; not in scope to add.

---

## Out of scope (for this change)

- Replacing `zoompinch` (deferred per Issue 4 analysis).
- Adding swipe-tracking visual feedback (image follows finger before threshold).
- Refactoring the parent `library-page.vue` media list management.
