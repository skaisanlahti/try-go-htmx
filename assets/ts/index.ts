import htmx from "htmx.org";
declare global {
  interface Window {
    htmx: typeof htmx;
  }
}

window.htmx = htmx;
console.log("scripts loaded", new Date().toLocaleTimeString());
