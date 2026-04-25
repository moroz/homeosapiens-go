import { mount } from "svelte";
import "./app.css";
import App from "./App.svelte";

const app = mount(App, {
  target: document.getElementById("app")!,
  props: {
    pp: "https://d3n1g0yg3ja4p3.cloudfront.net/019beef9-ad4c-736f-9bb0-965b59ca21ae.png",
    title: "What prevents me from moving on?",
    date: "2026-02-08",
    locale: "en",
    host: "Dr Asher Shaikh",
  },
});

export default app;
