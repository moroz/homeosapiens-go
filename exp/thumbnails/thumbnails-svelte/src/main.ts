import { mount, unmount } from "svelte";
import "./app.css";
import App from "./App.svelte";

export interface ThumbnailProps {
  pp: string;
  title: string;
  date: string;
  host: string;
  locale: "en" | "pl";
}

let app: any = null;

export function renderThumbnail(props: ThumbnailProps | string) {
  if (typeof props === "string") {
    props = JSON.parse(props) as ThumbnailProps;
  }
  app = mount(App, {
    target: document.getElementById("app")!,
    props,
  });
}

export async function destroyApp() {
  if (!app) return;
  unmount(app);
  app = null;
}
