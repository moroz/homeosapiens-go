<script lang="ts">
  import { Temporal } from "@js-temporal/polyfill";

  interface Props {
    title: string;
    pp: string;
    locale: "en" | "pl";
    host: string;
    date: string | Temporal.PlainDate;
  }

  const { pp, date, title, host, locale }: Props = $props();

  const dateParsed = $derived(Temporal.PlainDate.from(date));
</script>

<div class="w-full h-screen grid place-items-center bg-slate-100">
  <div
    class="aspect-video bg-indigo-950 w-160 flex overflow-hidden outline-pink-300 outline-3"
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
        style="color: white; width: 220px; height: auto"
      >
        <use href="/logo.svg#logo" />
      </svg>

      <!-- Title -->
      <div class="flex flex-col gap-3">
        <h1 class="text-white font-bold text-[40px] leading-tight">{title}</h1>
      </div>

      <!-- Host & Date -->
      <div class="flex flex-col gap-1">
        <p class="text-indigo-200 font-semibold text-3xl">{host}</p>
        <p class="text-white/90 font-semibold text-3xl">
          {dateParsed.toLocaleString(locale, { dateStyle: "long" })}
        </p>
      </div>
    </div>
  </div>
</div>
