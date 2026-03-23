<script lang="ts">
  import { Search, Plus } from "@lucide/svelte";
  import * as Select from "$lib/components/ui/select/index";
  import { Button } from "$lib/components/ui/button/index";
  import { Input } from "$lib/components/ui/input/index";

  let searchQuery = $state("");
  let selectedCategory = $state<string>("All");
  const categories = ["All", "Cars", "Tracks"];
</script>

<div class="p-6 mx-50">
  <div class="flex items-center justify-between mb-6 gap-4">
    <div class="flex items-center gap-4">
      <!-- Search bar -->
      <div class="relative">
        <Search
          class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground"
        />
        <Input
          type="text"
          placeholder="Search mods..."
          class="pl-10 min-w-64 h-8"
          bind:value={searchQuery}
        />
      </div>

      <!-- Category selection -->
      <Select.Root type="single" bind:value={selectedCategory}>
        <Select.Trigger class="w-40 min-h-8"
          >{selectedCategory === "All"
            ? "All Categories"
            : selectedCategory}</Select.Trigger
        >
        <Select.Content>
          {#each categories as category}
            <Select.Item value={category}>
              {category === "All" ? "All Categories" : category}
            </Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>

    <!-- Add mod button -->
    <div>
      <Button class="h-8 min-w-30">
        <Plus class="size-4 mr-2" />
        Add Mod
      </Button>
    </div>
  </div>
</div>
