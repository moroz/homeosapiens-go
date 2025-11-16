import "./app.css";

import videojs from "video.js";

document.querySelectorAll(".video-js").forEach((el) => {
	videojs(el, { controls: true, fluid: true });
});
