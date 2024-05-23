/* @refresh reload */
import { render } from "solid-js/web";

import "./index.css";
import App from "./App";
import { Toaster } from "./components/ui/toast";
import axios from "axios";

if (import.meta.env.VITE_BACKEND_HOST === undefined) {
  throw new Error("please set VITE_BACKEND_HOST variable");
}
axios.defaults.withCredentials = true;

const root = document.getElementById("root");

render(
  () => (
    <>
      <App />
      <Toaster />
    </>
  ),
  root!
);
