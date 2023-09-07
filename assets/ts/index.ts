import htmx from "htmx.org";
import { hello } from "../../todos/templates/todo_page";
declare global {
  interface Window {
    htmx: typeof htmx;
  }
}

window.htmx = htmx;
hello();
