<template>
  <span class="text-fit">
    <span><span v-html="text" /></span>
    <span
      aria-hidden="true"
      v-html="text"
    />
  </span>
</template>

<script setup lang="ts">
defineProps<{
  text: string;
}>();
</script>

<style>
/* Thanks to https://kizu.dev/fit-to-width/ */
.text-fit:not(.text-fit *) {
  padding: 0 0.5em;
  --max-font-size: 10em;
}

.text-fit {
  display: flex;
  container-type: inline-size;
  padding: 0 0.5em;
  line-height: 0.85;

  --captured-length: initial;
  --support-sentinel: var(--captured-length, 9999px);

  & > [aria-hidden] {
    visibility: hidden;
  }

  & > :not([aria-hidden]) {
    flex-grow: 1;
    container-type: inline-size;

    --captured-length: 100cqi;
    --available-space: var(--captured-length);

    & > * {
      --support-sentinel: inherit;
      --captured-length: 100cqi;
      --ratio: tan(atan2(var(--available-space), var(--available-space) - var(--captured-length)));
      --font-size: clamp(1em, 1em * var(--ratio), var(--max-font-size, infinity * 1px) - var(--support-sentinel));
      inline-size: var(--available-space);

      &:not(.text-fit) {
        display: block;
        font-size: var(--font-size);

        @container (inline-size > 0) {
          white-space: nowrap;
        }
      }

      &.text-fit {
        --captured-length2: var(--font-size);
        font-variation-settings: 'opsz' tan(atan2(var(--captured-length2), 1px));
      }
    }
  }
}

@property --captured-length {
  syntax: '<length>';
  initial-value: 0px;
  inherits: true;
}

@property --captured-length2 {
  syntax: '<length>';
  initial-value: 0px;
  inherits: true;
}
</style>
