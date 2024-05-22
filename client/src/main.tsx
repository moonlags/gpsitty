import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { Toaster } from "sonner";

if (process.env.REACT_APP_BACKEND_HOST == null) {
  throw new Error("Please set REACT_APP_BACKEND_HOST variable");
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
    <Toaster />
  </React.StrictMode>
);
