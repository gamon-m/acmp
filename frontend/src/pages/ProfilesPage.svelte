<script lang="ts">
  import { Search, Plus, Play, Pencil } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button/index";
  import { Input } from "$lib/components/ui/input/index";

  interface Profile {
    id: string;
    name: string;
    category: string;
    modCount: number;
    active: boolean;
  }

  let searchQuery = $state("");

  const mockProfiles: Profile[] = [
    {
      id: "1",
      name: "Drift",
      category: "Cards",
      modCount: 12,
      active: false,
    },
    {
      id: "2",
      name: "Touge",
      category: "Tracks",
      modCount: 45,
      active: true,
    },
    {
      id: "3",
      name: "Other",
      category: "Cars",
      modCount: 8,
      active: false,
    },
    {
      id: "4",
      name: "Drift",
      category: "Tracks",
      modCount: 167,
      active: false,
    },
  ];

  let filteredProfiles = $derived(
    mockProfiles.filter((p) =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()),
    ),
  );

  function toggleActive(profile: Profile) {
    profile.active = !profile.active;
  }

  function editProfile(profile: Profile) {
    console.log("Edit profile:", profile.name);
  }
</script>

<div class="p-6 mx-50">
  <div class="flex items-center justify-between mb-6">
    <h2 class="text-2xl font-semibold">PROFILES</h2>
    <div class="flex items-center gap-4">
      <div class="relative">
        <Search
          class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground"
        />
        <Input
          type="text"
          placeholder="Search profiles..."
          class="pl-10 w-64"
          bind:value={searchQuery}
        />
      </div>
      <Button>
        <Plus class="size-4 mr-2" />
        Add Profile
      </Button>
    </div>
  </div>

  <div class="space-y-1">
    {#each filteredProfiles as profile (profile.id)}
      <div
        class="grid grid-cols-[1fr_120px_100px_auto] gap-4 items-center p-4 rounded-lg transition-colors border-l-8 {profile.active
          ? 'bg-muted border-primary hover:bg-accent'
          : 'bg-card hover:bg-accent/25'}"
      >
        <span class="text-sm border-r border-border pr-4 truncate"
          >{profile.name}</span
        >
        <span
          class="text-sm text-muted-foreground border-r border-border pr-4 text-center"
          >{profile.category}</span
        >
        <span
          class="text-sm text-muted-foreground border-r border-border pr-4 text-center"
          >{profile.modCount} mods</span
        >
        <div class="flex items-center gap-2">
          <Button
            variant={profile.active ? "default" : "outline"}
            size="sm"
            onclick={() => toggleActive(profile)}
            class="min-w-25"
          >
            <Play class="w-5 mr-1" />
            {profile.active ? "Active" : "Activate"}
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onclick={() => editProfile(profile)}
          >
            <Pencil class="size-4" />
          </Button>
        </div>
      </div>
    {/each}
  </div>
</div>
