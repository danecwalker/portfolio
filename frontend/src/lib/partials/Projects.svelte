<script lang="ts">
  import { fetchProjects, type Project } from "$lib/api";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";

  let parallaxTargetRef: HTMLImageElement | null = $state(null);
  let hoveredCard: number = $state(0);
  let projects: Project[] = $state([]);
  let imagesLoaded = $state(0);
  let imagesReady = $state(false);
  let loadedMap = new Map<string, boolean>();

  const MAX_MOVE = 40; // px – how far the image can move
  let currentX = 0;
  let currentY = 0;
  let targetX = 0;
  let targetY = 0;
  const easeNorm = 0.05; // slower glide when returning to center
  const easeHome = 0.025; // slower glide when returning to center
  let ease = easeNorm; // normal follow speed
  let frameId: number;

  const startParallaxLoop = () => {
    const loop = () => {
      // interpolate towards target (0.15 = smoothing factor)
      currentX += (targetX - currentX) * ease;
      currentY += (targetY - currentY) * ease;

      if (parallaxTargetRef) {
        parallaxTargetRef.style.transform = `translate(calc(-50% + ${currentX}px), calc(-50% + ${currentY}px))`;
      }

      frameId = requestAnimationFrame(loop);
    };

    frameId = requestAnimationFrame(loop);
  };

  const handleMouseMove = (e: MouseEvent) => {
    ease = easeNorm; // restore normal speed

    if (!parallaxTargetRef) return;
    const rect = (e.currentTarget as HTMLDivElement).getBoundingClientRect();

    // Normalised position of mouse over button: -0.5 to 0.5
    const x = (e.clientX - rect.left) / rect.width - 0.5;
    const y = (e.clientY - rect.top) / rect.height - 0.5;

    targetX = x * MAX_MOVE;
    targetY = y * MAX_MOVE;
  };

  const handleMouseLeave = () => {
    ease = easeHome; // switch to slow mode
    targetX = 0;
    targetY = 0;
  };

  const preloadImages = (urls: string[], timeout: number = 0) => {
    return Promise.all(
      urls.map(
        (url) =>
          new Promise<void>(async (resolve) => {
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
    fetchProjects().then((projectsResp) => {
      projects = projectsResp.projects;
      preloadImages(projects.map((p) => p.projectImageURL)).then(() => {
        setTimeout(() => {
          imagesReady = true;
        }, 500);
      });
    });

    startParallaxLoop();

    return () => {
      cancelAnimationFrame(frameId);
    };
  });
</script>

<section
  id="projects"
  class="w-full max-w-7xl min-h-[760px] mx-auto p-8 text-text"
>
  <div class="w-full flex items-start">
    <span class="text-sm mr-2 text-footnote">[1]</span>
    <h1 class="font-bold text-2xl">Projects</h1>
  </div>

  <div class="flex gap-8 mt-8">
    {#if !imagesReady && projects.length > 0}
      <div class="w-full h-full flex items-center justify-center">
        <div class="grid grid-cols-4 gap-1">
          {#each projects as _, i}
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
      <div class="lg:flex lg:flex-col lg:w-1/3 w-full gap-8 lg:gap-0 grid grid-cols-1 md:grid-cols-2">
        {#each projects as project, i}
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <a href="{project.projectURL}" target="_blank" rel="noopener noreferrer" data-umami-event="project-link-click" data-umami-event-url={project.projectURL}>
            <img 
              src={project.projectImageURL}
              alt={project.title}
              class="w-full aspect-5/4 object-cover rounded-lg block lg:hidden"
            >
            <div
              onmouseenter={() => (hoveredCard = i)}
              onmousemove={handleMouseMove}
              onmouseleave={handleMouseLeave}
              data-selected={hoveredCard === i}
              class="lg:data-[selected=true]:bg-inner bg-transparent py-4 px-0 lg:px-4 lg:timing-card rounded-lg"
            >
              <div id="card-info">
                <h3 class="font-semibold text-lg">{project.title}</h3>
                <p class="text-sm text-footnote">{project.date}</p>
              </div>
            </div>
          </a>
        {/each}
      </div>

      <div class="relative w-2/3 flex-1 hidden lg:block">
        {#key hoveredCard}
          <img
            in:fade={{ duration: 300 }}
            out:fade={{ duration: 300 }}
            bind:this={parallaxTargetRef}
            src={projects[hoveredCard || 0]?.projectImageURL}
            alt={projects[hoveredCard || 0]?.title}
            class="absolute top-1/2 left-1/2 w-full aspect-5/4 object-cover rounded-lg shadow-lg will-change-transform {imagesReady
              ? 'opacity-100'
              : 'opacity-0'}"
            style="transform: translate(-50%, -50%);"
          />
        {/key}
      </div>
    {/if}
  </div>
</section>
