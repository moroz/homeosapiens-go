<script lang="ts">
  import { Temporal } from "@js-temporal/polyfill";
  import { twMerge } from "tailwind-merge";

  interface Props {
    title: string;
    pp: string;
    locale: "en" | "pl";
    host: string;
    date: string | Temporal.PlainDate;
  }

  const { pp, date, title: baseTitle, host, locale }: Props = $props();

  const dateParsed = $derived(Temporal.PlainDate.from(date));

  const [title, subtitle] = $derived.by(() => {
    const [t, s] = baseTitle.split(":").map((t) => t.trim());
    if (s && /\b\d\b/iu.test(s)) return [t, s];
    return [baseTitle];
  });

  const smallerText = Boolean(subtitle) || baseTitle.length > 33;
</script>

<div
  class="aspect-video bg-slate-900 w-160 flex overflow-hidden outline-pink-300 outline-3"
>
  <!-- Host photo -->
  <div class="relative h-full shrink-0 overflow-hidden p-4 w-68">
    <img
      src={pp}
      class="w-full h-full object-cover rounded-md outline outline-white/70"
      alt=""
    />
  </div>

  <!-- Content -->
  <div class="flex-1 flex flex-col justify-between py-6 pr-8 pl-2">
    <!-- Logo -->
    <svg
      viewBox="0 0 1538 361"
      class={twMerge("w-[220px] text-white h-auto", smallerText && "w-50")}
    >
      <use href="/logo.svg#logo" />
    </svg>

    <!-- Title -->
    <h1
      class={twMerge(
        "text-white font-bold text-[40px] leading-tight",
        smallerText && "text-4xl",
      )}
    >
      {title}
      {#if subtitle}
        <small class="text-white/90 mt-3 block">{subtitle}</small>
      {/if}
    </h1>

    <!-- Host & Date -->
    <div class="flex flex-col gap-1">
      <p
        class={twMerge(
          "text-slate-300 font-semibold text-3xl",
          smallerText && "text-2xl",
        )}
      >
        {host}
      </p>
      <p
        class={twMerge(
          "text-slate-100 font-semibold text-3xl",
          smallerText && "text-2xl",
        )}
      >
        {dateParsed.toLocaleString(locale, { dateStyle: "long" })}
      </p>
    </div>
  </div>
</div>
