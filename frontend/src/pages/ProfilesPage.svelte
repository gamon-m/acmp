<script lang="ts">
  import {
    Search,
    Plus,
    Play,
    Pencil,
    Trash2,
    EllipsisVertical,
  } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button/index";
  import { Input } from "$lib/components/ui/input/index";
  import * as Card from "$lib/components/ui/card/index";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index";

  interface Profile {
    id: string;
    name: string;
    category: string;
    modCount: number;
    active: boolean;
  }

  let searchQuery = $state("");

  const mockProfiles: Profile[] = $state([
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
    {
      id: "5",
      name: "Drift",
      category: "Cards",
      modCount: 12,
      active: false,
    },
    {
      id: "6",
      name: "Touge",
      category: "Tracks",
      modCount: 45,
      active: true,
    },
    {
      id: "7",
      name: "Other",
      category: "Cars",
      modCount: 8,
      active: false,
    },
    {
      id: "8",
      name: "Drift",
      category: "Tracks",
      modCount: 167,
      active: false,
    },
  ]);

  let filteredProfiles = $derived(
    mockProfiles.filter((p) =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()),
    ),
  );

  function toggleActive(profile: Profile) {
    profile.active = !profile.active;
  }

  function deleteProfile(profile: Profile) {
    console.log("deleting profile:", profile.name);
  }

  function editProfile(profile: Profile) {
    console.log("deleting profile:", profile.name);
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
          class="pl-10 w-64 h-8"
          bind:value={searchQuery}
        />
      </div>
      <Button class="h-8">
        <Plus class="size-4 mr-2" />
        Add Profile
      </Button>
    </div>
  </div>

  <div class="space-y-2">
    {#each filteredProfiles as profile (profile.id)}
      <Card.Root
        class="transition-colors border-l-6 hover:bg-accent/50 data-[active=true]:bg-muted data-[active=true]:hover:bg-accent data-[active=true]:border-primary"
        data-active={profile.active ? "true" : "false"}
      >
        <Card.Content class="grid grid-cols-[1fr_100px_100px_auto] gap-6 py-3">
          <div class="border-r border-border pr-4">
            <span class="text-xs text-muted-foreground font-medium"
              >PROFILE NAME</span
            >
            <p class="text-sm truncate">{profile.name}</p>
          </div>
          <div class="border-r border-border pr-4 text-center">
            <span class="text-xs text-muted-foreground font-medium"
              >CATEGORY</span
            >
            <p class="text-sm">{profile.category}</p>
          </div>
          <div class="border-r border-border pr-4 text-center">
            <span class="text-xs text-muted-foreground font-medium"
              >ACTIVE MODS</span
            >
            <p class="text-sm">{profile.modCount} mods</p>
          </div>
          <div class="flex items-center gap-2">
            <Button
              variant={profile.active ? "default" : "outline"}
              size="sm"
              onclick={() => toggleActive(profile)}
              class="min-w-25 cursor-pointer h-8"
            >
              <Play class="w-5 mr-1" />
              {profile.active ? "Active" : "Activate"}
            </Button>
            <DropdownMenu.Root>
              <DropdownMenu.Trigger>
                <Button variant="ghost" size="sm" class="h-8 w-8 p-0">
                  <EllipsisVertical class="size-4" />
                </Button>
              </DropdownMenu.Trigger>
              <DropdownMenu.Content>
                <DropdownMenu.Item onclick={() => editProfile(profile)}>
                  <Pencil class="size-4 mr-2" />
                  Edit
                </DropdownMenu.Item>
                <DropdownMenu.Separator />
                <DropdownMenu.Item
                  class="text-destructive focus:text-destructive"
                  onclick={() => deleteProfile(profile)}
                >
                  <Trash2 class="size-4 mr-2" />
                  Delete
                </DropdownMenu.Item>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          </div>
        </Card.Content>
      </Card.Root>
    {/each}
  </div>
  <div class="min-h-20"></div>
</div>
