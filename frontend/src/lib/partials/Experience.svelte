<script lang="ts">
  import { fetchContent, fetchExperiences, type Experience } from "$lib/api";
  import CompanyIcon from "$lib/components/CompanyIcon.svelte";
  import { onMount } from "svelte";
  import { scale } from "svelte/transition";

  let hoveredCard: number = $state(0);

  let experiences: Experience[] = $state([]);
  let imagesLoaded = $state(0);
  let imagesReady = $state(false);
  let loadedMap = new Map<string, boolean>();

  const preloadImages = (urls: string[], timeout: number = 0) => {
    return Promise.all(
      urls.map(
        (url) =>
          new Promise<void>(async (resolve) => {
            const e = experiences.find((exp) => exp.logoURL === url);
            if (e) {
              await fetchContent(e.pageId).then((content) => {
                e.hightlights = content.content;
              });
            }
            const img = new Image();
            img.src = url;
            img.onload = async () => {
              loadedMap.set(url, true);
              imagesLoaded += 1;
              await new Promise((r) => setTimeout(r, timeout)); // fake delay
              resolve();
            };
            img.onerror = async () => {
              loadedMap.set(url, false);
              imagesLoaded += 1;
              await new Promise((r) => setTimeout(r, timeout)); // fake delay
              resolve();
            };
          })
      )
    );
  };

  onMount(() => {
    fetchExperiences().then((experiencesResp) => {
      experiences = experiencesResp.experience;
      preloadImages(experiences.map((p) => p.logoURL)).then(() => {
        setTimeout(() => {
          imagesReady = true;
        }, 500);
      });
    });
  });
</script>

<section
  id="experience"
  class="w-full max-w-7xl mx-auto p-8 text-text min-h-[400px]"
>
  <div class="w-full flex items-start">
    <span class="text-sm mr-2 text-footnote">[2]</span>
    <h1 class="font-bold text-2xl">Experience</h1>
  </div>

  <div class="flex gap-8 mt-8">
    {#if !imagesReady && experiences.length > 0}
      <div class="w-full h-full flex items-center justify-center">
        <div class="grid grid-cols-5 gap-1">
          {#each experiences as _, i}
            <div class="flip-card {i < imagesLoaded ? 'loaded' : ''}">
              <div class="flip-card-inner">
                <div class="flip-card-front"></div>
                <div class="flip-card-back"></div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {:else}
      <div class="lg:flex lg:flex-col lg:w-1/2 w-full gap-8 grid grid-cols-1">
        {#each experiences as experience, i}
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div>
            <div
              data-selected={hoveredCard === i}
              onmouseenter={() => (hoveredCard = i)}
              class="flex gap-4 items-center lg:data-[selected=true]:bg-inner bg-transparent py-4 lg:px-4 px-0 lg:timing-card rounded-lg flex-row-reverse lg:flex-row justify-between lg:justify-start cursor-pointer"
            >
              <div class="w-12 aspect-square rounded-lg overflow-hidden">
                {#if !experience.logoURL}
                  <div
                    class="w-full h-full bg-inner flex items-center justify-center text-footnote text-sm"
                  >
                    <CompanyIcon />
                  </div>
                {:else}
                  <img
                    src={experience.logoURL}
                    alt={experience.company}
                    class="w-full h-full object-contain"
                  />
                {/if}
              </div>
              <div>
                <h3 class="font-semibold text-lg">
                  {experience.position}
                  <span class="hidden lg:block">@ {experience.company}</span>
                </h3>
                <p class="text-sm text-footnote">
                  {experience.start} - {experience.end ?? "Present"}
                </p>
              </div>
            </div>

            <div class="lg:hidden block">
              <h3 class="font-semibold">{experience.company}</h3>

              <div class="experience-highlights mt-4">
                {#if experience.hightlights}
                  {@html experience.hightlights}
                {/if}
              </div>
            </div>
          </div>
        {/each}
      </div>

      <div class="relative w-1/2 flex-1 h-[440px] hidden lg:block">
        {#key hoveredCard}
          <div
            in:scale={{ start: 0.9, duration: 500 }}
            class="absolute top-0 left-0 {imagesReady
              ? 'opacity-100'
              : 'opacity-0'}"
          >
            {#if !experiences[hoveredCard || 0]?.logoURL}
              <div
                class="w-24 h-24 bg-inner flex items-center justify-center text-footnote text-sm mb-4"
              >
                <CompanyIcon />
              </div>
            {:else}
              <img
                src={experiences[hoveredCard || 0]?.logoURL}
                alt={experiences[hoveredCard || 0]?.company}
                class="w-24 h-24 object-contain mb-4"
              />
            {/if}
            <h4 class="text-lg font-medium">
              {experiences[hoveredCard || 0]?.company}
            </h4>

            <div class="experience-highlights mt-4">
              {#if experiences[hoveredCard || 0]?.hightlights}
                {@html experiences[hoveredCard || 0]?.hightlights}
              {/if}
            </div>
          </div>
        {/key}
      </div>
    {/if}
  </div>
</section>
