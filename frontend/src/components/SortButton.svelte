<script lang="ts">
  import { ChevronUp } from "@lucide/svelte";
  let {
    displayName,
    catName,
    sortField,
    toggleSort,
    sortDirection,
    center = true,
  } = $props();
</script>

<button
  class="flex items-center gap-1 {center
    ? 'justify-center'
    : ''} transition-colors {sortField === catName
    ? 'text-accent-foreground hover:text-foreground'
    : 'hover:text-accent-foreground'}"
  onclick={() => toggleSort(catName)}
>
  {displayName}
  <ChevronUp
    class="size-3 transition-transform {sortField !== catName
      ? 'invisible'
      : ''} {sortDirection === 'desc' ? 'rotate-180' : ''}"
  />
</button>
