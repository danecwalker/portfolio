<script lang="ts">
  import {
    fetchLinks,
    type Link,
  } from "$lib/api";
  import Affiliation from "$lib/partials/Affiliation.svelte";
  import Experience from "$lib/partials/Experience.svelte";
  import Projects from "$lib/partials/Projects.svelte";
  import { onMount } from "svelte";

  let links: Link[] = $state([]);

  onMount(() => {
    fetchLinks().then((linksResp) => {
      links = linksResp.links;
    });
  });
</script>

<div class="w-full min-h-screen bg-background">
  <div class="p-8 max-w-7xl mx-auto flex justify-between items-center">
    <div class="font-medium text-text">Dane Walker</div>
    <nav>
      <ul class="flex items-center gap-8">
        <li>
          <a href="#projects" class="text-link hover:underline">Projects</a>
        </li>
        <li>
          <a href="#experience" class="text-link hover:underline">Experience</a>
        </li>
        <li>
          <a href="#affiliations" class="text-link hover:underline">Affiliations</a>
        </li>
      </ul>
    </nav>
  </div>

  <Projects />

  <Experience />

  <Affiliation />

  <footer>
    <div
      class="w-full border-t border-border mt-16 py-8 text-center border-inner-border text-sm text-footnote"
    >
      <ul class="flex items-center gap-8 w-full max-w-7xl mx-auto justify-end">
        {#each links as link}
          <li>
            {#if link.type === "file"}
              <a
                href={link.url}
                class="hover:text-link hover:underline duration-100"
                target="_blank"
                rel="noopener noreferrer">{link.title}</a
              >
            {:else}
              <a
                href={link.url}
                class="hover:text-link hover:underline duration-100"
                >{link.title}</a
              >
            {/if}
          </li>
        {/each}
      </ul>
    </div>
  </footer>
</div>
