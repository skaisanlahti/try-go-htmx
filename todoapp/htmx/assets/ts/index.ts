import htmx from "htmx.org";
import { hello } from "./todo_page";

declare global {
  interface Window {
    htmx: typeof htmx;
  }
}

window.htmx = htmx;
hello();
