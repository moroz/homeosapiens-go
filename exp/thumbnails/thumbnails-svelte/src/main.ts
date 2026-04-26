import { mount, unmount, tick } from "svelte";
import "./app.css";
import App from "./Thumbnail.svelte";
import Test from "./Test.svelte";

export interface ThumbnailProps {
  pp: string;
  title: string;
  date: string;
  host: string;
  locale: "en" | "pl";
}

let app: any = null;

export async function renderThumbnail(props: ThumbnailProps | string) {
  if (typeof props === "string") {
    props = JSON.parse(props) as ThumbnailProps;
  }
  app = mount(App, {
    target: document.getElementById("app")!,
    props,
  });
  await tick();
}

export async function destroyApp() {
  if (!app) return;
  await unmount(app);
  app = null;
}

export function renderTest() {
  mount(Test, {
    target: document.getElementById("app")!,
  });
}

if (import.meta.env.DEV && !navigator.webdriver) {
  renderTest();
}
