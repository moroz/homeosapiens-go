import "./app.css";

import videojs from "video.js";
import "vanilla-hamburger";

document.querySelectorAll(".video-js").forEach((el) => {
	videojs(el, { controls: true, fluid: true });
});

document.querySelectorAll(".user-dropdown").forEach((el) => {
	const dropdown = el.querySelector(".dropdown");
	if (!dropdown) return;
	let timeout: number;

	el.addEventListener("mouseover", () => {
		clearTimeout(timeout);
		dropdown.classList.remove("hidden");
	});

	el.addEventListener("mouseout", () => {
		timeout = setTimeout(() => {
			dropdown.classList.add("hidden");
		}, 500);
	});
});

async function setTimezone() {
	const localTz = new Intl.DateTimeFormat().resolvedOptions().timeZone;
	const storedTz = document.querySelector("meta[name=user-timezone]")?.getAttribute("content");
	if (localTz === storedTz || !localTz) return;

	const qs = new URLSearchParams({ tz: localTz }).toString();

	fetch(`/api/v1/prefs/timezone?${qs}`, { method: "POST", credentials: "include" });
}

setTimezone();
