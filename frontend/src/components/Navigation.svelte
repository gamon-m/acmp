<script lang="ts">
  import { Settings2 } from "@lucide/svelte";
  import * as NavigationMenu from "$lib/components/ui/navigation-menu/index";

  type Page = "profiles" | "mods" | "settings";

  let { activeTab = $bindable<Page>("profiles") } = $props();
</script>

<nav
  class="grid grid-cols-[1fr_auto_1fr] items-center px-12 py-4 border-b border-border bg-header"
>
  <h1 class="text-3xl font-bold text-primary">AC MOD PROFILES</h1>

  <NavigationMenu.Root>
    <NavigationMenu.List>
      <NavigationMenu.Item>
        <NavigationMenu.Link
          href="#"
          data-active={activeTab === "profiles" ? "true" : "false"}
          onclick={(e) => {
            e.preventDefault();
            activeTab = "profiles";
          }}
        >
          PROFILES
        </NavigationMenu.Link>
      </NavigationMenu.Item>
      <NavigationMenu.Item>
        <NavigationMenu.Link
          href="#"
          data-active={activeTab === "mods" ? "true" : "false"}
          onclick={(e) => {
            e.preventDefault();
            activeTab = "mods";
          }}
        >
          MODS
        </NavigationMenu.Link>
      </NavigationMenu.Item>
    </NavigationMenu.List>
  </NavigationMenu.Root>

  <button
    class="hover:text-muted-foreground rounded-md transition-colors justify-self-end cursor-pointer data-[active=true]:text-primary"
    data-active={activeTab === "settings" ? "true" : "false"}
    onclick={(e) => {
      e.preventDefault();
      activeTab = "settings";
    }}
  >
    <Settings2 class="size-8" />
  </button>
</nav>
