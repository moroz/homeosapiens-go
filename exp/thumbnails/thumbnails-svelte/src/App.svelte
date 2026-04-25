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

  const dateParsed = Temporal.PlainDate.from(date);
</script>

<div class="w-full h-screen grid place-items-center">
  <div class="aspect-video bg-slate-100 w-160 outline-2 flex p-8 gap-8">
    <div class="h-full aspect-3/4 outline flex overflow-hidden shrink-0">
      <img src={pp} class="w-full h-full object-cover" alt="" />
    </div>
    <div class="flex flex-col gap-4">
      <h1 class="font-bold text-3xl">{host}</h1>
      <h2 class="font-bold text-4xl leading-tight">{title}</h2>
      <h3>{dateParsed.toLocaleString(locale, { dateStyle: "long" })}</h3>
    </div>
  </div>
</div>
